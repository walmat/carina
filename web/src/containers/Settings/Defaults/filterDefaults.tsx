import { DefaultEntry } from '../../../stores/Main/reducers/defaults';

export const filterDefaults = (d: DefaultEntry, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(d.name)) {
		return true;
	}

	if (matches(d.store.name) || matches(d.store.url)) {
		return true;
	}

	// TODO: Add missing defaults fields
	return matches(d.profiles.name);
}
