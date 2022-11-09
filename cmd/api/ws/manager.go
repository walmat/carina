package ws

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kataras/neffos"
	"nebula/cmd/api/service"
	"strconv"
	"strings"
	"time"
)

func HeartbeatManager(ws *neffos.Server) {
	for range time.Tick(time.Second * 10) {
		ws.Broadcast(nil, neffos.Message{
			Namespace: "nebula",
			Event:     "authPing",
		})
	}
}

func containsSession(strs []string, s string) bool {
	for _, str := range strs {
		if strings.HasSuffix(str, s) {
			return true
		}
	}
	return false
}

func AuthListener(ws *neffos.Server, redisClient *redis.Client, authService service.AuthService) {
	ctx := context.Background()
	sub := redisClient.Subscribe(context.Background(), "__keyevent@0__:set", "__keyevent@0__:del", "__keyevent@0__:expired")
	for msg := range sub.Channel() {
		payloadParts := strings.Split(msg.Payload, ":")
		if payloadParts[0] != "sessions" || len(payloadParts) != 3 {
			continue
		}

		uid, err := uuid.Parse(payloadParts[1])
		if err != nil {
			continue
		}

		sessions, err := authService.GetActiveSessions(ctx, uid)
		if err != nil {
			continue
		}

		sessionCountData := []byte(strconv.Itoa(len(sessions)))
		ws.Do(func(c *neffos.Conn) {
			if !containsSession(sessions, c.ID()) {
				return
			}
			c.Namespace("nebula").Emit("instanceCount", sessionCountData)
		}, false)
	}
}
