package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/joho/godotenv"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func isGAE() bool {
	_, set := os.LookupEnv("GAE_INSTANCE")
	return set
}

func expectEnv(name string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("expected env variable %s", name))
	}
	return v
}

func envOrDefault(name, def string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return def
	}
	return v
}

func main() {
	godotenv.Load()
	log.SetFlags(log.Lshortfile)

	serverConfig := Config{Debug: !isGAE()}
	if err := fromSecrets(&serverConfig); err != nil {
		log.Fatalln("failed to create config from secrets:", err)
	}

	// TODO: reconsider rsa impl
	/*
		if !serverConfig.Debug {
			privateKeyFile := os.Getenv("KEY_FILE")
			if privateKeyFile == "" {
				privateKeyFile = "config/rsa/private.pem"
			}

			privateKey, err := getRsaKey(privateKeyFile)
			if err != nil {
				log.Fatalln("could not get rsa key:", err)
			}

			serverConfig.PrivateKey = privateKey
		}
	*/

	app, err := CreateApp(serverConfig)
	if err != nil {
		log.Fatalln("failed to create app:", err)
	}

	app.Logger().SetLevel(envOrDefault("LOG_LEVEL", "info"))

	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = fmt.Sprintf(":%s", os.Getenv("PORT"))
	}

	if err = app.Listen(port); err != nil {
		log.Fatalln("server error:", err)
	}
}

type SecretDatabaseInfo struct {
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Name     string `json:"name"`
	Instance string `json:"inst"`
}

func fromSecrets(config *Config) error {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	config.RedisAddr = expectEnv("REDIS_ADDR")

	dbRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: expectEnv("DB_SECRET"),
	}
	result, err := client.AccessSecretVersion(ctx, dbRequest)
	if err != nil {
		return err
	}

	var dbInfo SecretDatabaseInfo
	if err = json.Unmarshal(result.Payload.Data, &dbInfo); err != nil {
		return err
	}
	config.DbUser = dbInfo.User
	config.DbPass = dbInfo.Pass
	config.DbName = dbInfo.Name
	if !config.Debug {
		config.DbHost = dbInfo.Instance
	} else {
		config.DbHost = expectEnv("DB_HOST")
		config.DbPort = expectEnv("DB_PORT")
	}

	sendgridRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: expectEnv("SENDGRID_SECRET"),
	}
	result, err = client.AccessSecretVersion(ctx, sendgridRequest)
	if err != nil {
		return err
	}
	config.SendgridKey = string(result.Payload.Data)

	intercomRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: expectEnv("INTERCOM_SECRET"),
	}

	result, err = client.AccessSecretVersion(ctx, intercomRequest)
	if err != nil {
		return err
	}

	config.IntercomSecret = result.Payload.Data

	tldashUserRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: expectEnv("TLDASH_USER"),
	}

	result, err = client.AccessSecretVersion(ctx, tldashUserRequest)
	if err != nil {
		return err
	}

	config.TLDashUser = result.Payload.Data

	tldashPassRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: expectEnv("TLDASH_PASS"),
	}

	result, err = client.AccessSecretVersion(ctx, tldashPassRequest)
	if err != nil {
		return err
	}

	config.TLDashPass = result.Payload.Data

	return nil
}

func getRsaKey(keyFile string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid key type, need rsa private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
