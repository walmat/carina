package proxies

import (
	"errors"
	"github.com/lithammer/shortuuid/v3"
	"nebula/pkg/infra/eventbus"
	"net/url"
	"sync"
)

var (
	proxyGroups  = make(map[string]ProxyGroup)
	proxyToGroup = make(map[string]string)
	proxiesMu    sync.RWMutex

	ProxyNotExistErr = errors.New("proxies: proxy does not exist")
	GroupNotExistErr = errors.New("proxies: group does not exist")
)

func init() {
	proxyGroups["default"] = ProxyGroup{Name: "Default", Proxies: newProxySet()}
}

func AddProxy(groupId string, rawUrl string) (string, error) {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	if !groupExists(groupId) {
		return "", GroupNotExistErr
	}

	pUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	proxyData := ProxyData{
		Id:  shortuuid.New(),
		Url: pUrl,
	}
	if err = proxyGroups[groupId].Proxies.Add(proxyData); err != nil {
		return "", err
	}
	proxyToGroup[proxyData.Id] = groupId
	eventbus.Publish(&ProxyAddedEvent{
		Id:      proxyData.Id,
		GroupId: groupId,
		Data:    proxyData,
	})
	return proxyData.Id, nil
}

func EditProxy(proxyId string, rawUrl string) error {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	if !proxyExists(proxyId) {
		return ProxyNotExistErr
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	groupId := proxyToGroup[proxyId]
	proxyData := ProxyData{
		Id:  proxyId,
		Url: u,
	}
	eventbus.Publish(&ProxyModifiedEvent{
		Id:   proxyId,
		Data: proxyData,
	})
	return proxyGroups[groupId].Proxies.Set(proxyData)
}

func RemoveProxy(proxyId string) error {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	if !proxyExists(proxyId) {
		return GroupNotExistErr
	}

	groupId := proxyToGroup[proxyId]
	delete(proxyToGroup, proxyId)
	eventbus.Publish(&ProxyRemovedEvent{
		Id:      proxyId,
		GroupId: groupId,
	})
	return proxyGroups[groupId].Proxies.Del(proxyId)
}

func GetProxy(proxyId string) (*ProxyData, error) {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()

	groupId := proxyToGroup[proxyId]
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	return proxyGroups[groupId].Proxies.Get(proxyId)
}

func AddGroup(name string) (string, error) {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	groupId := shortuuid.New()
	proxyGroups[groupId] = ProxyGroup{
		Name:    name,
		Proxies: newProxySet(),
	}
	eventbus.Publish(&GroupAddedEvent{
		Id: groupId,
	})
	return groupId, nil
}

func GetGroup(groupId string) (*ProxyGroup, error) {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	group := proxyGroups[groupId]
	return &group, nil
}

func GetGroups() []string {
	var groups []string
	for groupId := range proxyGroups {
		groups = append(groups, groupId)
	}
	return groups
}

func RemoveGroup(groupId string) error {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	if !groupExists(groupId) {
		return GroupNotExistErr
	}

	proxyGroup := proxyGroups[groupId]
	for _, proxy := range proxyGroup.Proxies.Values() {
		delete(proxyToGroup, proxy.Id)
	}

	delete(proxyGroups, groupId)
	eventbus.Publish(&GroupRemovedEvent{
		Id:   groupId,
		Data: proxyGroup,
	})
	return nil
}

func GetProxyGroup(proxyId string) (string, error) {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	if groupId, ok := proxyToGroup[proxyId]; ok {
		return groupId, nil
	} else {
		return "", ProxyNotExistErr
	}
}

func ProxyExists(proxyId string) bool {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	return proxyExists(proxyId)
}

func proxyExists(proxyId string) bool {
	_, ok := proxyToGroup[proxyId]
	return ok
}

func GroupExists(groupId string) bool {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	return groupExists(groupId)
}

func groupExists(groupId string) bool {
	_, ok := proxyGroups[groupId]
	return ok
}

func MoveProxy(proxyId, newGroupId string) error {
	proxiesMu.Lock()
	defer proxiesMu.Unlock()
	if !groupExists(newGroupId) {
		return GroupNotExistErr
	} else if !proxyExists(proxyId) {
		return ProxyNotExistErr
	}

	originalGroupId := proxyToGroup[proxyId]
	proxyData, err := proxyGroups[originalGroupId].Proxies.Get(proxyId)
	if err != nil {
		return err
	}
	if err = proxyGroups[originalGroupId].Proxies.Del(proxyId); err != nil {
		return err
	}
	proxyToGroup[proxyId] = newGroupId
	return proxyGroups[newGroupId].Proxies.Add(*proxyData)
}
