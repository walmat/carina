package app

import (
	"nebula/pkg/infra/profiles"
	"nebula/pkg/infra/proxies"
	"nebula/pkg/infra/tasks"
)

type Identifier struct {
	Id string `json:"id"`
}

type TaskGroupMap map[string][]string

func AddTasks(typeId string, modeId string, staticCtx interface{}, groups []Identifier, profileId string, proxyGroup string, amount int) (TaskGroupMap, error) {
	taskMap := make(TaskGroupMap)
	for _, groupId := range groups {
		tasksMade := 0
		for tasksMade < amount {
			taskId, err := tasks.AddTask(groupId.Id, typeId, modeId, staticCtx)
			if err != nil {
				return nil, err
			}

			if err = profiles.AssignProfileToTask(taskId, profileId); err != nil {
				_ = tasks.RemoveTask(taskId)
				return nil, err
			} else if proxyGroup != "" {
				if err = proxies.AssignProxyGroupToTask(taskId, proxyGroup); err != nil {
					_ = tasks.RemoveTask(taskId)
					return nil, err
				}
			}

			taskMap[groupId.Id] = append(taskMap[groupId.Id], taskId)

			tasksMade++
		}
	}
	return taskMap, nil
}

func StartTasks(ids []string) {
	go func() {
		for _, taskId := range ids {
			startId := taskId
			go func() {
				_ = tasks.RunTask(startId)
			}()
		}
	}()
}

func StopTasks(ids []string) {
	go func() {
		for _, taskId := range ids {
			_ = tasks.StopTask(taskId)
		}
	}()
}

func RemoveTasks(taskIds []string) error {
	for _, taskId := range taskIds {
		if err := tasks.RemoveTask(taskId); err != nil {
			return err
		}
	}
	return nil
}

func MoveTasks(taskIds []string, newGroupId string) error {
	for _, taskId := range taskIds {
		if err := tasks.MoveTask(taskId, newGroupId); err != nil {
			return err
		}
	}
	return nil
}
