export const ACCOUNT_COLUMNS = [
	{
		index: 0,
		accessor: "store.name",
		Header: "Store",
		sortKeys: ["store"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 1,
		accessor: "username",
		Header: "Username / Email Address",
		sortKeys: ["username"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 2,
		type: "password",
		accessor: "password",
		Header: "Password",
		sortKeys: ["password"],
		isHidden: false,
		isSortable: true,
	},
];
