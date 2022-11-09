package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"nebula/pkg/api/apiutil"
	"nebula/pkg/api/model"
	"nebula/pkg/security"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type AuthService interface {
	GetUserForEmail(context.Context, string) (*model.User, error)
	GetUserAndSession(context.Context, uuid.UUID) (*model.Session, *model.User, error)

	UpdateUserPassword(context.Context, uuid.UUID, string) error

	CreateSession(context.Context, uuid.UUID, security.HWID) (uuid.UUID, error)
	RefreshSession(context.Context, uuid.UUID, uuid.UUID) error
	GetActiveSessions(context.Context, uuid.UUID) ([]string, error)
	IsSessionActive(context.Context, uuid.UUID, uuid.UUID) (bool, error)
	ClearSessions(context.Context, uuid.UUID, ...uuid.UUID) error
	VerifySession(context.Context, uuid.UUID) error
	MaxInstancesForUser(context.Context, uuid.UUID) (int, error)

	LockAccount(context.Context, uuid.UUID, time.Duration) (*redislock.Lock, error)

	InvalidatePreviousResetTokens(context.Context, uuid.UUID) error
	CreateAndSendPasswordReset(context.Context, uuid.UUID, string) error
	IsValidResetToken(context.Context, string) bool
	GetUserIdFromResetToken(context.Context, string) (uuid.UUID, error)

	GetOtpForUser(context.Context, uuid.UUID) (*model.OtpData, error)
	SetOtpForUser(context.Context, uuid.UUID, string, []byte) error

	HashUser(email string) string
	TLDashAuthorization() string

	RegisterAndSendConfirmation(context.Context, string, string, string) error
	CompleteRegistration(context.Context, uuid.UUID) error
}

func NewAuthService(pgPool *sql.DB, redisClient *redis.Client, redisLock *redislock.Client, sendClient *sendgrid.Client, intercomSecret, TLDashUser, TLDashPass []byte) AuthService {
	return &authenticator{pgPool, redisClient, redisLock, sendClient, intercomSecret, TLDashUser, TLDashPass}
}

type authenticator struct {
	PgPool         *sql.DB
	RedisClient    *redis.Client
	RedisLock      *redislock.Client
	SendClient     *sendgrid.Client
	IntercomSecret []byte
	TLDashUser     []byte
	TLDashPass     []byte
}

func (a *authenticator) GetUserForEmail(ctx context.Context, email string) (*model.User, error) {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var user model.User
	err = conn.QueryRowContext(ctx, "SELECT uid, password_hash, key FROM public.users WHERE email = $1", email).Scan(&user.Uid, &user.PasswordHash, &user.Key)
	if err != nil {
		return nil, err
	}

	user.Email = email
	return &user, nil
}

func (a *authenticator) GetUserAndSession(ctx context.Context, sid uuid.UUID) (*model.Session, *model.User, error) {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	var user model.User
	var sess model.Session
	var hwid []byte
	err = conn.QueryRowContext(ctx, "SELECT sessions.sid, sessions.uid, sessions.hwid, sessions.verified, users.email, users.key FROM public.sessions INNER JOIN public.users ON sessions.uid = users.uid WHERE sid = $1", sid).Scan(&sess.Sid, &sess.Uid, &hwid, &sess.Verified, &user.Email, &user.Key)
	if err != nil {
		return nil, nil, err
	}

	if err = json.Unmarshal(hwid, &sess.Hwid); err != nil {
		return nil, nil, err
	}

	user.Uid = sess.Uid

	return &sess, &user, nil
}

func (a *authenticator) CreateSession(ctx context.Context, uid uuid.UUID, hwid security.HWID) (uuid.UUID, error) {
	sessionId := uuid.New()

	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return sessionId, err
	}
	defer conn.Close()

	hwidData, err := json.Marshal(hwid)
	if err != nil {
		return sessionId, nil
	}

	_, err = conn.ExecContext(ctx, "INSERT INTO public.sessions(sid, uid, hwid) VALUES ($1, $2, $3)", sessionId, uid, hwidData)
	if err != nil {
		return sessionId, err
	}

	return sessionId, nil
}

func (a *authenticator) RefreshSession(ctx context.Context, uid uuid.UUID, sid uuid.UUID) error {
	return a.RedisClient.Expire(ctx, fmt.Sprintf("sessions:%s:%s", uid, sid), time.Second*30).Err()
}

func (a *authenticator) IsSessionActive(ctx context.Context, uid uuid.UUID, sid uuid.UUID) (bool, error) {
	cnt, err := a.RedisClient.Exists(ctx, fmt.Sprintf("sessions:%s:%s", uid, sid)).Result()
	return cnt > 0, err
}

func (a *authenticator) ClearSessions(ctx context.Context, uid uuid.UUID, exclude ...uuid.UUID) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	r, err := conn.QueryContext(ctx, "DELETE FROM public.sessions WHERE uid = $1 AND NOT sid = ANY( $2 ) RETURNING sid", uid, exclude)
	if err != nil {
		return err
	}

	var sessions []string
	for r.Next() {
		var sid uuid.UUID
		if err = r.Scan(&sid); err != nil {
			return err
		}
		sessions = append(sessions, sid.String())
	}
	if err = r.Err(); err != nil {
		return err
	}

	if len(sessions) == 0 {
		return nil
	}

	a.RedisClient.Del(ctx, sessions...)
	return err
}

func (a *authenticator) VerifySession(ctx context.Context, sid uuid.UUID) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE public.sessions SET verified=true WHERE sid = $1", sid)
	if err != nil {
		return err
	}

	return err
}

func (a *authenticator) MaxInstancesForUser(ctx context.Context, uid uuid.UUID) (int, error) {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return 1, err
	}
	defer conn.Close()

	maxInstances := 1
	err = conn.QueryRowContext(ctx, "SELECT count FROM public.instances WHERE uid = $1 AND (expires_at >= now() OR expires_at IS NULL) ORDER BY instances.count DESC LIMIT 1", uid).Scan(&maxInstances)
	if err != nil && err != sql.ErrNoRows {
		return maxInstances, err
	}
	return maxInstances, nil
}

func (a *authenticator) GetActiveSessions(ctx context.Context, uid uuid.UUID) ([]string, error) {
	return a.RedisClient.Keys(ctx, fmt.Sprintf("sessions:%s:*", uid)).Result()
}

func (a *authenticator) LockAccount(ctx context.Context, uid uuid.UUID, duration time.Duration) (*redislock.Lock, error) {
	return a.RedisLock.Obtain(ctx, fmt.Sprintf("uid:%s", uid), duration, &redislock.Options{
		RetryStrategy: redislock.LinearBackoff(time.Millisecond * 250),
	})
}

func (a *authenticator) InvalidatePreviousResetTokens(ctx context.Context, uid uuid.UUID) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx, "UPDATE public.password_resets SET active = false WHERE uid = $1 AND active = true", uid)
	return err
}

func (a *authenticator) CreateAndSendPasswordReset(ctx context.Context, uid uuid.UUID, email string) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	token := uuid.New()

	_, err = conn.ExecContext(ctx, "INSERT INTO public.password_resets(uid, token, active) VALUES ($1, $2, $3)", uid, token, true)
	if err != nil {
		return err
	}

	from := mail.NewEmail("Nebulabots", "no-reply@nebulabots.com")
	to := mail.NewEmail("", email)
	subject := "Password Reset"
	plainTextContent := fmt.Sprintf("Reset Token: %s", token)
	htmlContent := fmt.Sprintf("<strong>Reset Token: %s</strong>", token)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	_, err = a.SendClient.Send(message)
	return err
}

func (a *authenticator) UpdateUserPassword(ctx context.Context, uuid uuid.UUID, password string) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	hashedPw, err := apiutil.HashPassword(password, nil)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "UPDATE public.users SET password_hash = $1 WHERE uid = $2", hashedPw, uuid)
	return err
}

func (a *authenticator) IsValidResetToken(ctx context.Context, token string) bool {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return false
	}
	defer conn.Close()

	var id uuid.UUID
	err = conn.QueryRowContext(ctx, "SELECT uid FROM public.password_resets WHERE token = $1 AND active = TRUE", token).Scan(&id)
	return err == nil
}

func (a *authenticator) GetUserIdFromResetToken(ctx context.Context, token string) (uuid.UUID, error) {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Close()

	var user model.User
	err = conn.QueryRowContext(ctx, "SELECT uid FROM public.password_resets WHERE token = $1", token).Scan(&user.Uid)
	if err != nil {
		return uuid.Nil, err
	}

	return user.Uid, nil
}

func (a *authenticator) GetOtpForUser(ctx context.Context, uid uuid.UUID) (*model.OtpData, error) {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var otpData model.OtpData
	err = conn.QueryRowContext(ctx, "SELECT secret, uid, recovery_code FROM public.otp WHERE uid = $1", uid).Scan(&otpData.Secret, &otpData.Uid, &otpData.RecoveryCodes)
	if err != nil {
		return nil, err
	}

	return &otpData, nil
}

func (a *authenticator) SetOtpForUser(ctx context.Context, uid uuid.UUID, secret string, recoveryData []byte) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "INSERT INTO public.otp(uid, secret, recovery_code) VALUES ($1, $2, $3)", uid, secret, recoveryData)
	return err
}

func (a *authenticator) HashUser(email string) string {
	h := hmac.New(sha256.New, a.IntercomSecret)
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *authenticator) TLDashAuthorization() string {
	auth := string(a.TLDashUser) + ":" + string(a.TLDashPass)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (a *authenticator) RegisterAndSendConfirmation(ctx context.Context, email, password, key string) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	regId := uuid.New()

	hashedPassword, err := apiutil.HashPassword(password, nil)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "INSERT INTO public.pending_users(id, email, password_hash, key) VALUES ($1, $2, $3, $4)", regId, email, hashedPassword, key)
	if err != nil {
		return err
	}

	from := mail.NewEmail("Nebulabots", "no-reply@nebulabots.com")
	to := mail.NewEmail("", email)
	subject := "Registration Confirmation"
	plainTextContent := fmt.Sprintf("Thank you for signing up for Nebulabots! Please follow the link here https://nebulabots.uc.r.appspot.com/api/auth/register?id=%s to complete your registration. If you did not sign up for Nebulabots, please contact support in the community discord or email via hello@nebulabots.com.", regId)
	htmlContent := fmt.Sprintf("<div>Thank you for signing up for Nebulabots! Please follow the link <a href=\"https://nebulabots.uc.r.appspot.com/api/auth/register?id=%s\">here</a> to complete your registration. If you did not sign up for Nebulabots, please contact support in the community discord or email via hello@nebulabots.com. </div>", regId)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	_, err = a.SendClient.Send(message)
	return err
}

func (a *authenticator) CompleteRegistration(ctx context.Context, id uuid.UUID) error {
	conn, err := a.PgPool.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var email, passwordHash, key string
	row := conn.QueryRowContext(ctx, "SELECT email, password_hash, key FROM public.pending_users WHERE id = $1", id)
	if err = row.Scan(&email, &passwordHash, &key); err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "INSERT INTO public.users(email, password_hash, key) VALUES ($1, $2, $3)", email, passwordHash, key)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, "DELETE FROM public.pending_users WHERE id = $1", id)
	return err
}
