package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/akavel/rsrc/rsrc"
)

var (
	GOOS         string
	GOARCH       string
	ldflags      []string
	appName      string
	outputDir    string
	binaryDir    string
	resourcesDir string

	binaryName string
	version    string

	errUnsupportedOS = errors.New("unsupported OS")
)

func init() {
	appName = "Carina"
	flag.StringVar(&GOOS, "GOOS", runtime.GOOS, "Operating system to target")
	flag.StringVar(&GOARCH, "GOARCH", runtime.GOARCH, "Architecture to target")
	flag.StringVar(&outputDir, "outputDir", "dist", "path to output for the finished application")
	flag.StringVar(&binaryDir, "binaryDir", "dist/bin", "path to output for the binary")
	flag.StringVar(&resourcesDir, "resourcesDir", "resources", "path to resources(icons, plist, syso files etc.)")

	if err := resetDir(outputDir); err != nil {
		panic(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	commitHash, err := getCommitHash(ctx)

	version = strings.TrimSuffix(string(commitHash), "\n")

	if err != nil {
		panic(err)
	}

	switch GOOS {
	case "windows":
		appName = appName + "-" + version

		binaryName = appName + ".exe"
		ldflags = append(ldflags, "-H=windowsgui")

		createWindowsSyso()
	case "darwin":
		binaryName = appName
		if err := createDarwinDir(); err != nil {
			panic(err)
		}
	default:
		panic(errUnsupportedOS)
	}
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancelFunc()
	}()

	if err := setCommitHash(ctx); err != nil {
		panic(err)
	}

	if os.Getenv("NEBULA_ENABLE_RSA") != "" {
		if err := handleRSA(); err != nil {
			panic(err)
		}
	}

	ldflags = append(ldflags, "-X nebula/pkg/api/client.ApiUrl=https://nebulabots.uc.r.appspot.com")
	ld := strings.Join(ldflags, " ")
	var binaryPath = filepath.Join(binaryDir, binaryName)
	args := []string{"build", "-o", binaryPath, "-ldflags", ld, "nebula/cmd/bot"}
	cmd := exec.CommandContext(ctx, "go", args...)

	baseEnv := os.Environ()
	baseEnv = append(baseEnv, fmt.Sprintf("GOOS=%s", GOOS), fmt.Sprintf("GOARCH=%s", GOARCH), "CGO_ENABLED=1")
	cmd.Env = baseEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	finish(ctx)
}

func finish(ctx context.Context) {
	var binaryAppPath string
	var binaryPath = filepath.Join(binaryDir, binaryName)
	switch GOOS {
	case "darwin":
		binaryAppPath = filepath.Join(outputDir, appName+".app", "Contents", "MacOS", binaryName)

		copyFile(binaryPath, binaryAppPath)

		cmd := exec.CommandContext(ctx, "chmod", "755", binaryAppPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		gon := Config{}
		gon.Source = []string{filepath.Join(outputDir, appName+".app")}
		gon.BundleId = "com.nebulabots.carina"

		gon.Sign = &Sign{}
		gon.Sign.ApplicationIdentity = "Developer ID Application: Matt Wall"

		gon.AppleId = &AppleId{}
		gon.AppleId.Password = "@env:AC_PASSWORD"

		gon.Dmg = &Dmg{}
		gon.Dmg.OutputPath = filepath.Join(outputDir, binaryName+"-"+version+".dmg")
		gon.Dmg.VolumeName = appName

		o, err := json.MarshalIndent(gon, "", " ")

		if err != nil {
			panic(err)
		}

		os.WriteFile(filepath.Join("scripts", "build-scripts", "gon.json"), o, 0644)

	case "windows":
		binaryAppPath = filepath.Join(outputDir, binaryName)

		copyFile(binaryPath, binaryAppPath)

		cmd := exec.CommandContext(ctx, "./scripts/Themida/Themida64.exe", "/protect", "./scripts/Themida/Carina.tmd", "/inputfile", binaryAppPath, "/outputfile", binaryAppPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	os.RemoveAll(binaryDir)
}

func createWindowsSyso() (err error) {
	var iconPath = filepath.Join(resourcesDir, "icons", "icon.ico")

	var sysoOutputPath = filepath.Join("cmd", "bot", "windows.syso")

	if err = rsrc.Embed(sysoOutputPath, "amd64", "", iconPath); err != nil {
		return
	}

	return
}

func createDarwinDir() (err error) {
	var contentsPath = filepath.Join(outputDir, appName+".app", "Contents")
	var macOSPath = filepath.Join(contentsPath, "MacOS")
	if err = os.MkdirAll(macOSPath, 0777); err != nil {
		return
	}

	var resourcesPath = filepath.Join(contentsPath, "Resources")
	if err = os.MkdirAll(resourcesPath, 0777); err != nil {
		return
	}

	if err = copyFile("resources/icons/icon.icns", resourcesPath+"/icon.icns"); err != nil {
		return
	}

	if err = copyFile("resources/Info.plist", contentsPath+"/Info.plist"); err != nil {
		return
	}

	return
}

func handleRSA() (err error) {
	rsaKeyFile, err := os.ReadFile("config/rsa/public.pem")
	if err != nil {
		return
	}

	pemBytes, err := json.Marshal(rsaKeyFile)
	if err != nil {
		return
	}

	ldflags = append(ldflags, fmt.Sprintf("-X nebula/cmd/bot/inject.RsaPem=%s", string(pemBytes)))
	return
}

func copyFile(originalPath string, newPath string) (err error) {
	sourceFile, err := os.Open(originalPath)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	newFile, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		return
	}

	return
}

func resetDir(p string) (err error) {
	// Remove
	os.RemoveAll(p)

	// Mkdir
	if err = os.MkdirAll(p, 0755); err != nil {
		return
	}
	return
}

func getCommitHash(ctx context.Context) (out []byte, err error) {
	return exec.CommandContext(ctx, "git", "rev-parse", "--short", "HEAD").Output()
}

func setCommitHash(ctx context.Context) (err error) {
	out, err := getCommitHash(ctx)

	if err != nil {
		return
	}

	ldflags = append(ldflags, fmt.Sprintf("-s -w -X nebula/cmd/bot/inject.CommitHash=%s", string(out)))
	return
}

// Config is the configuration structure for gon.
type Config struct {
	// Source is the list of binary files to sign.
	Source []string `json:"source,omitempty"`

	// BundleId is the bundle ID to use for the package that is created.
	// This should be in a format such as "com.example.app". The value can
	// be anything, this is required by Apple.
	BundleId string `json:"bundle_id,omitempty"`

	// Notarize is a single file (usually a .pkg installer or zip)
	// that is ready for notarization as-is
	Notarize []Notarize `json:"notarize,omitempty"`

	// Sign are the settings for code-signing the binaries.
	Sign *Sign `json:"sign,omitempty"`

	// AppleId are the credentials to use to talk to Apple.
	AppleId *AppleId `json:"apple_id,omitempty"`

	// Zip, if present, creates a notarized zip file as the output. Note
	// that zip files do not support stapling, so the final result will
	// require an internet connection on first use to validate the notarization.
	Zip *Zip `json:"zip,omitempty"`

	// Dmg, if present, creates a dmg file to package the signed `Source` files
	// into. Dmg files support stapling so this allows offline usage.
	Dmg *Dmg `json:"dmg,omitempty"`
}

// AppleId are the authentication settings for Apple systems.
type AppleId struct {
	// Username is your AC username, typically an email. This is required, but will
	// be read from the environment via AC_USERNAME if not specified via config.
	Username string `json:"username,omitempty"`

	// Password is the password for your AC account. This also accepts
	// two additional forms: '@keychain:<name>' which reads the password from
	// the keychain and '@env:<name>' which reads the password from an
	// an environmental variable named <name>. If omitted, it has the same effect
	// as passing '@env:AC_PASSWORD'.
	Password string `json:"password,omitempty"`

	// Provider is the AC provider. This is optional and only needs to be
	// specified if you're using an Apple ID account that has multiple
	// teams.
	Provider string `json:"provider,omitempty"`
}

// Notarize are the options for notarizing a pre-built file.
type Notarize struct {
	// Path is the path to the file to notarize. This can be any supported
	// filetype (dmg, pkg, app, zip).
	Path string `json:"path"`

	// BundleId is the bundle ID to use for notarizing this package.
	// If this isn't specified then the root bundle_id is inherited.
	BundleId string `json:"bundle_id"`

	// Staple, if true will staple the notarization ticket to the file.
	Staple bool `json:"staple,omitempty"`
}

// Sign are the options for codesigning the binaries.
type Sign struct {
	// ApplicationIdentity is the ID or name of the certificate to
	// use for signing binaries. This is used for all binaries in "source".
	ApplicationIdentity string `json:"application_identity"`
	// Specify a path to an entitlements file in plist format
	EntitlementsFile string `json:"entitlements_file,omitempty"`
}

// Dmg are the options for a dmg file as output.
type Dmg struct {
	// OutputPath is the path where the final dmg will be saved.
	OutputPath string `json:"output_path"`

	// Volume name is the name of the volume that shows up in the title
	// and sidebar after opening it.
	VolumeName string `json:"volume_name"`
}

// Zip are the options for a zip file as output.
type Zip struct {
	// OutputPath is the path where the final zip file will be saved.
	OutputPath string `json:"output_path"`
}
