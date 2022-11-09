package eventbus

import (
	"errors"
	"reflect"
	"sync"
)

var (
	listeners   = make(map[reflect.Type][]reflect.Value)
	listenersMu sync.RWMutex

	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

func Subscribe(handler interface{}) error {
	handlerVal := reflect.ValueOf(handler)
	handlerType := handlerVal.Type()
	if handlerType.Kind() != reflect.Func {
		return errors.New("eventbus: handlers must be functions")
	}

	if handlerType.NumIn() != 1 {
		return errors.New("eventbus: handlers must take 1 argument, the event data they are listening for")
	}

	inType := handlerType.In(0)
	if inType.Kind() != reflect.Ptr {
		return errors.New("eventbus: handlers must take a pointer to the event data they are listening for")
	}

	if handlerType.NumOut() > 1 || (handlerType.NumOut() == 1 && !handlerType.Out(1).Implements(errorInterface)) {
		return errors.New("eventbus: handlers must return nothing or an error type")
	}

	listenersMu.Lock()
	defer listenersMu.Unlock()
	if handlers, ok := listeners[inType]; ok {
		handlers = append(handlers, handlerVal)
		listeners[inType] = handlers
	} else {
		listeners[inType] = []reflect.Value{handlerVal}
	}

	return nil
}

func Publish(event interface{}) {
	go PublishSync(event)
}

func PublishSync(event interface{}) (interface{}, error) {
	eventVal := reflect.ValueOf(event)
	eventType := eventVal.Type()

	listenersMu.RLock()
	defer listenersMu.RUnlock()

	if handlers, ok := listeners[eventType]; ok {
		for _, handler := range handlers {
			out := handler.Call([]reflect.Value{eventVal})
			if len(out) == 1 {
				err := out[0].Interface().(error)
				if err != nil {
					return eventVal.Interface(), err
				}
			}
		}
	}

	return eventVal.Interface(), nil
}
