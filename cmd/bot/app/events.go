package app

import (
	"encoding/json"
	"fmt"
	"log"
	"nebula/cmd/bot/windows"
	"nebula/pkg/logger"
	"os/exec"
	"runtime"

	"github.com/Nebulabots/go-astilectron"
)

type eventMessage struct {
	Type string `json:"type"`
	Data string `json:"data"` // `JSON.stringify()`
}

type rpcData struct {
	Action string        `json:"action"`
	Input  []interface{} `json:"input"`
}

type rpcResp struct {
	Error  error       `json:"error,omitempty"`
	Output interface{} `json:"output,omitempty"`
}

type eventHandler struct {
	w *astilectron.Window
}

func (e *eventHandler) EventHandler(m *astilectron.EventMessage) interface{} {
	var ev eventMessage
	if err := m.Unmarshal(&ev); err != nil {
		return nil
	}

	switch ev.Type {
	case "showCollective":
		if windows.Windows.CollectiveWindow == nil {
			CollectiveWindow()
		}

		_ = windows.Windows.CollectiveWindow.Show()
		return nil
	case "close":
		if e.w == windows.Windows.MainWindow || e.w == windows.Windows.LoginWindow {
			_ = App.Quit()
		} else {
			_ = e.w.Close()
		}
		return nil
	case "hide":
		_ = e.w.Hide()
		return nil
	case "minimize":
		_ = e.w.Minimize()
		return nil
	case "rpc":
		var rpcCall rpcData
		if err := json.Unmarshal([]byte(ev.Data), &rpcCall); err != nil {
			return rpcResp{Error: err}
		}
		retVal, err := handleRPC(rpcCall.Action, rpcCall.Input)
		if err != nil {
			return rpcResp{Error: err}
		}

		return rpcResp{Output: retVal}
	case "openExternalUrl":
		openBrowser(ev.Data)
		return nil
	default:
		logger.Warn("app: unexpected event type:", ev.Type)
		return nil
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
