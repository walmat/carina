package app

import (
	"log"
	"nebula/cmd/bot/windows"

	"github.com/Nebulabots/go-astilectron"
	"github.com/asticode/go-astikit"
)

func CreateMenu(w *astilectron.Window) error {
	return w.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Edit"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label:       astikit.StrPtr("Undo"),
					Role:        astilectron.MenuItemRoleUndo,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+Z"),
				},
				{
					Label:       astikit.StrPtr("Redo"),
					Role:        astilectron.MenuItemRoleRedo,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+Y"),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label:       astikit.StrPtr("Cut"),
					Role:        astilectron.MenuItemRoleCut,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+X"),
				},
				{
					Label:       astikit.StrPtr("Copy"),
					Role:        astilectron.MenuItemRoleCopy,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+C"),
				},
				{
					Label:       astikit.StrPtr("Paste"),
					Role:        astilectron.MenuItemRolePaste,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+V"),
				},
				{
					Label:       astikit.StrPtr("Select All"),
					Role:        astilectron.MenuItemRoleSelectAll,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+A"),
				},
			},
		},
		{
			Label: astikit.StrPtr("View"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label:       astikit.StrPtr("Reload"),
					Role:        astilectron.MenuItemRoleReload,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+R"),
				},
				{
					Label:       astikit.StrPtr("Full Screen"),
					Role:        astilectron.MenuItemRoleToggleFullScreen,
					Accelerator: astilectron.NewAccelerator("Control+CommandOrControl+F"),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Show Collective"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if windows.Windows.CollectiveWindow == nil {
							CollectiveWindow()
						}

						err := windows.Windows.CollectiveWindow.Show()
						if err != nil {
							return false
						}
						return false
					},
					Accelerator: astilectron.NewAccelerator("CommandOrControl+~"),
				},
			},
		},
		{
			Label: astikit.StrPtr("Window"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label:       astikit.StrPtr("Hide"),
					Role:        astilectron.MenuItemRoleReload,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+H"),
				},
				{
					Label:       astikit.StrPtr("Minimize"),
					Role:        astilectron.MenuItemRoleReload,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+M"),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label:       astikit.StrPtr("Close"),
					Role:        astilectron.MenuItemRoleReload,
					Accelerator: astilectron.NewAccelerator("CommandOrControl+Q"),
				},
				{
					Label: astikit.StrPtr("Developer Tools"),
					Role:  astilectron.MenuItemRoleToggleDevTools,
				},
			},
		},
		{
			Label: astikit.StrPtr("Help"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					Role:  astilectron.MenuItemRoleAbout,
					OnClick: func(e astilectron.Event) bool {
						go func() {

							aboutUrl := ""

							w, err := CreateWindow("About", aboutUrl, 350, 400, true)
							if err != nil {
								log.Fatal(err)
							}
							err = w.Create()
							if err != nil {
								log.Fatal(err)
							}
						}()
						return false
					},
				},
				{
					Label: astikit.StrPtr("Support"),
					Role:  astilectron.MenuItemRoleHelp,
					OnClick: func(e astilectron.Event) bool {
						// TODO: Open intercom in browser
						return false
					},
				},
			},
		},
	}).Create()
}
