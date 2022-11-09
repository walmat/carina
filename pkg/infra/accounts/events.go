package accounts

type (
	AccountAddedEvent struct {
		Id      string      `json:"id"`
		GroupId string      `json:"groupId"`
		Data    AccountData `json:"data"`
	}

	AccountRemovedEvent struct {
		Id      string `json:"id"`
		GroupId string `json:"groupId"`
	}

	AccountModifiedEvent struct {
		Id   string      `json:"id"`
		Data AccountData `json:"data"`
	}

	AccountAssignedToTaskEvent struct {
		Id     string `json:"id"`
		TaskId string `json:"taskId"`
	}

	GroupAddedEvent struct {
		Id string `json:"id"`
	}

	GroupRemovedEvent struct {
		Id   string       `json:"id"`
		Data AccountGroup `json:"data"`
	}
)
