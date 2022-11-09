import { Rate } from '../../../stores/Main/reducers/rates';

export const filterRates = (rate: Rate, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(rate.location)) {
		return true;
	}

	if (matches(rate.rate)) {
		return true;
	}

	if (matches(rate.store.name) || matches(rate.store.url)) {
		return true;
	}

	// search for 'up-to' rate price
	if (!Number.isNaN(safeSearch) && Number(rate.price) >= Number(safeSearch)) {
		return true;
	}

	return matches(rate.carrier);
}
