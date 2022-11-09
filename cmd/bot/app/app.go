package app

import (
	"fmt"
	"nebula/cmd/bot/windows"
	"nebula/pkg/logger"
	"time"

	"github.com/Nebulabots/go-astilectron"
)

var App *astilectron.Astilectron

func TransitionToMain() {
	if windows.Windows.LoginWindow != nil {
		windows.Windows.LoginWindow.Hide()
		MainWindow()
		// wait to make sure callback for this is called before we completely close window
		go time.AfterFunc(time.Second*5, func() {
			windows.Windows.LoginWindow.Close()
		})
	}
}

func TransitionToLogin() {
	if windows.Windows.MainWindow != nil {
		windows.Windows.MainWindow.Hide()
		AuthWindow()
		// wait to make sure callback for this is called before we completely close window
		go time.AfterFunc(time.Second*5, func() {
			windows.Windows.MainWindow.Close()
		})
	}
}

func MainWindow() {
	port, err := GetUiPort(3000)
	if err != nil {
		logger.Error("failed to get main URL", err)
	}

	mainUrl := fmt.Sprintf("http://127.0.0.1:%d/index.html", port)
	w, err := CreateWindow("Carina", mainUrl, 1150, 765, true)
	if err != nil {
		logger.Error("error launching main window:", err)
	}

	windows.Windows.MainWindow = w

	go CollectiveWindow()
}

func AuthWindow() {
	port, err := GetUiPort(3001)
	if err != nil {
		logger.Error("failed to get main URL", err)
	}

	loginUrl := fmt.Sprintf("http://127.0.0.1:%d/login.html", port)
	w, err := CreateWindow("Carina", loginUrl, 925, 630, false)
	if err != nil {
		logger.Error("error launching auth window:", err)
	}

	windows.Windows.LoginWindow = w
}

func CollectiveWindow() {
	port, err := GetUiPort(3002)
	if err != nil {
		logger.Error("failed to get main URL", err)
	}

	harvesterUrl := fmt.Sprintf("http://127.0.0.1:%d/collective.html", port)
	w, err := CreateWindow("Collective", harvesterUrl, 410, 640, true)
	if err != nil {
		logger.Error("error launching collective window:", err)
	}

	windows.Windows.CollectiveWindow = w
}
