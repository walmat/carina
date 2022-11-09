package tasks

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"nebula/pkg/logger"
	"os"
	"reflect"
)

type storedGroup struct {
	Name  string       `json:"name"`
	Tasks []storedTask `json:"tasks"`
}

type storedTask struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Mode      string `json:"mode"`
	StaticCtx []byte `json:"ctx"`
}

type exportedState struct {
	Groups  map[string]storedGroup `json:"tasks"`
	ToGroup map[string]string      `json:"toGroup"`
}

func ImportGob(outPath string) error {
	data, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(data))

	var state exportedState
	if err = dec.Decode(&state); err != nil {
		return err
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()
	typesMutex.RLock()
	defer typesMutex.RUnlock()

	for groupId, group := range state.Groups {
		taskGroups[groupId] = TaskGroup{
			Name:  group.Name,
			Tasks: make(map[string]*Controller),
		}
		for _, task := range group.Tasks {
			handler, ok := types[task.Type][task.Mode]
			if !ok {
				logger.Warn("Tried importing unregistered task type:", task.Type)
				continue
			}

			newCtx := reflect.New(handler.staticCtx.Elem()).Interface()
			if err = json.Unmarshal(task.StaticCtx, &newCtx); err != nil {
				logger.Warnf("Could not unmarshal state for task (ID: %s, TYP: %s): %s\n", task.Id, task.Type, err)
				continue
			}

			if err = addTaskWithContext(groupId, task.Id, task.Type, task.Mode, handler, newCtx); err != nil {
				logger.Warnf("Failed to import task (ID: %s, TYP: %s): %s\n", task.Id, task.Type, err)
			}
		}
	}

	taskToGroup = state.ToGroup
	return nil
}

func ExportAllGob(outPath string) error {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	outState := exportedState{
		Groups:  make(map[string]storedGroup),
		ToGroup: taskToGroup,
	}
	for groupId, group := range taskGroups {
		stateGroup := storedGroup{
			Name:  group.Name,
			Tasks: make([]storedTask, len(group.Tasks)),
		}
		for taskId, task := range group.Tasks {
			marshalledCtx, err := json.Marshal(task.Container.StaticCtx)
			if err != nil {
				return err
			}

			stateGroup.Tasks = append(stateGroup.Tasks, storedTask{
				Id:        taskId,
				Type:      task.Type,
				Mode:      task.Mode,
				StaticCtx: marshalledCtx,
			})
		}
		outState.Groups[groupId] = stateGroup
	}

	var outBuf bytes.Buffer
	enc := gob.NewEncoder(&outBuf)
	if err := enc.Encode(outState); err != nil {
		return err
	}

	return os.WriteFile(outPath, outBuf.Bytes(), os.ModePerm)
}

func ImportJson(outPath string) error {
	taskBytes, err := os.ReadFile(outPath)
	if err != nil {
		return err
	}

	var parsedState exportedState
	if err = json.Unmarshal(taskBytes, &parsedState); err != nil {
		return err
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()
	typesMutex.RLock()
	defer typesMutex.RUnlock()

	for groupId, group := range parsedState.Groups {
		taskGroups[groupId] = TaskGroup{
			Name:  group.Name,
			Tasks: make(map[string]*Controller),
		}
		for _, task := range group.Tasks {
			handler, ok := types[task.Type][task.Mode]
			if !ok {
				logger.Warn("Tried importing unregistered task type:", task.Type)
				continue
			}

			newCtx := reflect.New(handler.staticCtx.Elem()).Interface()
			if err = json.Unmarshal(task.StaticCtx, &newCtx); err != nil {
				logger.Warnf("Could not unmarshal state for task (ID: %s, TYP: %s): %s\n", task.Id, task.Type, err)
				continue
			}

			if err = addTaskWithContext(groupId, task.Id, task.Type, task.Mode, handler, newCtx); err != nil {
				logger.Warnf("Failed to import task (ID: %s, TYP: %s): %s\n", task.Id, task.Type, err)
			}
		}
	}

	taskToGroup = parsedState.ToGroup
	return nil
}

func ExportAllJson(outPath string) error {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	state := exportedState{
		Groups:  make(map[string]storedGroup),
		ToGroup: taskToGroup,
	}
	for groupId, group := range taskGroups {
		stateGroup := storedGroup{
			Name:  group.Name,
			Tasks: make([]storedTask, len(group.Tasks)),
		}
		for taskId, task := range group.Tasks {
			marshalledCtx, err := json.Marshal(task.Container.StaticCtx)
			if err != nil {
				return err
			}

			stateGroup.Tasks = append(stateGroup.Tasks, storedTask{
				Id:        taskId,
				Type:      task.Type,
				StaticCtx: marshalledCtx,
			})
		}
		state.Groups[groupId] = stateGroup
	}

	exportedBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, exportedBytes, os.ModePerm)
}
