package app

import (
	"embed"
	"errors"
	"fmt"
	"mime"
	"nebula/cmd/bot/inject"
	"nebula/third_party/themida"
	"net/http"
	"strings"

	"net"
)

//go:embed web/*
var fs embed.FS

var (
	listener net.Listener = nil
)

// GetUiPort binds to the frontend listener and serves the frontend port or serves the default development port if isDev is true.
func GetUiPort(devPort int) (int, error) {
	themida.Macro(themida.EAGLE_BLACK_START)

	if inject.IsDev == "yes" {
		return devPort, nil
	}

	if listener != nil {
		return listener.Addr().(*net.TCPAddr).Port, nil
	}

	uiListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, errors.New("could not bind frontend listener")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		if file, err := fs.ReadFile(fmt.Sprintf("web%s", path)); err == nil {
			splitPath := strings.Split(path, ".")
			mimeType := mime.TypeByExtension(fmt.Sprintf(".%s", splitPath[len(splitPath)-1]))
			if mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			_, err = w.Write(file)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	themida.Macro(themida.EAGLE_BLACK_END)

	go func() {
		if err = http.Serve(uiListener, nil); err != nil {
			panic(err)
		}
	}()

	listener = uiListener
	return listener.Addr().(*net.TCPAddr).Port, nil
}
