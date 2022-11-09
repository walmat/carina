package shopify

import (
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/webhooks"
)

func addWebhook(taskCtx *tasks.Context, runningCtx *runningCtx, staticCtx *staticCtx, status webhooks.OrderStatus, mode string) {
	tGroupId := taskCtx.GroupIdentifier()
	taskGroup, err := tasks.GetGroup(tGroupId)
	if err != nil {
		return
	}

	pGroupId, err := profiles.GetProfileGroup(runningCtx.ProfileData.Id)
	if err != nil {
		return
	}

	profileGroup, err := profiles.GetGroup(pGroupId)
	var proxyGroupName = "None"
	var proxyDataUrl = "None"

	if runningCtx.ProxyData != nil {
		proxyGroupId, err := proxies.GetProxyGroup(runningCtx.ProxyData.Id)
		if err != nil {
			return
		}

		proxyGroup, err := proxies.GetGroup(proxyGroupId)
		if err != nil {
			return
		}
		proxyDataUrl = runningCtx.ProxyData.Url.String()
		proxyGroupName = proxyGroup.Name
	}


	coData := webhooks.CheckoutData{
		Task: webhooks.Task{
			ID:        taskCtx.Identifier(),
			Mode:      mode,
			GroupName: taskGroup.Name,
		},
		Product: webhooks.Product{
			Name:     runningCtx.Product.Name,
			Price:    runningCtx.Product.Price,
			Image:    runningCtx.Product.Image,
			Size:    runningCtx.Product.Size,
			Quantity: "1",
			Url:      runningCtx.Product.Url,
		},
		Store: webhooks.Store{
			Name: staticCtx.Default.Store.Name,
			Url:  runningCtx.StoreUrl.String(),
		},
		Profile: webhooks.ProfileEmbed{
			ID:        runningCtx.ProfileData.Id,
			GroupName: profileGroup.Name,
			Email:     runningCtx.ProfileData.Payment.Email,
			Name:      runningCtx.ProfileData.Name,
		},
		Proxy: webhooks.GroupedObject{
			GroupName: proxyGroupName,
			Name:      proxyDataUrl,
		},
		Order: webhooks.Order{
			ID:  "",
			Url: "",
		},
		Status: status,
	}

	webhooks.Queue(coData)
}
