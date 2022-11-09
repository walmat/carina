package tasks

import (
	"context"
	"image/color"
	"nebula/pkg/api/client/rest"
	"nebula/pkg/api/model"
	"nebula/pkg/infra/eventbus"
	"reflect"
	"runtime"
)

type Context struct {
	*Container
	groupId string

	statusCh chan StatusEvent
	context.Context
}

func (ctx *Context) Identifier() string {
	return ctx.taskId
}

func (ctx *Context) GroupIdentifier() string {
	return ctx.groupId
}

func (ctx *Context) SendStatus(message string) {
	ctx.SendStatusColored(message, defaultStatusColor)
}

func (ctx *Context) SendStatusColored(message string, rgba color.RGBA) {
	eventbus.Publish(&StatusEvent{
		TaskId:  ctx.taskId,
		GroupId: ctx.groupId,
		Message: message,
		Color:   rgba,
	})
}

func (ctx *Context) PushProductUpdate() {
	eventbus.Publish(&ProductUpdateEvent{
		TaskId: ctx.taskId,
	})
}

func getCallerFile() *model.LoggerFile {
	_, file, ln, _ := runtime.Caller(2)

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return &model.LoggerFile{
		Name: file,
		Line: ln,
	}
}

func (ctx *Context) Log(entry model.LoggerEntry) {
	entry.File = getCallerFile()
	go ctx.LogSync(entry)
}

func (ctx *Context) LogSync(entry model.LoggerEntry) {
	if entry.File == nil {
		_, file, ln, _ := runtime.Caller(1)

		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short

		entry.File = &model.LoggerFile{
			Name: file,
			Line: ln,
		}
	}

	if entry.Labels == nil {
		entry.Labels = make(map[string]string)
	}
	entry.Labels["taskId"] = ctx.Identifier()
	entry.Labels["taskType"], _ = GetTaskType(ctx.Identifier())

	for i := 0; i < 5; i++ {
		if err := rest.Log(entry); err == nil {
			break
		}
	}
}

func (ctx *Context) Debug(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Debug,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Info(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Info,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Notice(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Notice,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Warning(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Warning,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Error(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Error,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Critical(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Critical,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Alert(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Alert,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

func (ctx *Context) Emergency(payload interface{}, labels ...map[string]string) {
	var label map[string]string
	if len(labels) == 1 {
		label = labels[0]
	}

	ctx.Log(model.LoggerEntry{
		Severity: model.Emergency,
		Payload:  payload,
		Labels:   label,
		File:     getCallerFile(),
	})
}

type StaticContext interface {
	Metadata() Metadata

	// Validate is meant to verify input fields, if everything is ok it will return a nil slice.
	Validate() []error

	GetProductInfo() *model.CheckoutProduct
}

type InitContext interface {
	Init(*Context) error
}

type InitContextExt interface {
	Init(*Context, interface{}) error
}

type CleanupContext interface {
	Cleanup(*Context)
}

type DefaultData struct {
	Store Store    `json:"store"`
	Sku   string   `json:"sku"`
	Sizes []string `json:"sizes"`
}

type Store struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Modes []Mode `json:"modes"`
}

type Mode struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type Metadata struct {
	DisplayName string `json:"displayName"`

	StoreUrl string  `json:"storeUrl"`
	Stores   []Store `json:"subStores"`

	UsesAccounts bool `json:"usesAccounts"`
}

type FullMetadata struct {
	Metadata  `json:"metadata"`
	InputData []InputData `json:"inputData"`
}

func GetMetadataForTaskType(typeId string, modeId string) (FullMetadata, error) {
	handler, err := getTaskHandler(typeId, modeId)
	if err != nil {
		return FullMetadata{}, err
	}
	partialMeta := reflect.New(handler.staticCtx).Elem().Interface().(StaticContext).Metadata()

	inputData, err := GetInputDataForTaskType(typeId, modeId)
	if err != nil {
		return FullMetadata{}, err
	}

	return FullMetadata{
		Metadata:  partialMeta,
		InputData: inputData,
	}, nil
}

func GetStaticContextForTaskType(typeId string, modeId string) (StaticContext, error) {
	handler, err := getTaskHandler(typeId, modeId)
	if err != nil {
		return nil, err
	}

	return reflect.New(handler.staticCtx).Elem().Interface().(StaticContext), nil
}

func GetStaticContextForTaskId(taskId string) (StaticContext, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	if groupId, ok := taskToGroup[taskId]; ok {
		return taskGroups[groupId].Tasks[taskId].Container.StaticCtx.(StaticContext), nil
	} else {
		return nil, TaskNotExistErr
	}
}
