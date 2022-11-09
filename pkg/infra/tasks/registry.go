package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"nebula/pkg/infra/eventbus"
	"reflect"
	"sync"

	"github.com/lithammer/shortuuid/v3"
)

var (
	taskGroups  = make(map[string]TaskGroup)
	taskToGroup = make(map[string]string)
	tasksMu     sync.RWMutex

	TaskNotExistErr  = errors.New("tasks: task does not exist")
	GroupNotExistErr = errors.New("tasks: group does not exist")

	types      = make(map[string]Handlers)
	typesMutex = sync.RWMutex{}

	TypeRegisteredErr    = errors.New("tasks: task type already registered")
	TypeNotRegisteredErr = errors.New("tasks: task type is not registered")
	ModeNotRegisteredErr = errors.New("tasks: mode type is not registered")
)

type TaskGroup struct {
	Name  string                 `json:"name"`
	Tasks map[string]*Controller `json:"tasks"`
}

type Controller struct {
	Type      string     `json:"type"`
	Mode      string     `json:"mode"`
	Container *Container `json:"container"`
}

func (e *Controller) IsRunning() bool {
	if e == nil {
		return false
	}
	return e.Container.IsRunning()
}

func (e *Controller) Run(ctx context.Context) error {
	return e.Container.Run(ctx)
}

func (e *Controller) Stop() error {
	return e.Container.Stop()
}

func init() {
	taskGroups["default"] = TaskGroup{Name: "Default", Tasks: make(map[string]*Controller)}
}

// RegisterType is used to register task types to a handler. Types should be unique.
// Will return a TypeRegisteredErr if the type is already registered.
func RegisterType(id string, h Handlers) error {
	typesMutex.RLock()
	if _, ok := types[id]; !ok {
		typesMutex.RUnlock()

		typesMutex.Lock()
		types[id] = h
		typesMutex.Unlock()

		return nil
	} else {
		typesMutex.RUnlock()
		return TypeRegisteredErr
	}
}

// AddTask adds a new task of the specified typeId.
// It will return an error if the specified type is not registered.
func AddTask(groupId, typeId string, modeId string, staticCtxObj interface{}) (string, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if !groupExists(groupId) {
		return "", GroupNotExistErr
	}

	handler, err := getTaskHandler(typeId, modeId)
	if err != nil {
		return "", err
	}

	taskId := shortuuid.New()

	staticCtx := reflect.New(handler.staticCtx.Elem()).Interface()
	staticCtxJson, _ := json.Marshal(staticCtxObj)
	if err = json.Unmarshal(staticCtxJson, staticCtx); err != nil {
		return "", err
	}

	if err = addTaskWithContext(groupId, taskId, typeId, modeId, handler, staticCtx); err != nil {
		return "", err
	}
	taskToGroup[taskId] = groupId
	eventbus.Publish(&TaskAddedEvent{
		Id:      taskId,
		GroupId: groupId,
	})
	return taskId, nil
}

// addTaskWithContext is used to add a task very specifically with a handler and static context.
// It will return an error if the static context does not match the handlers type.
// It does not acquire any locks to the resources it accesses, that is to be handled by the caller.
// It will only return an error if it cannot make a Container for the task.
func addTaskWithContext(groupId, taskId, typeId, modeId string, handler *Handler, staticCtx interface{}) error {
	taskContainer, err := newContainerWithContext(taskId, handler, staticCtx)
	if err != nil {
		return err
	}

	taskGroups[groupId].Tasks[taskId] = &Controller{
		Type:      typeId,
		Mode:      modeId,
		Container: taskContainer,
	}
	return nil
}

func EditTask(taskId, typeId string, modeId string, staticCtxObj interface{}) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if !taskExists(taskId) {
		return TaskNotExistErr
	}

	handler, err := getTaskHandler(typeId, modeId)
	if err != nil {
		return err
	}

	groupId := taskToGroup[taskId]

	if taskGroups[groupId].Tasks[taskId].IsRunning() {
		return errors.New("tasks: cannot edit task while it is running")
	}

	staticCtx := reflect.New(handler.staticCtx.Elem()).Interface()
	staticCtxJson, _ := json.Marshal(staticCtxObj)
	if err = json.Unmarshal(staticCtxJson, staticCtx); err != nil {
		return err
	}

	if err = addTaskWithContext(groupId, taskId, typeId, modeId, handler, staticCtx); err != nil {
		return err
	}

	eventbus.Publish(&TaskModifiedEvent{
		Id:      taskId,
		GroupId: groupId,
	})
	return nil
}

// RemoveTask will completely remove a task from the registry by id.
// It will return a TaskNotExistErr if the task does not exist.
func RemoveTask(taskId string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		task := taskGroups[groupId].Tasks[taskId]
		if task.IsRunning() {
			_ = task.Stop()
		}
		delete(taskGroups[groupId].Tasks, taskId)
		delete(taskToGroup, taskId)
		eventbus.Publish(&TaskRemovedEvent{
			Id:      taskId,
			GroupId: groupId,
		})
		return nil
	} else {
		return TaskNotExistErr
	}
}

// RunTask runs a task with the context.Background context, it will not be able to be cancelled.
// It will return a TaskNotExistErr if the task does not exist.
func RunTask(taskId string) error {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		task := taskGroups[groupId].Tasks[taskId]
		return task.Run(context.Background())
	} else {
		return TaskNotExistErr
	}
}

// RunTaskWithContext runs a task with the specified context, if it expires or a cancel function is called the task
// will be cancelled whenever the task returns from whatever state it was in. The task execution cannot
// be stopped until it returns from it's current state.
// It will return a TaskNotExistErr if the task does not exist.
func RunTaskWithContext(taskId string, ctx context.Context) error {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		task := taskGroups[groupId].Tasks[taskId]
		return task.Run(ctx)
	} else {
		return TaskNotExistErr
	}
}

// StopTask handles the call to a task contexts context.CancelFunc.
// It will return an error if tasks cancelFunc is nil, meaning the task is not running.
// It will return a TaskNotExistErr if the task does not exist.
func StopTask(taskId string) error {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		task := taskGroups[groupId].Tasks[taskId]
		return task.Stop()
	} else {
		return TaskNotExistErr
	}
}

func IsRunning(taskId string) bool {
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	if groupId, ok := taskToGroup[taskId]; ok {
		return taskGroups[groupId].Tasks[taskId].IsRunning()
	} else {
		return false
	}
}

// TaskExists returns true if a task with the specified id taskExists, and false if it does not.
func TaskExists(taskId string) bool {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	return taskExists(taskId)
}

func taskExists(taskId string) bool {
	_, ok := taskToGroup[taskId]
	return ok
}

// GetTasks gets all the task id's in the registry and returns them in a slice.
func GetTasks() []string {
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	var tasksList []string
	for _, group := range taskGroups {
		for taskId := range group.Tasks {
			tasksList = append(tasksList, taskId)
		}
	}
	return tasksList
}

func GetTaskType(taskId string) (string, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		return taskGroups[groupId].Tasks[taskId].Type, nil
	} else {
		return "", TaskNotExistErr
	}
}

func GetTaskMode(taskId string) (string, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		return taskGroups[groupId].Tasks[taskId].Mode, nil
	} else {
		return "", TaskNotExistErr
	}
}

func getTaskHandler(taskType string, modeId string) (*Handler, error) {
	typesMutex.RLock()
	defer typesMutex.RUnlock()

	if _, ok := types[taskType]; !ok {
		return nil, TypeNotRegisteredErr
	}

	if _, ok := types[taskType][modeId]; !ok {
		return nil, ModeNotRegisteredErr
	}

	return types[taskType][modeId], nil
}

func GetTypes() []string {
	var arrTypes []string
	for typeId := range types {
		arrTypes = append(arrTypes, typeId)
	}
	return arrTypes
}

func GetModes(typ string) []string {
	var arrModes []string
	for typeId := range types[typ] {
		arrModes = append(arrModes, typeId)
	}
	return arrModes
}

func AddGroup(name string) (string, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	groupId := shortuuid.New()
	taskGroups[groupId] = TaskGroup{
		Name:  name,
		Tasks: make(map[string]*Controller),
	}
	eventbus.Publish(&GroupAddedEvent{
		Id: groupId,
	})
	return groupId, nil
}

func GetGroup(groupId string) (*TaskGroup, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if !groupExists(groupId) {
		return nil, GroupNotExistErr
	}

	group := taskGroups[groupId]
	return &group, nil
}

func GetGroups() []string {
	var groups []string
	for groupId := range taskGroups {
		groups = append(groups, groupId)
	}
	return groups
}

func RemoveGroup(groupId string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if !groupExists(groupId) {
		return GroupNotExistErr
	}

	taskGroup := taskGroups[groupId]
	for taskId := range taskGroup.Tasks {
		delete(taskToGroup, taskId)
	}

	delete(taskGroups, groupId)
	eventbus.Publish(&GroupRemovedEvent{
		Id:   groupId,
		Data: taskGroup,
	})
	return nil
}

func GetTaskGroup(taskId string) (string, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		return groupId, nil
	} else {
		return "", TaskNotExistErr
	}
}

func GroupExists(groupId string) bool {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	return groupExists(groupId)
}

func groupExists(groupId string) bool {
	_, ok := taskGroups[groupId]
	return ok
}

func MoveTask(taskId, newGroupId string) error {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	if !groupExists(newGroupId) {
		return GroupNotExistErr
	} else if !taskExists(taskId) {
		return TaskNotExistErr
	}

	originalGroupId := taskToGroup[taskId]
	controller, ok := taskGroups[originalGroupId].Tasks[taskId]
	if !ok {
		return TaskNotExistErr
	}
	delete(taskGroups[originalGroupId].Tasks, taskId)
	taskGroups[newGroupId].Tasks[taskId] = controller
	taskToGroup[taskId] = newGroupId
	return nil
}
