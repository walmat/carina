package handlers

import (
	"cloud.google.com/go/logging"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	"nebula/cmd/api/service"
	"nebula/pkg/api/model"
	"strings"
)

const logName = "tasks"

func LoggerHandler(client *logging.Client, authService service.AuthService) iris.Handler {
	logger := client.Logger(logName)
	return func(ctx iris.Context) {
		authHeader := ctx.GetHeader("Authorization")
		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			ctx.StatusCode(iris.StatusUnauthorized)
			return
		}

		sessionId, err := uuid.Parse(authParts[1])
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		s, u, err := authService.GetUserAndSession(ctx.Request().Context(), sessionId)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		var logEntry model.LoggerEntry
		if err = ctx.ReadJSON(&logEntry); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		labels := make(map[string]string)
		if logEntry.Labels != nil {
			labels = logEntry.Labels
		}
		labels["sessionId"] = s.Sid.String()
		labels["userId"] = u.Uid.String()

		entry := logging.Entry{
			Severity: logging.Severity(logEntry.Severity),
			Payload:  logEntry.Payload,
			Labels:   labels,
		}
		if logEntry.File != nil {
			entry.SourceLocation = &loggingpb.LogEntrySourceLocation{
				File:     logEntry.File.Name,
				Line:     int64(logEntry.File.Line),
			}
		}

		logger.Log(entry)
	}
}
