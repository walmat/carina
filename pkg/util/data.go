package util

import (
	"os"
	"path"
)

func GetDataDirectory() (string, error) {
	if home, err := os.UserHomeDir(); err != nil {
		return "", err
	} else {
		return path.Join(home, ".nebula"), nil
	}
}
