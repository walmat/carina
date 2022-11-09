package accounts

import (
	"errors"
	"github.com/lithammer/shortuuid/v3"
	"nebula/pkg/infra/eventbus"
	"sync"
)

var (
	accountGroups  = make(map[string]AccountGroup)
	accountToGroup = make(map[string]string)
	accountsMu     sync.RWMutex

	AccountNotExistErr = errors.New("accounts: account does not exist")
	GroupNotExistErr   = errors.New("accounts: group does not exist")
)

func init() {
	accountGroups["default"] = AccountGroup{Name: "Default", Accounts: make(map[string]AccountData)}
}

func AddAccount(groupId string, data AccountData) (string, error) {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	if !groupExists(groupId) {
		return "", GroupNotExistErr
	}

	accountId := shortuuid.New()
	data.Id = accountId

	accountGroups[groupId].Accounts[accountId] = data
	accountToGroup[accountId] = groupId
	eventbus.Publish(&AccountAddedEvent{
		Id:      accountId,
		GroupId: groupId,
		Data:    data,
	})
	return accountId, nil
}

func EditAccount(accountId string, data AccountData) error {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	if !accountExists(accountId) {
		return AccountNotExistErr
	}

	groupId := accountToGroup[accountId]
	accountGroups[groupId].Accounts[groupId] = data
	eventbus.Publish(&AccountModifiedEvent{
		Id:   accountId,
		Data: data,
	})
	return nil
}

func RemoveAccount(accountId string) error {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	if !accountExists(accountId) {
		return GroupNotExistErr
	}

	groupId := accountToGroup[accountId]
	delete(accountGroups[groupId].Accounts, accountId)
	delete(accountToGroup, accountId)
	eventbus.Publish(&AccountRemovedEvent{
		Id:      accountId,
		GroupId: groupId,
	})
	return nil
}

func GetAccount(accountId string) (*AccountData, error) {
	accountsMu.RLock()
	defer accountsMu.RUnlock()

	groupId := accountToGroup[accountId]
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	if acct, ok := accountGroups[groupId].Accounts[accountId]; ok {
		return &acct, nil
	} else {
		return nil, AccountNotExistErr
	}
}

func AddGroup(name string) (string, error) {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	groupId := shortuuid.New()
	accountGroups[groupId] = AccountGroup{
		Name:     name,
		Accounts: make(map[string]AccountData),
	}
	eventbus.Publish(&GroupAddedEvent{
		Id: groupId,
	})
	return groupId, nil
}

func GetGroup(groupId string) (*AccountGroup, error) {
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	group := accountGroups[groupId]
	return &group, nil
}

func GetGroups() []string {
	var groups []string
	for groupId := range accountGroups {
		groups = append(groups, groupId)
	}
	return groups
}

func RemoveGroup(groupId string) error {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	if !groupExists(groupId) {
		return GroupNotExistErr
	}

	for accountId := range accountGroups[groupId].Accounts {
		delete(accountToGroup, accountId)
		eventbus.Publish(&AccountRemovedEvent{
			Id:      accountId,
			GroupId: groupId,
		})
	}

	delete(accountGroups, groupId)
	eventbus.Publish(&GroupRemovedEvent{
		Id: groupId,
	})
	return nil
}

func GetAccountGroup(accountId string) (string, error) {
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	if groupId, ok := accountToGroup[accountId]; ok {
		return groupId, nil
	} else {
		return "", AccountNotExistErr
	}
}

func AccountExists(accountId string) bool {
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	return accountExists(accountId)
}

func accountExists(accountId string) bool {
	_, ok := accountToGroup[accountId]
	return ok
}

func GroupExists(groupId string) bool {
	accountsMu.RLock()
	defer accountsMu.RUnlock()
	return groupExists(groupId)
}

func groupExists(groupId string) bool {
	_, ok := accountGroups[groupId]
	return ok
}

func MoveAccount(accountId, newGroupId string) error {
	accountsMu.Lock()
	defer accountsMu.Unlock()
	if !groupExists(newGroupId) {
		return GroupNotExistErr
	} else if !accountExists(accountId) {
		return AccountNotExistErr
	}

	originalGroupId := accountToGroup[accountId]
	accountData := accountGroups[originalGroupId].Accounts[accountId]
	delete(accountGroups[originalGroupId].Accounts, accountId)
	accountGroups[newGroupId].Accounts[accountId] = accountData
	accountToGroup[accountId] = newGroupId
	return nil
}
