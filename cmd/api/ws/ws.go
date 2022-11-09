package ws

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"nebula/cmd/api/service"
	"strings"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

func sessionIdGenerator(ctx iris.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 {
		return "undefined"
	}
	return authParts[1]
}

func Configure(app *iris.Application, authService service.AuthService) *neffos.Server {
	ws := neffos.New(websocket.DefaultGorillaUpgrader, neffos.Namespaces{
		"nebula": neffos.Events{
			"authPing": func(conn *neffos.NSConn, message neffos.Message) error {
				ctx := context.Background()

				sid, err := uuid.Parse(conn.Conn.ID())
				if err != nil {
					return err
				}

				// TODO: stores sid and uid in conn
				session, user, err := authService.GetUserAndSession(ctx, sid)
				if err != nil {
					return err
				}

				if active, err := authService.IsSessionActive(ctx, user.Uid, session.Sid); !active {
					conn.Conn.Close()
				} else if err != nil {
					return err
				}

				maxInstances, err := authService.MaxInstancesForUser(ctx, user.Uid)
				if err != nil {
					return err
				}

				sessions, err := authService.GetActiveSessions(ctx, user.Uid)
				if err != nil {
					return err
				} else if len(sessions) > maxInstances {
					conn.Conn.Close()
				}

				return authService.RefreshSession(ctx, user.Uid, session.Sid)
			},
		},
	})

	ws.OnConnect = func(c *neffos.Conn) error {
		ctx := context.Background()

		req := c.Socket().Request()
		authHeader := req.Header.Get("Authorization")
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			return ErrNotAuthenticated
		}

		authType := authParts[0]
		authToken := authParts[1]
		if authType != "Nebula" {
			return ErrNotAuthenticated
		}

		sessionId, err := uuid.Parse(authToken)
		if err != nil {
			return ErrNotAuthenticated
		}

		session, user, err := authService.GetUserAndSession(ctx, sessionId)
		if err != nil {
			return ErrNotAuthenticated
		}

		active, err := authService.IsSessionActive(ctx, user.Uid, session.Sid)
		if err != nil || !active {
			return ErrNotAuthenticated
		}

		sessions, err := authService.GetActiveSessions(ctx, user.Uid)
		if err != nil {
			return ErrNotAuthenticated
		} else if len(sessions) > 1 {
			return ErrNotAuthenticated
		}

		if err = authService.RefreshSession(ctx, user.Uid, session.Sid); err != nil {
			return ErrNotAuthenticated
		}

		return nil
	}

	app.Get("/ws", websocket.Handler(ws, sessionIdGenerator))

	return ws
}
