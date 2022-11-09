import { Defaults } from "../../../stores/Main/reducers/defaults";

export const DEFAULTS_COLUMNS = [
	{
		index: 0,
		maxWidth: 112,
		minWidth: 80,
		width: 80,
		accessor: "name",
		Header: "Name",
		sortKeys: ["name"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 1,
		accessor: "stores.name",
		Header: "Store",
		sortKeys: ["store"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 2,
		accessor: "sizes",
		Header: "Sizes",
		sortKeys: ["sizes"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 3,
		accessor: "profileGroup",
		Header: "Profiles",
		sortKeys: ["profiles"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 4,
		accessor: "proxies",
		Header: "Proxies",
		sortKeys: ["proxies"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 5,
		accessor: "accounts",
		Header: "Accounts",
		sortKeys: ["accounts"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 6,
		accessor: "rate",
		Header: "Rates",
		sortKeys: ["rate"],
		isHidden: false,
		isSortable: true,
	},
];

export const extractAllStores = (groups: Defaults) => {
	let byId: any = {};

	for (const group of Object.values(groups)) {
		for (const d of Object.values(group.byId)) {
			byId[d.id] = d;
		}
	}

	return { byId };
};