package app

import "github.com/ncruces/zenity"

func SelectFile(title string) (string, error) {
	return zenity.SelectFile(zenity.Title(title))
}

func SelectFileFiltered(title string, filterDesc string, filterExt ...string) (string, error) {
	return zenity.SelectFile(zenity.Title(title), zenity.FileFilter{Name: filterDesc, Patterns: filterExt})
}
