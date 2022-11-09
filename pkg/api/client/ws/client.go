package ws

import (
	"context"
	"fmt"
	"nebula/pkg/api/model"
	"nebula/pkg/clientutil"
	"nebula/pkg/infra/update"
	stdhttp "net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	gorilla "github.com/gorilla/websocket"
	"github.com/kataras/neffos"
	ng "github.com/kataras/neffos/gorilla"
)

var (
	nebulaNsConn *neffos.NSConn
)

func Connect(sessionId uuid.UUID) error {
	ws, err := neffos.Dial(context.Background(), ng.Dialer(gorilla.DefaultDialer, stdhttp.Header{"Authorization": []string{fmt.Sprintf("%s %s", model.WebsocketAuthType, sessionId)}}), "ws://localhost:8080/ws", neffos.Namespaces{
		"nebula": neffos.Events{
			"authPing": func(c *neffos.NSConn, msg neffos.Message) error {
				c.Emit("authPing", nil)
				return nil
			},
			"instanceCount": func(c *neffos.NSConn, msg neffos.Message) error {
				instanceCount, err := strconv.Atoi(string(msg.Body))
				if err != nil {
					return err
				}

				clientutil.ModifyUser(func(user *model.AuthUser) {
					user.Instances = instanceCount
				})
				return nil
			},
			"autoupdate": func(conn *neffos.NSConn, msg neffos.Message) error {
				updateInfo := strings.Split(string(msg.Body), ":")
				_, platform, _ := updateInfo[0], updateInfo[1], updateInfo[2]

				if platform == runtime.GOOS {
					update.NotifyFrontend()
				}

				return nil
			},
		},
	})
	if err != nil {
		return err
	}

	nebulaNsConn, err = ws.Connect(context.Background(), "nebula")
	if err != nil {
		return err
	}

	go reconnectHandler(sessionId, ws.NotifyClose)

	return nil
}

func reconnectHandler(sid uuid.UUID, notifyClose <-chan struct{}) {
	<-notifyClose
	for i := 0; i < 4; i++ {
		if err := Connect(sid); err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}
}
