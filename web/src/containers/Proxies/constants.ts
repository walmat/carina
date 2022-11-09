export const PROXY_COLUMNS = [
	{
		index: 0,
		maxWidth: 256,
		minWidth: 112,
		width: 192,
		accessor: "host",
		Header: "Hostname",
		sortKeys: ["host"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 1,
		maxWidth: 112,
		minWidth: 64,
		width: 92,
		accessor: "port",
		Header: "Port",
		sortKeys: ["port"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 2,
		maxWidth: 112,
		minWidth: 96,
		width: 112,
		accessor: "username",
		Header: "Username",
		sortKeys: ["username"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 3,
		maxWidth: 112,
		minWidth: 96,
		width: 112,
		type: "password",
		accessor: "password",
		Header: "Password",
		sortKeys: ["password"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 4,
		maxWidth: 112,
		minWidth: 64,
		width: 92,
		accessor: "speed",
		Header: "Speed",
		sortKeys: ["speed"],
		isHidden: false,
		isSortable: true,
	},
	{
		index: 5,
		maxWidth: 112,
		minWidth: 64,
		width: 92,
		accessor: "active",
		Header: "Active",
		sortKeys: ["active"],
		isHidden: false,
		isSortable: true,
	},
];