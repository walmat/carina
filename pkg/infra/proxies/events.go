package proxies

type (
	ProxyAddedEvent struct {
		Id      string    `json:"id"`
		GroupId string    `json:"groupId"`
		Data    ProxyData `json:"data"`
	}

	ProxyRemovedEvent struct {
		Id      string `json:"id"`
		GroupId string `json:"groupId"`
	}

	ProxyModifiedEvent struct {
		Id   string    `json:"id"`
		Data ProxyData `json:"data"`
	}

	GroupAssignedToTaskEvent struct {
		Id     string `json:"id"`
		TaskId string `json:"taskId"`
	}

	GroupAddedEvent struct {
		Id string `json:"id"`
	}

	GroupRemovedEvent struct {
		Id   string     `json:"id"`
		Data ProxyGroup `json:"data"`
	}
)
