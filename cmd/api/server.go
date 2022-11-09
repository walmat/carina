package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"nebula/cmd/api/controller"
	"nebula/cmd/api/handlers"
	"nebula/cmd/api/middleware"
	"nebula/cmd/api/service"
	"nebula/cmd/api/ws"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/iris-contrib/middleware/cors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/sendgrid/sendgrid-go"
)

type Config struct {
	Debug      bool
	PrivateKey *rsa.PrivateKey

	DbUser string
	DbPass string
	DbName string
	DbHost string
	DbPort string

	RedisAddr string

	SendgridKey string

	IntercomSecret []byte

	TLDashUser []byte
	TLDashPass []byte
}

func CreateApp(config Config) (*iris.Application, error) {
	app := iris.New()
	if !config.Debug {
		app.Configure(iris.WithoutStartupLog)
	}

	app.Use(recover.New())
	app.Use(logger.New())

	// uncommnent to enable cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})
	app.UseRouter(crs)

	app.Get("/health", handlers.HealthHandler(time.Now()))
	if config.Debug {
		app.Get("/graphiql", handlers.GraphiQLHandler)
	}

	if err := SetupApiRoutes(config, app); err != nil {
		return nil, err
	}

	return app, nil
}

func SetupApiRoutes(config Config, app *iris.Application) error {
	var api iris.Party
	if config.Debug || config.PrivateKey == nil {
		api = app.Party("/api")
	} else {
		api = app.Party("/api", middleware.DecryptBody(config.PrivateKey), middleware.SignResponse(config.PrivateKey))
	}

	api.Post("/graphql", handlers.GraphQLHandler)

	api.Get("/download", handlers.DownloadHandler)

	api.Get("/update", handlers.CheckForUpdateHandler)

	pgPool, err := newPgPool(nil, config)
	if err != nil {
		return err
	}
	configurePgPool(pgPool)

	redisClient, err := newRedisClient(nil, config)
	if err != nil {
		return err
	}
	redisLock := redislock.New(redisClient)

	sendClient := newSendClient(nil, config)

	authService := service.NewAuthService(pgPool, redisClient, redisLock, sendClient, config.IntercomSecret, config.TLDashUser, config.TLDashPass)

	logClient, err := logging.NewClient(context.Background(), expectEnv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		return err
	}
	api.Post("/log", handlers.LoggerHandler(logClient, authService))

	mvc.Configure(api.Party("/auth"), func(app *mvc.Application) {
		app.Register(
			pgPool,
			redisClient,
			redisLock,
			authService,
		)

		app.Handle(new(controller.AuthController))
	})

	wsServer := ws.Configure(app, authService)
	go ws.HeartbeatManager(wsServer)
	go ws.AuthListener(wsServer, redisClient, authService)

	redisSub := redisClient.Subscribe(context.Background(), "carina:autoupdate")
	updateChan := redisSub.Channel()
	go func() {
		for update := range updateChan {
			updateInfo := strings.Split(update.Payload, ":")
			file, hash := updateInfo[0], updateInfo[1]

			platform := ""
			if strings.HasSuffix(file, ".exe") {
				platform = "windows"
			} else if strings.HasSuffix(file, ".dmg") {
				platform = "darwin"
			}

			wsServer.Broadcast(nil, neffos.Message{
				Namespace: "nebula",
				Event:     "autoupdate",
				Body:      []byte(fmt.Sprintf("%s:%s", platform, hash)),
			})
		}
	}()

	psc, err := pubsub.NewClient(context.Background(), "nebulabots")
	if err != nil {
		return err
	}

	sub := psc.Subscription("AUTOUPDATER_V1-sub")
	go func() {
		_ = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
			m.Ack()

			if m.Attributes["eventType"] == "OBJECT_FINALIZE" {
				storageEvent, err := UnmarshalStorageEvent(m.Data)
				if err != nil {
					return
				}

				if strings.HasSuffix(storageEvent.Name, "Carina.exe") || strings.HasSuffix(storageEvent.Name, "Carina.dmg") {
					redisClient.Publish(context.Background(), "carina:autoupdate", fmt.Sprintf("%s:%s:%s", storageEvent.ID, storageEvent.Name, storageEvent.Md5Hash))
				}
			}
		})
	}()

	return nil
}

func newPgPool(_ iris.Context, config Config) (*sql.DB, error) {
	var dbUri string
	if !config.Debug {
		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}

		dbUri = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", config.DbUser, config.DbPass, config.DbName, socketDir, config.DbHost)
	} else {
		dbUri = fmt.Sprintf("user=%s password=%s database=%s host=localhost port=%s", config.DbUser, config.DbPass, config.DbName, config.DbPort)
	}
	return sql.Open("pgx", dbUri)
}

func configurePgPool(db *sql.DB) {
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1800)
}

func newRedisClient(_ iris.Context, config Config) (*redis.Client, error) {
	var redisAddr string
	if redisEnv := config.RedisAddr; redisEnv != "" {
		redisAddr = redisEnv
	} else {
		redisAddr = "127.0.0.1:6379"
	}

	return redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    redisAddr,
	}), nil
}

func newSendClient(_ iris.Context, config Config) *sendgrid.Client {
	return sendgrid.NewSendClient(config.SendgridKey)
}
