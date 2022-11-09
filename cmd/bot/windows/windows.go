package windows

import "github.com/Nebulabots/go-astilectron"

var App *astilectron.Astilectron

var Windows = struct {
	LoginWindow      *astilectron.Window
	MainWindow       *astilectron.Window
	CollectiveWindow *astilectron.Window
}{}
