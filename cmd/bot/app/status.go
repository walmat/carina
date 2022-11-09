package app

import (
	"image/color"
	"nebula/cmd/bot/windows"
	"nebula/pkg/infra/eventbus"
	"nebula/pkg/infra/harvesters"
	"nebula/pkg/infra/tasks"
	"sync"
	"time"
)

// TODO(cc): refactor below

type EventType string

const (
	EventTypeStatus        EventType = "status"
	EventTypeTaskStart     EventType = "taskStart"
	EventTypeTaskCancelled EventType = "taskCancelled"
	EventTypeTaskCompleted EventType = "taskCompleted"
	EventTypeProductUpdate EventType = "taskProductUpdate"
)

var (
	eventQueue []FrontendEvent
	eventMutex = sync.RWMutex{}

	statusQueue      = make(map[string]tasks.StatusEvent)
	statusQueueMutex = sync.RWMutex{}
)

type StrippedStatusEvent struct {
	Message string     `json:"message"`
	Color   color.RGBA `json:"color"`
}

type FrontendEvent struct {
	EventType EventType `json:"type"`
	EventData interface{}     `json:"data"`
}

func init() {
	eventbus.Subscribe(TaskStartHandler)
	eventbus.Subscribe(StatusUpdateHandler)
	eventbus.Subscribe(ProductUpdateHandler)
	eventbus.Subscribe(TaskCompletedHandler)
	eventbus.Subscribe(TaskCancelledHandler)

	go func() {
		for range time.Tick(time.Millisecond * 100) {

			statusQueueMutex.RLock()
			if len(statusQueue) == 0 {
				statusQueueMutex.RUnlock()
				continue
			}

			statusEvents := make(map[string]map[string]StrippedStatusEvent)
			for _, event := range statusQueue {
				if statusEvents[event.GroupId] == nil {
					statusEvents[event.GroupId] = make(map[string]StrippedStatusEvent)
				}

				statusEvents[event.GroupId][event.TaskId] = StrippedStatusEvent{
					Message: event.Message,
					Color:   event.Color,
				}
			}
			statusQueueMutex.RUnlock()

			dispatchEvent(EventTypeStatus, statusEvents)

			statusQueueMutex.Lock()
			statusQueue = make(map[string]tasks.StatusEvent)
			statusQueueMutex.Unlock()
		}
	}()

	go func() {
		for range time.Tick(time.Millisecond * 100) {

			eventMutex.RLock()
			if len(eventQueue) == 0 {
				eventMutex.RUnlock()
				continue
			}

			if windows.Windows.MainWindow != nil && windows.Windows.MainWindow.IsShown() {
				windows.Windows.MainWindow.SendMessage(struct {
					Type string      `json:"type"`
					Data interface{} `json:"data"`
				}{
					Type: "taskEvents",
					Data: eventQueue,
				})
			}
			eventMutex.RUnlock()

			eventMutex.Lock()
			eventQueue = nil
			eventMutex.Unlock()
		}
	}()
}

func dispatchEvent(eventType EventType, data interface{}) {
	eventMutex.Lock()
	eventQueue = append(eventQueue, FrontendEvent{
		EventType: eventType,
		EventData: data,
	})
	eventMutex.Unlock()
}

var (
	lastStatuses      = make(map[string]tasks.StatusEvent)
	lastStatusesMutex = sync.RWMutex{}
)

func clearLastStatus(taskId string) {
	lastStatusesMutex.Lock()
	delete(lastStatuses, taskId)
	lastStatusesMutex.Unlock()
}

func StatusUpdateHandler(event *tasks.StatusEvent) {
	lastStatusesMutex.RLock()
	if lastStatuses[event.TaskId] == *event {
		lastStatusesMutex.RUnlock()
		return
	}
	lastStatusesMutex.RUnlock()

	lastStatusesMutex.Lock()
	lastStatuses[event.TaskId] = *event
	lastStatusesMutex.Unlock()

	statusQueueMutex.Lock()
	statusQueue[event.TaskId] = *event
	statusQueueMutex.Unlock()
}

func TaskStartHandler(event *tasks.TaskStartEvent) {
	clearLastStatus(event.TaskId)
	dispatchEvent(EventTypeTaskStart, *event)
}

func TaskCancelledHandler(event *tasks.TaskCancelledEvent) {
	harvesters.RemoveSolve(event.TaskId)
	clearLastStatus(event.TaskId)
	dispatchEvent(EventTypeTaskCancelled, *event)
}

func TaskCompletedHandler(event *tasks.TaskCompletedEvent) {
	harvesters.RemoveSolve(event.TaskId)
	clearLastStatus(event.TaskId)
	dispatchEvent(EventTypeTaskCompleted, *event)
}

func ProductUpdateHandler(event *tasks.ProductUpdateEvent) {
	dispatchEvent(EventTypeProductUpdate, *event)
}
