package controller

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"nebula/cmd/api/service"
	"nebula/pkg/api/apiutil"
	"nebula/pkg/api/client"
	"nebula/pkg/api/model"
	"nebula/pkg/security"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"github.com/pquerna/otp/totp"
)

type AuthController struct {
	Ctx         iris.Context
	AuthService service.AuthService
	RedisClient *redis.Client
}

type TLDashResponse struct {
	UserID          int         `json:"userId"`
	UserStatus      string      `json:"userStatus"`
	Email           string      `json:"email"`
	PlanType        string      `json:"planType"`
	DiscordID       string      `json:"discordId"`
	DiscordUsername string      `json:"discordUsername"`
	DiscordAvatar   string      `json:"discordAvatar"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	ExpiresAt       interface{} `json:"expiresAt"`
}

func (c *AuthController) Post(authRequest model.AuthRequest) model.AuthResponse {
	ctx := c.Ctx.Request().Context()

	user, err := c.AuthService.GetUserForEmail(ctx, authRequest.Email)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Invalid Email/Password"}
	}

	acctLock, err := c.AuthService.LockAccount(ctx, user.Uid, time.Second*5)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}
	defer acctLock.Release(ctx)

	maxInstances, err := c.AuthService.MaxInstancesForUser(ctx, user.Uid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	sessions, err := c.AuthService.GetActiveSessions(ctx, user.Uid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	} else if len(sessions) >= maxInstances {
		return model.AuthResponse{Success: false, Message: "Too Many Active Sessions"}
	}

	if ok, err := apiutil.VerifyPassword(authRequest.Pass, user.PasswordHash); err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	} else if !ok {
		return model.AuthResponse{Success: false, Message: "Invalid Email/Password"}
	}

	// TODO: check that key has valid subscription, possibly even before verifying password?

	sessionId, err := c.AuthService.CreateSession(ctx, user.Uid, authRequest.Hwid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	_, err = c.AuthService.GetOtpForUser(ctx, user.Uid)
	if err == nil {
		return model.AuthResponse{Success: false, Message: "Submit OTP", SessionId: &sessionId}
	}

	err = c.RedisClient.SetEX(ctx, fmt.Sprintf("sessions:%s:%s", user.Uid, sessionId), "", time.Second*30).Err()
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	authUser := model.AuthUser{
		Id:           user.Uid.String(),
		Hash:         c.AuthService.HashUser(user.Email),
		Email:        user.Email,
		Type:         "F&F", // TODO: yea
		Instances:    len(sessions) + 1,
		MaxInstances: maxInstances,
	}

	return model.AuthResponse{Success: true, User: &authUser, SessionId: &sessionId}
}

func (c *AuthController) PostSession(sessionRequest model.ActivateSessionRequest) model.AuthResponse {
	ctx := c.Ctx.Request().Context()

	session, user, err := c.AuthService.GetUserAndSession(ctx, sessionRequest.SessionId)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	if security.HwidVariance(session.Hwid, sessionRequest.Hwid) >= 2 {
		// TODO: invalidate session and log?
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	_, err = c.AuthService.GetOtpForUser(ctx, user.Uid)
	if err == nil && !session.Verified {
		return model.AuthResponse{Success: false, Message: "Submit OTP"}
	}

	acctLock, err := c.AuthService.LockAccount(ctx, user.Uid, time.Second*5)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}
	defer acctLock.Release(ctx)

	active, err := c.AuthService.IsSessionActive(ctx, user.Uid, session.Sid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	} else if active {
		return model.AuthResponse{Success: false, Message: "Session Already Active"}
	}

	maxInstances, err := c.AuthService.MaxInstancesForUser(ctx, user.Uid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	sessions, err := c.AuthService.GetActiveSessions(ctx, user.Uid)
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	} else if len(sessions) >= maxInstances {
		return model.AuthResponse{Success: false, Message: "Too Many Active Sessions"}
	}

	err = c.RedisClient.SetEX(ctx, fmt.Sprintf("sessions:%s:%s", user.Uid, session.Sid), "", time.Second*30).Err()
	if err != nil {
		return model.AuthResponse{Success: false, Message: "Internal Server Error"}
	}

	authUser := model.AuthUser{
		Id:           user.Uid.String(),
		Hash:         c.AuthService.HashUser(user.Email),
		Email:        user.Email,
		Type:         "F&F", // TODO: yea
		Instances:    len(sessions) + 1,
		MaxInstances: maxInstances,
	}

	return model.AuthResponse{Success: true, User: &authUser}
}

func (c *AuthController) PostDeactivate(sessionRequest model.DeactivateSessionRequest) model.AuthResponse {
	ctx := c.Ctx.Request().Context()

	session, user, err := c.AuthService.GetUserAndSession(ctx, sessionRequest.SessionId)
	if err != nil {
		return model.AuthResponse{Success: false}
	}

	/*
		active, err := c.Service.IsSessionActive(ctx, user.Uid, session.Sid)
		if err != nil {
			return nil, err
		} else if active {
			return nil, errors.New("session already active")
		}
	*/

	err = c.RedisClient.Del(ctx, fmt.Sprintf("sessions:%s:%s", user.Uid, session.Sid)).Err()
	if err != nil {
		return model.AuthResponse{Success: false}
	}

	return model.AuthResponse{Success: true}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmail(email string) bool {
	if len(email) < 5 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

/*
	POST /api/auth/reset
	fields: email, recaptcha_response

	This is the initial request for resetting a user's password.
	This handles the creation of the password reset request
	but doesn't handle changing the actual password.
*/
func (c *AuthController) PostReset(resetReq model.PasswordResetRequest) model.PasswordResetResponse {
	ctx := c.Ctx.Request().Context()
	email := resetReq.Email
	recaptchaResponse := resetReq.RecaptchaResponse

	if !isValidEmail(email) {
		return model.PasswordResetResponse{Success: false}
	}

	// TODO: Check recaptcha
	if strings.Compare(recaptchaResponse, "abc") != 0 {
		return model.PasswordResetResponse{Success: false}
	}

	// Return a successful message even if the email isn't found.
	// This is to prevent enumeration attacks.
	user, err := c.AuthService.GetUserForEmail(ctx, email)
	if err != nil {
		return model.PasswordResetResponse{Success: true}
	}

	// When a password reset is requested, all previous reset tokens
	// for that user should be invalidated and a new one should be
	// created.
	if err = c.AuthService.InvalidatePreviousResetTokens(ctx, user.Uid); err != nil {
		return model.PasswordResetResponse{Success: false}
	} else if err = c.AuthService.CreateAndSendPasswordReset(ctx, user.Uid, user.Email); err != nil {
		return model.PasswordResetResponse{Success: false}
	}

	return model.PasswordResetResponse{Success: true}
}

/*
	POST /api/auth/change
	fields: token, password

	This handles the changing of the user's password.
	The token is created in AuthController#PostReset
	and is provided in a link that the user will receive
	in their email.
*/
func (c *AuthController) PostChange(changeReq model.PasswordChangeRequest) model.PasswordResetResponse {
	ctx := c.Ctx.Request().Context()
	token := changeReq.Token
	password := changeReq.Password

	if len(token) > 64 || len(password) > 128 {
		return model.PasswordResetResponse{Success: false}
	}

	valid := c.AuthService.IsValidResetToken(ctx, token)
	if !valid {
		return model.PasswordResetResponse{Success: false}
	}

	uid, err := c.AuthService.GetUserIdFromResetToken(ctx, token)
	if err != nil {
		return model.PasswordResetResponse{Success: false}
	}

	if err = c.AuthService.UpdateUserPassword(ctx, uid, password); err != nil {
		return model.PasswordResetResponse{Success: false}
	} else if err = c.AuthService.InvalidatePreviousResetTokens(ctx, uid); err != nil {
		return model.PasswordResetResponse{Success: false}
	} else if err = c.AuthService.ClearSessions(ctx, uid); err != nil {
		return model.PasswordResetResponse{Success: false}
	}

	return model.PasswordResetResponse{Success: true}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *AuthController) PostOtpEnable(otpReq model.OtpEnableRequest) model.OtpEnableResponse {
	ctx := c.Ctx.Request().Context()

	_, u, err := c.AuthService.GetUserAndSession(ctx, otpReq.SessionId)
	if err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	_, err = c.AuthService.GetOtpForUser(ctx, u.Uid)
	if err == nil {
		return model.OtpEnableResponse{Success: false}
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NebulaBots",
		AccountName: u.Email,
	})
	if err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	rawRecovery, err := GenerateRandomBytes(80)
	if err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	recoveryCodes := make([]string, 0, 15)
	for i := 0; i < len(rawRecovery)-5; i += 5 {
		pt1 := hex.EncodeToString(rawRecovery[i : i+3])
		pt2 := hex.EncodeToString(rawRecovery[i+3 : i+5])
		recoveryCodes = append(recoveryCodes, fmt.Sprintf("%s-%s%s", pt1[0:5], string(pt1[5]), pt2))
	}

	img, err := key.Image(256, 256)
	if err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	var imgBuf bytes.Buffer
	if err = png.Encode(&imgBuf, img); err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	if err = c.AuthService.ClearSessions(ctx, u.Uid, otpReq.SessionId); err != nil {
		return model.OtpEnableResponse{Success: false}
	} else if err = c.AuthService.VerifySession(ctx, otpReq.SessionId); err != nil {
		return model.OtpEnableResponse{Success: false}
	} else if err = c.AuthService.SetOtpForUser(ctx, u.Uid, key.Secret(), rawRecovery); err != nil {
		return model.OtpEnableResponse{Success: false}
	}

	return model.OtpEnableResponse{Success: true, QrCode: "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBuf.Bytes()), RecoveryCodes: recoveryCodes}
}

func (c *AuthController) PostOtpVerify(otpReq model.OtpVerifyRequest) model.OtpVerifyResponse {
	ctx := c.Ctx.Request().Context()

	_, u, err := c.AuthService.GetUserAndSession(ctx, otpReq.SessionId)
	if err != nil {
		return model.OtpVerifyResponse{Success: false}
	}

	otpData, err := c.AuthService.GetOtpForUser(ctx, u.Uid)
	if err != nil {
		return model.OtpVerifyResponse{Success: false}
	}

	isValid := totp.Validate(otpReq.Code, otpData.Secret)
	if isValid {
		if err = c.AuthService.VerifySession(ctx, otpReq.SessionId); err != nil {
			return model.OtpVerifyResponse{Success: false}
		} else {
			return model.OtpVerifyResponse{Success: true}
		}
	} else {
		return model.OtpVerifyResponse{Success: false}
	}
}

func verifyKey(key, auth string) (TLDashResponse, error) {
	resp, err := client.R().
		SetHeader("Authorization", "Basic "+auth).
		SetQueryParam("licenseKey", key).
		Get("https://api.tldash.io/v1/user")

	if err != nil {
		return TLDashResponse{}, err
	}

	if resp.StatusCode() != 200 {
		return TLDashResponse{}, errors.New("status code invalid")
	}

	var dashResponse TLDashResponse
	if err = json.Unmarshal(resp.Body(), &dashResponse); err != nil {
		return TLDashResponse{}, err
	}

	return dashResponse, nil
}

func (c *AuthController) PostRegister(registerReq model.RegisterRequest) model.RegisterResponse {
	ctx := c.Ctx.Request().Context()
	email := registerReq.Email
	recaptchaResponse := registerReq.RecaptchaResponse

	a := c.AuthService.TLDashAuthorization()

	resp, err := verifyKey(registerReq.Key, a)

	if err != nil {
		return model.RegisterResponse{Success: false}
	}

	if !isValidEmail(email) || resp.Email != email {
		return model.RegisterResponse{Success: false}
	}

	// TODO: Check recaptcha
	if strings.Compare(recaptchaResponse, "abc") != 0 {
		return model.RegisterResponse{Success: false}
	}

	// Return a successful message even if the email isn't found.
	// This is to prevent enumeration attacks.
	_, err = c.AuthService.GetUserForEmail(ctx, email)
	if err == nil {
		return model.RegisterResponse{Success: true}
	}

	if err = c.AuthService.RegisterAndSendConfirmation(ctx, registerReq.Email, registerReq.Password, registerReq.Key); err != nil {
		return model.RegisterResponse{Success: false}
	}

	return model.RegisterResponse{Success: true}
}

func (c *AuthController) GetRegister(completeReq model.CompleteRegistrationRequest) model.CompleteRegistrationResponse {
	if err := c.AuthService.CompleteRegistration(c.Ctx.Request().Context(), completeReq.Id); err != nil {
		return model.CompleteRegistrationResponse{Success: false}
	} else {
		return model.CompleteRegistrationResponse{Success: true}
	}
}
