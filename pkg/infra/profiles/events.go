package profiles

type (
	ProfileAddedEvent struct {
		Id      string      `json:"id"`
		GroupId string      `json:"groupId"`
		Data    ProfileData `json:"data"`
	}

	ProfileRemovedEvent struct {
		Id      string `json:"id"`
		GroupId string `json:"groupId"`
	}

	ProfileModifiedEvent struct {
		Id   string      `json:"id"`
		Data ProfileData `json:"data"`
	}

	ProfileAssignedToTaskEvent struct {
		Id     string `json:"id"`
		TaskId string `json:"taskId"`
	}

	GroupAddedEvent struct {
		Id string `json:"id"`
	}

	GroupRemovedEvent struct {
		Id   string       `json:"id"`
		Data ProfileGroup `json:"data"`
	}
)
