package accounts

import (
	"errors"
	"nebula/pkg/infra/eventbus"
	"nebula/pkg/infra/tasks"
	"sync"
)

var (
	links   = make(map[string]string)
	linksMu sync.RWMutex

	NotLinkedErr = errors.New("accounts: account not linked")
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

	// remove tasks that are linked to this account
	_ = eventbus.Subscribe(func(e *AccountRemovedEvent) {
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

func AssignAccountToTask(taskId, accountId string) error {
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	if !accountExists(accountId) {
		return AccountNotExistErr
	}

	linksMu.Lock()
	defer linksMu.Unlock()
	links[taskId] = accountId
	eventbus.Publish(&AccountAssignedToTaskEvent{
		Id:     accountId,
		TaskId: taskId,
	})
	return nil
}

func GetAccountForTask(taskId string) (*AccountData, error) {
	linksMu.RLock()
	defer linksMu.RUnlock()
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	accountId, ok := links[taskId]
	if !ok {
		return nil, NotLinkedErr
	} else if !accountExists(accountId) {
		return nil, AccountNotExistErr
	}

	groupId := accountToGroup[accountId]
	accountData := accountGroups[groupId].Accounts[accountId]
	return &accountData, nil
}
