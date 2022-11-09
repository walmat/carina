package accounts

type AccountGroup struct {
	Name     string                 `json:"name"`
	Accounts map[string]AccountData `json:"accounts"`
}

type Store struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type AccountData struct {
	Id       string `json:"id"`
	Store    Store 	`json:"store"`
	Username string `json:"username"`
	Password string `json:"password"`
}
