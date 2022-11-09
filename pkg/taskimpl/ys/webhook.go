package ys

import (
	"fmt"
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
	"nebula/pkg/infra/webhooks"
	"strconv"
)

func addWebhook(taskCtx *tasks.Context, runningCtx *runningCtx, status webhooks.OrderStatus) {
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

	proxyGroupId, err := proxies.GetProxyGroup(runningCtx.ProxyData.Id)
	if err != nil {
		return
	}

	proxyGroup, err := proxies.GetGroup(proxyGroupId)
	if err != nil {
		return
	}

	coData := webhooks.CheckoutData{
		Task: webhooks.Task{
			ID:        taskCtx.Identifier(),
			Mode:      "Normal", // TODO: Get mode from task ctx or something
			GroupName: taskGroup.Name,
		},
		Product: webhooks.Product{
			Name:     runningCtx.ProductName,
			Price:    strconv.FormatInt(runningCtx.ProductPrice, 10),
			Image:    runningCtx.ProductImage,
			Size:     runningCtx.Size,
			Quantity: "1",
			Url:      fmt.Sprintf("https://www.yeezysupply.com/product/%s", runningCtx.ProductId),
		},
		Store: webhooks.Store{
			Name: "Yeezy Supply",
			Url:  runningCtx.StoreUrl.String(),
		},
		Profile: webhooks.ProfileEmbed{
			ID:        runningCtx.ProfileData.Id,
			GroupName: profileGroup.Name,
			Email:     runningCtx.ProfileData.Payment.Email,
			Name:      runningCtx.ProfileData.Name,
		},
		Proxy: webhooks.GroupedObject{
			GroupName: proxyGroup.Name,
			Name:      runningCtx.ProxyData.Url.String(),
		},
		Order: webhooks.Order{
			ID:  "",
			Url: "",
		},
		Status: status,
	}

	webhooks.Queue(coData)
}
