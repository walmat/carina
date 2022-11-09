import { Account } from '../../../stores/Main/reducers/accounts';

export const filterAccounts = (account: Account, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(account.store.name) || matches(account.store.url)) {
		return true;
	}

	return matches(account.username) || matches(account.password);
}
