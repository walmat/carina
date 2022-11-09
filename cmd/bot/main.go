package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"log"
	"math/rand"
	"nebula/cmd/bot/app"
	"nebula/cmd/bot/inject"
	"nebula/cmd/bot/windows"
	"nebula/pkg/api/client"
	"nebula/pkg/api/client/rest"
	"nebula/pkg/api/client/ws"
	"nebula/pkg/clientutil"
	"nebula/pkg/infra/accounts"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/rates"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/webhooks"
	"nebula/pkg/integrations"
	"nebula/pkg/integrations/autosolve"
	"nebula/pkg/logger"
	"nebula/pkg/security/hwid"
	_ "nebula/pkg/taskimpl"
	"nebula/pkg/userpref"
	"nebula/pkg/util"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"

	"github.com/Nebulabots/go-astilectron"
)

const (
	IconsDir = "resources/icons"

	IntegrationsFile    = "integrations.json"
	UserPreferencesFile = "preferences.json"

	TasksFile    = "tasks.neb"
	ProfilesFile = "profiles.neb"
	ProxiesFile  = "proxies.neb"
	AccountsFile = "accounts.neb"

	WebhooksFile = "webhooks.json"

	HarvestersFile = "harvesters.json"

	GmailsFile = "gmails.json"
)

func sharedMain() {
	// themida.Macro(themida.EAGLE_BLACK_START)
	rand.Seed(time.Now().UnixNano())

	dataDir, err := util.GetDataDirectory()
	if err != nil {
		panic(err)
	} else if err = os.MkdirAll(dataDir, os.ModePerm); err != nil {
		logger.Error("error creating data directory:", err)
	}

	if inject.RsaPem != "" {
		setupRsa()
	}

	integrationsFile := path.Join(dataDir, IntegrationsFile)
	if err = integrations.ImportIntegrations(integrationsFile); err != nil {
		logger.Error("failed to import integrations:", err)
	}

	preferencesFile := path.Join(dataDir, UserPreferencesFile)
	if err = userpref.ImportPreferences(preferencesFile); err != nil {
		logger.Error("failed to import user preferences:", err)
	}

	tasksFile := path.Join(dataDir, TasksFile)
	if err = tasks.ImportGob(tasksFile); err != nil {
		logger.Error("failed to import task state:", err)
	}

	profilesFile := path.Join(dataDir, ProfilesFile)
	if err = profiles.ImportGob(profilesFile); err != nil {
		logger.Error("failed to import profiles state:", err)
	}

	proxiesFile := path.Join(dataDir, ProxiesFile)
	if err = proxies.ImportGob(proxiesFile); err != nil {
		logger.Error("failed to import proxies state:", err)
	}

	accountsFile := path.Join(dataDir, AccountsFile)
	if err = accounts.ImportGob(accountsFile); err != nil {
		logger.Error("failed to import accounts state:", err)
	}

	webhooksFile := path.Join(dataDir, WebhooksFile)
	if err = webhooks.ImportState(webhooksFile); err != nil {
		logger.Error("failed to import webhooks state:", err)
	}

	harvestersFile := path.Join(dataDir, HarvestersFile)
	if err = harvesters.ImportState(harvestersFile); err != nil {
		logger.Error("failed to import harvesters state:", err)
	}

	gmailsFile := path.Join(dataDir, GmailsFile)
	if err = harvesters.ImportGmailState(gmailsFile); err != nil {
		logger.Error("failed to import harvester gmails state:", err)
	}

	if err = addDefaultGroups(harvestersFile); err != nil {
		logger.Error("failed to create default groups:", err)
		return
	}

	go periodicExport(dataDir, time.Second*15)
	defer func() {
		if err = integrations.ExportIntegrations(integrationsFile); err != nil {
			logger.Error("failed to export integrations:", err)
		}
	}()
	defer func() {
		if err = userpref.ExportPreferences(preferencesFile); err != nil {
			logger.Error("failed to export user preferences:", err)
		}
	}()
	defer func() {
		if err = tasks.ExportAllGob(tasksFile); err != nil {
			logger.Error("failed to export tasks:", err)
		}
	}()
	defer func() {
		if err = profiles.ExportAllGob(profilesFile); err != nil {
			logger.Error("failed to export profiles:", err)
		}
	}()
	defer func() {
		if err = proxies.ExportAllGob(proxiesFile); err != nil {
			logger.Error("failed to export proxies:", err)
		}
	}()
	defer func() {
		if err = accounts.ExportAllGob(accountsFile); err != nil {
			logger.Error("failed to export accounts:", err)
		}
	}()
	defer func() {
		if err = webhooks.ExportState(webhooksFile); err != nil {
			logger.Error("failed to export webhooks:", err)
		}
	}()
	defer func() {
		if err = harvesters.ExportState(harvestersFile); err != nil {
			logger.Error("failed to export harvesters:", err)
		}
	}()
	defer func() {
		if err = harvesters.ExportGmailState(gmailsFile); err != nil {
			logger.Error("failed to export harvester gmails:", err)
		}
	}()

	if err = rates.Initialize(); err != nil {
		logger.Error("failed to initialize rates:", err)
	}

	ints := integrations.Get()

	if ints.Aycd.Active {
		_ = autosolve.Connect(ints.Aycd.Credentials[0].Value, ints.Aycd.Credentials[1].Value)
	}

	var opts astilectron.Options
	if runtime.GOOS == "windows" {
		opts = astilectron.Options{
			AppName:            "Carina",
			SingleInstance:     true,
			AppIconDefaultPath: "resources/icons/icon.png",

			BaseDirectoryPath:  dataDir,
			DataDirectoryPath:  dataDir,
			VersionAstilectron: "0.50.8",
			VersionElectron:    "15.2.0",
		}
	} else {
		opts = astilectron.Options{
			AppName:            "Carina",
			SingleInstance:     true,
			BaseDirectoryPath:  dataDir,
			DataDirectoryPath:  dataDir,
			VersionAstilectron: "0.50.8",
			VersionElectron:    "15.2.0",
		}
	}

	a, err := astilectron.New(log.New(io.Discard, "", 0), opts)
	if err != nil {
		logger.Fatal("main: creating astilectron failed:", err)
	}
	defer a.Close()

	app.App = a
	windows.App = a
	app.App.HandleSignals()

	if err = app.App.Start(); err != nil {
		logger.Fatal("main: starting astilectron failed:", err)
	}

	// themida.Macro(themida.EAGLE_BLACK_END)

	go windowController()
	defer func() {
		// themida.Macro(themida.EAGLE_BLACK_START)

		s := clientutil.GetSession()
		if s == nil {
			return
		}
		_, err = rest.DeactivateSession(*s)
		if err != nil {
			logger.Error(err)
		}

		// themida.Macro(themida.EAGLE_BLACK_END)
	}()

	app.App.Wait()
}

func windowController() {
	// themida.Macro(themida.EAGLE_BLACK_START)

	sid := clientutil.GetSession()
	if sid == nil {
		app.AuthWindow()
		return
	}

	resp, err := rest.ActivateSession(*sid, hwid.Get())
	if err != nil || !resp.Success {
		app.AuthWindow()
		return
	}

	if err = ws.Connect(*sid); err != nil {
		app.AuthWindow()
		return
	}

	clientutil.SetUser(resp.User)

	// themida.Macro(themida.EAGLE_BLACK_END)

	app.MainWindow()
}

func setupRsa() {
	var pemBytes []byte
	err := json.Unmarshal([]byte(inject.RsaPem), &pemBytes)
	if err != nil {
		logger.Fatal(err)
	}

	keyBlock, _ := pem.Decode(pemBytes)
	if err != nil {
		logger.Fatal(err)
	}

	pubKey, err := x509.ParsePKIXPublicKey(keyBlock.Bytes)
	if err != nil {
		logger.Fatal(err)
	}

	if rsaPubKey, ok := pubKey.(*rsa.PublicKey); ok {
		client.InitRSA(rsaPubKey)
	} else {
		logger.Fatalf("expected pub key of type *crypto/rsa.PublicKey got type *%s.%s", reflect.TypeOf(pubKey).Elem().PkgPath(), reflect.TypeOf(pubKey).Elem().Name())
	}
}

func periodicExport(dataDir string, dur time.Duration) {
	for range time.Tick(dur) {
		_ = integrations.ExportIntegrations(path.Join(dataDir, IntegrationsFile))
		_ = userpref.ExportPreferences(path.Join(dataDir, UserPreferencesFile))
		_ = tasks.ExportAllGob(path.Join(dataDir, TasksFile))
		_ = profiles.ExportAllGob(path.Join(dataDir, ProfilesFile))
		_ = proxies.ExportAllGob(path.Join(dataDir, ProxiesFile))
		_ = accounts.ExportAllGob(path.Join(dataDir, AccountsFile))
		_ = webhooks.ExportState(path.Join(dataDir, WebhooksFile))
		_ = harvesters.ExportState(path.Join(dataDir, HarvestersFile))
		_ = harvesters.ExportGmailState(path.Join(dataDir, GmailsFile))
	}
}

func addDefaultGroups(harvestersPath string) error {
	if _, err := harvesters.GetHarvester("default"); err == harvesters.ErrHarvesterNotFound {
		harvesters.AddDefaultHarvester()
		if err = harvesters.ExportHarvesters(harvestersPath); err != nil {
			return err
		}
	}

	return nil
}
