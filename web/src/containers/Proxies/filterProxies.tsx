import { Proxy } from '../../stores/Main/reducers/proxies';

export const filterProxies = (proxy: Proxy, search: string) => {
	const safeSearch = search.trim().toLowerCase();

	const matches = (value: string) => value.toLowerCase().includes(safeSearch);

	if (matches(proxy.host) || matches(proxy.port)) {
		return true;
	}

	if (proxy.username && matches(proxy.username)) {
		return true;
	}

	if (proxy.password && matches(proxy.password)) {
		return true;
	}

	if (proxy.active && matches(`${proxy.active}`)) {
		return true;
	}

	return !!(proxy.speed && matches(`${proxy.active}`));
}
