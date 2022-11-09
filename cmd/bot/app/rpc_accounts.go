package app

import (
	"nebula/pkg/infra/accounts"
)

func AddAccounts(groupId string, data []accounts.AccountData) ([]string, error) {
	var accountIds []string
	for _, acct := range data {
		accountId, err := accounts.AddAccount(groupId, acct)
		if err != nil {
			return nil, err
		}
		accountIds = append(accountIds, accountId)
	}
	return accountIds, nil
}

func EditAccounts(data []accounts.AccountData) error {
	for _, acct := range data {
		if err := accounts.EditAccount(acct.Id, acct); err != nil {
			return err
		}
	}
	return nil
}

func MoveAccounts(acctIds []string, newGroupId string) error {
	for _, acctId := range acctIds {
		if err := accounts.MoveAccount(acctId, newGroupId); err != nil {
			return err
		}
	}
	return nil
}

func RemoveAccounts(acctIds []string) error {
	for _, acctId := range acctIds {
		if err := accounts.RemoveAccount(acctId); err != nil {
			return err
		}
	}
	return nil
}
