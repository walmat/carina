package proxies

import (
	"errors"
	"nebula/pkg/infra/eventbus"
	"nebula/pkg/infra/tasks"
	"sync"
)

var (
	links        = make(map[string]string)
	linksMu      sync.RWMutex
	roundRobin   = make(map[string]int)
	roundRobinMu sync.RWMutex

	NotLinkedErr = errors.New("proxies: proxy not linked")
)

func init() {
	// remove links related to deleted tasks
	_ = eventbus.Subscribe(func(e *tasks.TaskRemovedEvent) {
		linksMu.Lock()
		defer linksMu.Unlock()
		for k := range links {
			if k == e.Id {
				delete(links, k)
			}
		}
	})

	// remove tasks that are linked to this profile
	_ = eventbus.Subscribe(func(e *GroupRemovedEvent) {
		linksMu.Lock()
		defer linksMu.Unlock()
		for k, v := range links {
			if v == e.Id {
				_ = tasks.RemoveTask(k)
				delete(links, k)
			}
		}
	})
}

func AssignProxyGroupToTask(taskId, groupId string) error {
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	if !groupExists(groupId) {
		return GroupNotExistErr
	}

	linksMu.Lock()
	defer linksMu.Unlock()
	links[taskId] = groupId
	return nil
}

func GetProxyForTask(taskId string) (*ProxyData, error) {
	linksMu.RLock()
	defer linksMu.RUnlock()
	proxiesMu.RLock()
	defer proxiesMu.RUnlock()
	groupId, ok := links[taskId]
	if !ok {
		return nil, NotLinkedErr
	} else if !groupExists(groupId) {
		return nil, ProxyNotExistErr
	}

	proxies := proxyGroups[groupId].Proxies.Values()
	if len(proxies) == 0 {
		return nil, NotLinkedErr
	}

	roundRobinMu.Lock()
	defer roundRobinMu.Unlock()
	idx, ok := roundRobin[groupId]
	if !ok {
		roundRobin[groupId] = 0
	} else {
		if len(proxies) == idx+1 {
			roundRobin[groupId] = 0
		} else {
			roundRobin[groupId] = idx + 1
		}
	}
	return &proxies[idx], nil
}
