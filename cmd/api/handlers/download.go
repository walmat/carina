package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/kataras/iris/v12"
	ua "github.com/mileusna/useragent"
	"google.golang.org/api/iterator"
)

type storageConnection struct {
	Client *storage.Client
}

var (
	client *storageConnection
	once   sync.Once
)

const (
	// GCSBucket name
	GCSBucket = "carina-autoupdater-1"
	// ProjectID Google Project ID name
	ProjectID = "Nebulabots"
)

// GetGCSClient gets singleton object for Google Storage
func GetGCSClient(ctx context.Context) (*storage.Client, error) {
	var clientErr error
	once.Do(func() {
		storageClient, err := storage.NewClient(ctx)
		if err != nil {
			clientErr = fmt.Errorf("failed to create GCS client ERROR:%s", err.Error())
		} else {
			client = &storageConnection{
				Client: storageClient,
			}
		}
	})
	return client.Client, clientErr
}

// Download gets a file from GCS bucket, Takes file path as a path param from request
func DownloadHandler(ctx iris.Context) {
	clientCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filename string
	useragent := ctx.GetHeader("user-agent")
	if useragent == "" {
		return
	}

	client, err := GetGCSClient(clientCtx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	u := ua.Parse(useragent)

	if u.OS != "Windows" && u.OS != "macOS" {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	it := client.Bucket(GCSBucket).UserProject(ProjectID).Objects(clientCtx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		n := attrs.Name
		// check if its an exe or a dmg
		if strings.Contains(n, "exe") && !strings.Contains(n, "exe.bak") && u.OS == "Windows" {
			filename = n
		}

		if strings.Contains(n, "dmg") && u.OS == "macOS" {
			filename = n
		}
	}

	reader, err := client.Bucket(GCSBucket).UserProject(ProjectID).Object(filename).NewReader(clientCtx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	defer reader.Close()

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	p := strings.Split(filename, "-")
	name := p[0]

	q := strings.Split(filename, ".")
	ext := q[1]

	formattedName := name + "." + ext

	ctx.Header("Content-Type", reader.Attrs.ContentType)
	ctx.Header("Content-Disposition", "attachment"+"; filename="+formattedName)
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.StatusCode(iris.StatusOK)

	ctx.Write(content)
}

func CheckForUpdateHandler(ctx iris.Context) {
	clientCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	platform := ctx.URLParam("platform")
	if platform == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	client, err := GetGCSClient(clientCtx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	var latest string
	it := client.Bucket(GCSBucket).UserProject(ProjectID).Objects(clientCtx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		n := attrs.Name

		if strings.Contains(n, "exe") && platform == "windows" {
			latest = n
		}

		if strings.Contains(n, "dmg") && platform == "darwin" {
			latest = n
		}
	}

	ctx.JSON(map[string]interface{}{"latest": latest})
	ctx.StatusCode(iris.StatusOK)
}
