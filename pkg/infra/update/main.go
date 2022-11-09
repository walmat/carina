package update

import (
	"bytes"
	"errors"
	"fmt"
	"nebula/cmd/bot/inject"
	"nebula/cmd/bot/windows"
	"nebula/pkg/api/client/rest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/inconshreveable/go-update"
)

func NotifyFrontend() {
	if windows.Windows.MainWindow != nil {
		windows.Windows.MainWindow.ExecuteJavaScript(fmt.Sprintf(`
		let evt = new CustomEvent("%s", {detail: "%s"});
		window.dispatchEvent(evt);
		`, "autoupdater", "A new version of Carina is available."))
		return
	}
}

func CheckForUpdates() (bool, error) {
	res, err := rest.CheckForUpdates()

	if err == nil && res.Latest != "" {
		version := inject.CommitHash

		var ext string
		switch runtime.GOOS {
		case "darwin":
			ext = ".dmg"
		case "windows":
			ext = ".exe"
		}

		re := regexp.MustCompile(fmt.Sprintf(`-(.*)%s`, ext))
		match := re.FindStringSubmatch(res.Latest)
		if len(match) > 1 {
			if version != "dev" && match[1] != version {
				return true, nil
			}
		} else {
			err = errors.New("error parsing version from exec name")
			return false, err
		}
	}

	return false, nil
}

// Update initiates the self update process
func Update() error {
	windows.Windows.MainWindow.Hide()
	var readerBytes []byte

	resp, err := rest.Download()
	if err != nil {
		return err
	}

	readerBytes = resp

	if runtime.GOOS == "darwin" {
		ex, err := os.Executable()
		if err != nil {
			return err
		}
		exPath := filepath.Dir(ex)

		t := filepath.Join(exPath, "Carina.dmg")
		err = os.WriteFile(t, resp, 0755)
		if err != nil {
			return err
		}

		a, err := exec.Command("hdiutil", "attach", t, "-nobrowse").Output()
		if err != nil {
			return err
		}

		defer func() {
			//todo: maybe this doesnt work. does it matter if we delete dmg anyway?
			exec.Command("hdiutil", "detach", string(a)).Run()
			os.Remove(t)
			windows.App.Close()
		}()

		// todo: may need to get most recent dmg if they have stray dmgs
		b, err := os.ReadFile("/Volumes/Carina/Carina.app/Contents/MacOS/Carina")
		if err != nil {
			return err
		}

		readerBytes = b
	}

	return update.Apply(bytes.NewReader(readerBytes), update.Options{})
}
