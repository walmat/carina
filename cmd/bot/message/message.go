package message

import (
	"errors"
	"nebula/cmd/bot/windows"
)

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func SendMessageToMain(event string, message interface{}) error {
	if windows.Windows.MainWindow == nil {
		return errors.New("window is nil")
	}

	err := windows.Windows.MainWindow.SendMessage(Message{Event: event, Data: message})

	return err
}

func SendMessageToCollective(event string, message interface{}) error {
	if windows.Windows.CollectiveWindow == nil {
		return errors.New("window is nil")
	}

	err := windows.Windows.CollectiveWindow.SendMessage(Message{Event: event, Data: message})

	return err
}
