package app

import (
	"nebula/cmd/bot/inject"
	"nebula/pkg/api/model"
	"nebula/pkg/clientutil"
	"nebula/pkg/infra/accounts"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/rates"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/update"
	"nebula/pkg/infra/webhooks"
	"nebula/pkg/integrations"
	"nebula/pkg/userpref"

	"github.com/Nebulabots/go-astilectron"
)

type State struct {
	User         model.AuthUser
	Tasks        map[string]StateTaskGroup
	Profiles     map[string]StateProfileGroup
	Proxies      map[string]StateProxyGroup
	Accounts     map[string]StateAccountGroup
	Webhooks     map[string]StateWebhook
	Harvesters   StateHarvesterGroup
	Gmails       StateGmailsGroup
	Stores       []StateStoreObject
	Integrations integrations.Integrations
	Settings     userpref.UserPreferences
	Theme        uint8
	Version      string
	Update       bool
}

type StateTask struct {
	Id      string       `json:"id"`
	Type    string       `json:"type"`
	Store   TableStore   `json:"store"`
	Profile TableProfile `json:"profile"`
	Proxies TableProxies `json:"proxies"`
	Product TableProduct `json:"product"`
	Status  string       `json:"status"`
}

type StateHarvester struct {
	ID      string                `json:"id"`
	Name    string                `json:"name"`
	Index   int                   `json:"index"`
	Focused bool                  `json:"focused"`
	Store   *harvesters.Store     `json:"store"`
	Account *harvesters.GmailData `json:"account"`
	Proxy   string                `json:"proxy"`
	Solver  *harvesters.Solver    `json:"solver"`
}

type StateGmailData struct {
	ID      string                       `json:"id"`
	Email   string                       `json:"email"`
	Cookies *[]astilectron.SessionCookie `json:"cookies"`
}

type TableStore struct {
	Name    string `json:"name"`
	Meta    string `json:"meta"`
	Details string `json:"details"`
}

type TableProduct struct {
	Name    string `json:"name"`
	Meta    string `json:"meta"`
	Details string `json:"details"`
}

type TableProxies struct {
	Name    string `json:"name"`
	Meta    string `json:"meta"`
	Details string `json:"details"`
}

type StateStoreObject struct {
	Options      []tasks.Store `json:"options"`
	Type         string        `json:"type"`
	Index        int           `json:"index"`
	Label        string        `json:"label"`
	UsesRates    bool          `json:"usesRates"`
	UsesAccounts bool          `json:"usesAccounts"`
}

type StateTaskGroup struct {
	Id   string               `json:"id"`
	Name string               `json:"name"`
	ById map[string]StateTask `json:"byId"`
}

type StateHarvesterGroup struct {
	ById map[string]StateHarvester `json:"byId"`
}

type StateGmailsGroup struct {
	ById map[string]StateGmailData `json:"byId"`
}

type StateProfile struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Billing      profiles.Address `json:"billing,omitempty"`
	Shipping     profiles.Address `json:"shipping,omitempty"`
	Payment      profiles.Payment `json:"payment,omitempty"`
	MaxCheckouts int              `json:"maxCheckouts,omitempty"`
}

type TableProfile struct {
	Name    string `json:"name"`
	Meta    string `json:"meta"`
	Details string `json:"details"`
}

type StateProfileGroup struct {
	Id   string                  `json:"id"`
	Name string                  `json:"name"`
	ById map[string]StateProfile `json:"byId"`
}

type StateProxy struct {
	Id     string  `json:"id"`
	Active bool    `json:"active"`
	Speed  *int    `json:"speed"`
	Host   string  `json:"host"`
	Port   string  `json:"port"`
	User   *string `json:"username"`
	Pass   *string `json:"password"`
}

type StateProxyGroup struct {
	Id   string                `json:"id"`
	Name string                `json:"name"`
	ById map[string]StateProxy `json:"byId"`
}

type StateAccount struct {
	Id       string         `json:"id"`
	Store    accounts.Store `json:"store"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}

type StateAccountGroup struct {
	Id   string                  `json:"id"`
	Name string                  `json:"name"`
	ById map[string]StateAccount `json:"byId"`
}

type StateWebhook struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Active      bool                    `json:"active"`
	Type        string                  `json:"type"`
	Url         string                  `json:"url"`
	Profiles    []webhooks.Profile      `json:"profiles"`
	Fields      []webhooks.WebhookField `json:"fields"`
	Declines    bool                    `json:"declines"`
	Sensitivity bool                    `json:"sensitivity"`
}

func getState(stateType string) State {
	if stateType == "auth" {
		return State{
			Theme: userpref.GetThemeState(),
		}
	}

	if stateType == "collective" {
		var harvestersData StateHarvesterGroup
		var gmailsData StateGmailsGroup
		harvestersData.ById = make(map[string]StateHarvester)
		gmailsData.ById = make(map[string]StateGmailData)

		harvesterIds := harvesters.GetHarvesterIds()
		for _, harvesterId := range harvesterIds {
			harvester, err := harvesters.GetHarvester(harvesterId)
			if err != nil {
				continue
			}

			focus := false
			if harvesterId == "default" {
				focus = true
			}

			harvestersData.ById[harvesterId] = StateHarvester{
				ID:      harvesterId,
				Index:   harvester.Index,
				Focused: focus,
				Name:    harvester.Name,
				Store:   harvester.Store,
				Account: harvester.Account,
				Proxy:   harvester.Proxy,
			}
		}

		if harvestersData.ById == nil {
			harvestersData.ById["default"] = StateHarvester{
				ID:      "default",
				Index:   0,
				Focused: true,
				Name:    "Default",
			}
		}

		gmailIds := harvesters.GetGmailIds()

		for _, gmailId := range gmailIds {
			gmail, err := harvesters.GetGmail(gmailId)
			if err != nil {
				continue
			}

			gmailsData.ById[gmailId] = StateGmailData{
				ID:      gmailId,
				Email:   gmail.Email,
				Cookies: gmail.Cookies,
			}
		}

		return State{
			Theme:      userpref.GetThemeState(),
			Harvesters: harvestersData,
			Gmails:     gmailsData,
		}
	}

	taskGroups := make(map[string]StateTaskGroup)
	taskGroupIds := tasks.GetGroups()
	for _, groupId := range taskGroupIds {
		group, err := tasks.GetGroup(groupId)
		if err != nil {
			continue
		}

		byId := make(map[string]StateTask)
		for taskId := range group.Tasks {
			typeId, err := tasks.GetTaskType(taskId)
			if err != nil {
				continue
			}

			modeId, err := tasks.GetTaskMode(taskId)
			if err != nil {
				continue
			}

			profile, err := profiles.GetProfileForTask(taskId)
			if err != nil {
				continue
			}
			profileGroupId, err := profiles.GetProfileGroup(profile.Id)
			if err != nil {
				continue
			}
			profileGroup, err := profiles.GetGroup(profileGroupId)
			if err != nil {
				continue
			}

			proxy, _ := proxies.GetProxyForTask(taskId)
			var proxyGroup *proxies.ProxyGroup
			var proxyGroupId string
			if proxy != nil {
				proxyGroupId, _ = proxies.GetProxyGroup(proxy.Id)
				proxyGroup, _ = proxies.GetGroup(proxyGroupId)
			}

			meta, err := tasks.GetMetadataForTaskType(typeId, modeId)
			if err != nil {
				continue
			}

			taskCtx, err := tasks.GetStaticContextForTaskId(taskId)
			if err != nil {
				continue
			}

			productName := ""
			if product := taskCtx.GetProductInfo(); product != nil {
				if product.Name != nil {
					productName = *product.Name
				}
			}

			proxyGroupName := "No Group"
			if proxyGroup != nil {
				proxyGroupName = proxyGroup.Name
			}

			m, _ := tasks.GetTaskMode(taskId)

			byId[taskId] = StateTask{
				Id:   taskId,
				Type: typeId,
				Store: TableStore{
					Name:    meta.DisplayName,
					Meta:    meta.StoreUrl,
					Details: m,
				},
				Profile: TableProfile{
					Name:    profile.Name,
					Meta:    profile.Id,
					Details: profileGroup.Name,
				},
				Proxies: TableProxies{
					Name:    proxyGroupName,
					Meta:    proxyGroupId,
					Details: "None",
				},
				Product: TableProduct{
					Name:    productName,
					Details: "Random", // TODO: Get sizes based on task id
				},
				Status: "Idle",
			}
		}

		taskGroups[groupId] = StateTaskGroup{
			Id:   groupId,
			Name: group.Name,
			ById: byId,
		}
	}

	webhooksData := make(map[string]StateWebhook)
	webhookIds := webhooks.GetWebhookIds()
	for _, webhookId := range webhookIds {
		webhook, err := webhooks.GetWebhook(webhookId)
		if err != nil {
			continue
		}

		webhooksData[webhookId] = StateWebhook{
			Id:          webhookId,
			Name:        webhook.Name,
			Active:      webhook.Active,
			Type:        webhook.Type,
			Url:         webhook.Url,
			Profiles:    webhook.Profiles,
			Fields:      webhook.Fields,
			Declines:    webhook.Declines,
			Sensitivity: webhook.Sensitivity,
		}
	}
	profileGroups := make(map[string]StateProfileGroup)
	profileGroupIds := profiles.GetGroups()
	for _, groupId := range profileGroupIds {
		group, err := profiles.GetGroup(groupId)
		if err != nil {
			continue
		}

		byId := make(map[string]StateProfile)
		for profileId := range group.Profiles {
			profile, err := profiles.GetProfile(profileId)
			if err != nil {
				continue
			}

			byId[profileId] = StateProfile{
				Id:           profileId,
				Name:         profile.Name,
				Billing:      profile.Billing,
				Shipping:     profile.Shipping,
				Payment:      profile.Payment,
				MaxCheckouts: profile.MaxCheckouts,
			}
		}

		profileGroups[groupId] = StateProfileGroup{
			Id:   groupId,
			Name: group.Name,
			ById: byId,
		}
	}

	proxyGroups := make(map[string]StateProxyGroup)
	proxyGroupIds := proxies.GetGroups()
	for _, groupId := range proxyGroupIds {
		group, err := proxies.GetGroup(groupId)
		if err != nil {
			continue
		}

		byId := make(map[string]StateProxy)
		for _, proxy := range group.Proxies.Values() {
			if err != nil {
				continue
			}

			var user, pass *string
			if proxy.Url.User != nil {
				username := proxy.Url.User.Username()
				user = &username
				password, set := proxy.Url.User.Password()
				if set {
					pass = &password
				}
			}

			byId[proxy.Id] = StateProxy{
				Id:   proxy.Id,
				Host: proxy.Url.Hostname(),
				Port: proxy.Url.Port(),
				User: user,
				Pass: pass,
			}
		}

		proxyGroups[groupId] = StateProxyGroup{
			Id:   groupId,
			Name: group.Name,
			ById: byId,
		}
	}

	accountGroups := make(map[string]StateAccountGroup)
	accountGroupIds := accounts.GetGroups()
	for _, groupId := range accountGroupIds {
		group, err := accounts.GetGroup(groupId)
		if err != nil {
			continue
		}

		byId := make(map[string]StateAccount)
		for acctId, acct := range group.Accounts {
			accountId := acctId

			byId[accountId] = StateAccount{
				Id:       accountId,
				Store:    acct.Store,
				Username: acct.Username,
				Password: acct.Password,
			}
		}

		accountGroups[groupId] = StateAccountGroup{
			Id:   groupId,
			Name: group.Name,
			ById: byId,
		}
	}

	var stateStoreObjects []StateStoreObject
	for _, typ := range tasks.GetTypes() {
		m := tasks.GetModes(typ)
		staticCtx, err := tasks.GetStaticContextForTaskType(typ, m[0])
		if err != nil {
			continue
		}
		metadata := staticCtx.Metadata()

		var stores []tasks.Store
		if len(metadata.Stores) == 0 {
			stores = []tasks.Store{
				{
					Name: metadata.DisplayName,
					Url:  metadata.StoreUrl,
				},
			}
		} else {
			stores = metadata.Stores
		}

		_, rateFetcherSupport := staticCtx.(rates.RateFetcher)
		stateStoreObjects = append(stateStoreObjects, StateStoreObject{
			Options:      stores,
			Type:         typ,
			Index:        len(stateStoreObjects),
			Label:        metadata.DisplayName,
			UsesRates:    rateFetcherSupport,
			UsesAccounts: metadata.UsesAccounts,
		})
	}

	needsUpdate, err := update.CheckForUpdates()

	if err != nil {
		panic(err)
	}

	return State{
		User:         *clientutil.GetUser(),
		Tasks:        taskGroups,
		Profiles:     profileGroups,
		Proxies:      proxyGroups,
		Accounts:     accountGroups,
		Webhooks:     webhooksData,
		Stores:       stateStoreObjects,
		Integrations: integrations.Get(),
		Settings:     userpref.Get(),
		Theme:        userpref.GetThemeState(),
		Version:      inject.CommitHash,
		Update:       needsUpdate,
	}
}
