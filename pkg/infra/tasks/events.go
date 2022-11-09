package tasks

import "image/color"

var (
	defaultStatusColor = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
)

type (
	StatusEvent struct {
		TaskId  string     `json:"taskId"`
		GroupId string     `json:"groupId"`
		Message string     `json:"message"`
		Color   color.RGBA `json:"color"`
	}

	TaskAddedEvent struct {
		Id      string `json:"id"`
		GroupId string `json:"groupId"`
	}

	TaskModifiedEvent struct {
		Id      string    `json:"id"`
		GroupId string    `json:"groupId"`
	}

	TaskRemovedEvent struct {
		Id      string `json:"id"`
		GroupId string `json:"groupId"`
	}

	TaskStartEvent struct {
		TaskId  string `json:"taskId"`
		GroupId string `json:"groupId"`
	}

	TaskCompletedEvent struct {
		TaskId  string `json:"taskId"`
		GroupId string `json:"groupId"`
		Error   string `json:"error"`
	}

	TaskCancelledEvent struct {
		TaskId  string `json:"taskId"`
		GroupId string `json:"groupId"`
	}

	ProductUpdateEvent struct {
		TaskId string `json:"taskId"`
	}

	GroupAssignedToTaskEvent struct {
		Id     string `json:"id"`
		TaskId string `json:"taskId"`
	}

	GroupAddedEvent struct {
		Id string `json:"id"`
	}

	GroupRemovedEvent struct {
		Id   string    `json:"id"`
		Data TaskGroup `json:"data"`
	}
)
