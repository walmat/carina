package app

import (
	"nebula/cmd/bot/inject"

	"github.com/Nebulabots/go-astilectron"
	"github.com/asticode/go-astikit"
)

func CreateWindow(title string, url string, width, height int, resizable bool) (*astilectron.Window, error) {
	w, err := App.NewWindow(url, &astilectron.WindowOptions{
		Center:           astikit.BoolPtr(true),
		Title:            astikit.StrPtr(title),
		Height:           astikit.IntPtr(height),
		MinHeight:        astikit.IntPtr(height),
		MinWidth:         astikit.IntPtr(width),
		Width:            astikit.IntPtr(width),
		Resizable:        astikit.BoolPtr(resizable),
		Frame:            astikit.BoolPtr(false),
		Show:             astikit.BoolPtr(false),
		Transparent:      astikit.BoolPtr(false),
		AcceptFirstMouse: astikit.BoolPtr(true),
		WebPreferences: &astilectron.WebPreferences{
			AllowRunningInsecureContent: astikit.BoolPtr(false),
			BackgroundThrottling:        astikit.BoolPtr(true),
			ContextIsolation:            astikit.BoolPtr(false),
			DevTools:                    astikit.BoolPtr(inject.IsDev == "yes"),
			Plugins:                     astikit.BoolPtr(true),
		},
	})
	if err != nil {
		return nil, err
	}

	w.OnMessage((&eventHandler{w}).EventHandler)

	w.On(astilectron.EventNameWindowEventReadyToShow, func(e astilectron.Event) bool {
		if title != "Collective" {
			_ = w.Show()
		}
		//
		//if inject.IsDev == "yes" {
		//	_ = w.OpenDevTools()
		//}
		return true
	})

	w.On(astilectron.EventNameWindowCmdClose, func(e astilectron.Event) bool {
		if inject.IsDev == "yes" {
			_ = w.CloseDevTools()
		}

		return false
	})

	_ = w.Create()
	_ = CreateMenu(w)

	return w, nil
}
