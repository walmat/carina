package handlers

import (
	// "github.com/gofiber/fiber/v2"
	"github.com/kataras/iris/v12"
	"time"
)

func HealthHandler(startTime time.Time) iris.Handler {
	return func(ctx iris.Context) {
		_, err := ctx.JSON(map[string]interface{}{"uptime": time.Now().Sub(startTime)})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}
	}
}
