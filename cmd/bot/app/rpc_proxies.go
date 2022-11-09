package app

import (
	"nebula/pkg/infra/proxies"
)

func AddProxies(groupId string, urls []string) ([]string, error) {
	var proxyIds []string
	for _, rawUrl := range urls {
		proxyId, err := proxies.AddProxy(groupId, rawUrl)
		if err != nil {
			return nil, err
		}
		proxyIds = append(proxyIds, proxyId)
	}
	return proxyIds, nil
}

func TestProxies(proxyIds []string, groupId string, urls []string) {
	// TODO: Send batched events back to main window to update speed
}

type proxyData struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func EditProxies(proxyData []proxyData) error {
	for _, proxy := range proxyData {
		if err := proxies.EditProxy(proxy.Id, proxy.Url); err != nil {
			return err
		}
	}
	return nil
}

func MoveProxies(proxyIds []string, newGroupId string) error {
	for _, proxyId := range proxyIds {
		if err := proxies.MoveProxy(proxyId, newGroupId); err != nil {
			return err
		}
	}
	return nil
}

func RemoveProxies(proxyIds []string) error {
	for _, proxyId := range proxyIds {
		if err := proxies.RemoveProxy(proxyId); err != nil {
			return err
		}
	}
	return nil
}
