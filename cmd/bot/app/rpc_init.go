package app

import (
	"nebula/pkg/infra/accounts"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/rates"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/update"
	"nebula/pkg/infra/webhooks"
	"nebula/pkg/integrations"
	"nebula/pkg/integrations/autosolve"
	"nebula/pkg/userpref"
)

type LoginStatus struct {
	Success bool
	Message string
}

func expectRPC(name string, fn interface{}) {
	if err := RegisterRPC(name, fn); err != nil {
		panic(err)
	}
}

func init() {
	// region AUTH RPC
	expectRPC("auth:login", authLogin)
	expectRPC("auth:forgot", forgotPassword)
	expectRPC("auth:change", changePassword)
	expectRPC("auth:2fa", submitOtp)
	expectRPC("auth:register", register)
	// endregion

	// region updater RPC
	expectRPC("updater:update", update.Update)
	expectRPC("updater:check", update.CheckForUpdates)
	// endregion

	// region STATE RPC
	expectRPC("state:get", getState)
	// endregion

	// region TASKS RPC
	expectRPC("tasks:add", AddTasks)
	expectRPC("tasks:start", StartTasks)
	expectRPC("tasks:stop", StopTasks)
	expectRPC("tasks:delete", RemoveTasks)
	expectRPC("tasks:move", MoveTasks)
	expectRPC("tasks:edit", tasks.EditTask)
	expectRPC("tasks:addGroup", tasks.AddGroup)
	expectRPC("tasks:deleteGroup", tasks.RemoveGroup)
	expectRPC("tasks:importJson", tasks.ImportJson)
	expectRPC("tasks:exportAllJson", tasks.ExportAllJson)
	expectRPC("tasks:importGob", tasks.ImportGob)
	expectRPC("tasks:exportAllGob", tasks.ExportAllGob)
	// endregion

	// region PROFILES RPC
	expectRPC("profiles:add", AddProfiles)
	expectRPC("profiles:edit", EditProfiles)
	expectRPC("profiles:move", MoveProfiles)
	expectRPC("profiles:delete", RemoveProfiles)
	expectRPC("profiles:assign", profiles.AssignProfileToTask)
	expectRPC("profiles:addGroup", profiles.AddGroup)
	expectRPC("profiles:deleteGroup", profiles.RemoveGroup)
	expectRPC("profiles:exportAllGob", profiles.ExportAllGob)
	expectRPC("profiles:importGob", profiles.ImportGob)
	expectRPC("profiles:exportAllJson", profiles.ExportAllJson)
	expectRPC("profiles:importJson", profiles.ImportJson)
	// endregion

	// region PROXIES RPC
	expectRPC("proxies:add", AddProxies)
	expectRPC("proxies:edit", EditProxies)
	expectRPC("proxies:move", MoveProxies)
	expectRPC("proxies:delete", RemoveProxies)
	expectRPC("proxies:assign", proxies.AssignProxyGroupToTask)
	expectRPC("proxies:addGroup", proxies.AddGroup)
	expectRPC("proxies:deleteGroup", proxies.RemoveGroup)
	expectRPC("proxies:test", TestProxies)
	expectRPC("proxies:exportAllGob", proxies.ExportAllGob)
	expectRPC("proxies:importGob", proxies.ImportGob)
	expectRPC("proxies:exportAllJson", proxies.ExportAllJson)
	expectRPC("proxies:importJson", proxies.ImportJson)
	// endregion

	// region ACCOUNTS RPC
	expectRPC("accounts:add", AddAccounts)
	expectRPC("proxies:edit", EditAccounts)
	expectRPC("proxies:move", MoveAccounts)
	expectRPC("accounts:delete", RemoveAccounts)
	expectRPC("accounts:assign", accounts.AssignAccountToTask)
	expectRPC("accounts:addGroup", accounts.AddGroup)
	expectRPC("accounts:deleteGroup", accounts.RemoveGroup)
	expectRPC("proxies:exportAllGob", accounts.ExportAllGob)
	expectRPC("proxies:importGob", accounts.ImportGob)
	expectRPC("proxies:exportAllJson", accounts.ExportAllJson)
	expectRPC("proxies:importJson", accounts.ImportJson)
	// endregion

	// region WEBHOOK RPC
	expectRPC("webhooks:add", webhooks.AddWebhook)
	expectRPC("webhooks:edit", webhooks.EditWebhook)
	expectRPC("webhooks:delete", webhooks.DeleteWebhook)
	expectRPC("webhooks:sendToWebhook", webhooks.SendToWebhook)
	expectRPC("webhooks:sendWebhook", webhooks.Queue)
	// endregion

	// region INTEGRATIONS RPC
	expectRPC("integrations:get", integrations.Get)
	expectRPC("integrations:set", integrations.Set)
	expectRPC("autosolve:connect", autosolve.Connect)
	expectRPC("autosolve:close", autosolve.Close)
	// endregion

	// region THEME RPC
	expectRPC("theme:set", userpref.SetThemeState)
	expectRPC("theme:get", userpref.GetThemeState)
	// endregion

	// region USERPREF RPC
	expectRPC("preferences:get", userpref.Get)
	expectRPC("preferences:set", userpref.Set)
	expectRPC("preferences:setCheckoutSound", userpref.SetCheckoutSound)
	expectRPC("preferences:setCheckoutVolume", userpref.SetCheckoutVolume)
	expectRPC("preferences:playCheckout", userpref.PlayCheckout)
	expectRPC("preferences:setHarvesterSound", userpref.SetHarvesterSound)
	expectRPC("preferences:setHarvesterVolume", userpref.SetHarvesterVolume)
	expectRPC("preferences:playHarvester", userpref.PlayHarvester)
	expectRPC("preferences:stopAllSounds", userpref.StopAllSounds)
	// endregion

	// region FS RPC
	expectRPC("fs:selectFile", SelectFile)
	expectRPC("fs:selectFileFiltered", SelectFileFiltered)
	// endregion

	// region RATES RPC
	expectRPC("rates:fetch", rates.FetchRates)
	// endregion

	// region COLLECTIVE RPC
	expectRPC("collective:add", harvesters.AddHarvester)
	expectRPC("collective:edit", harvesters.EditHarvester)
	expectRPC("collective:focus", harvesters.FocusHarvester)
	expectRPC("collective:show", harvesters.ShowHarvester)
	expectRPC("collective:hide", harvesters.HideHarvester)
	expectRPC("collective:delete", harvesters.DeleteHarvester)
	expectRPC("collective:import", harvesters.ImportHarvesters)
	expectRPC("collective:export", harvesters.ExportHarvesters)
	// endregion

	// region GMAIL RPC
	expectRPC("gmail:add", harvesters.AddGmail)
	expectRPC("gmail:select", harvesters.SelectGmail)
	expectRPC("gmail:delete", harvesters.DeleteGmail)
	// endregion

	// region WEBHOOKS RPC
	expectRPC("webhooks:add", webhooks.AddWebhook)
	// expectRPC("webhook:test", webhooks.TestWebhook)
	expectRPC("webhooks:edit", webhooks.EditWebhook)
	expectRPC("webhooks:delete", webhooks.DeleteWebhook)
	expectRPC("webhooks:export", webhooks.ExportWebhooks)
	expectRPC("webhooks:import", webhooks.ImportWebhooks)
	// endregion
}
