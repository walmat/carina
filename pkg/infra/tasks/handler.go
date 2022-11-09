// TODO: more comments
package tasks

import (
	"context"
	"errors"
	"fmt"
	"nebula/pkg/infra/eventbus"
	"reflect"
)

var (
	CompletedErr = errors.New("task complete")

	stringType             = reflect.TypeOf("")
	errorInterface         = reflect.TypeOf((*error)(nil)).Elem()
	staticContextInterface = reflect.TypeOf((*StaticContext)(nil)).Elem()
	taskContextType        = reflect.TypeOf(&Context{})
)

type HandlerMap map[interface{}]interface{}

type reflectMap map[interface{}]reflect.Value

type Handlers map[string]*Handler

type Handler struct {
	staticCtx     reflect.Type
	runningCtx    reflect.Type
	initialState  interface{}
	stateHandlers reflectMap
}

func NewHandler(initialState interface{}, stateHandlers HandlerMap) (*Handler, error) {
	var initialStateType, staticCtxType, runningCtxType reflect.Type = reflect.TypeOf(initialState), nil, nil
	if !initialStateType.ConvertibleTo(stringType) {
		return nil, fmt.Errorf("expected state type that is convertible to string, type %s is not", initialStateType.Name())
	}

	handlerMap := make(reflectMap)
	for state, handler := range stateHandlers {
		stateValue := reflect.ValueOf(state)
		if stateValue.IsZero() {
			return nil, errors.New("cannot have zero state value")
		}

		stateType := stateValue.Type()
		if stateType == stringType {
			// use enum-like types you monkey
			return nil, fmt.Errorf("state type must be convertible to a string, not a literal string")
		} else if stateType != initialStateType {
			// cannot mix and match state types
			return nil, fmt.Errorf("expected state type of %s got %s", initialStateType.Name(), stateType.Name())
		}

		handlerValue := reflect.ValueOf(handler)
		if handlerValue.Kind() != reflect.Func {
			// handler has to be a function
			return nil, fmt.Errorf("expected handler to be a function type, it is a %d", handlerValue.Kind())
		}

		handlerType := handlerValue.Type()
		if handlerType.NumIn() != 3 {
			return nil, fmt.Errorf("expected handler function to accept 2 arguments (static ctx, running ctx), it only takes %d", handlerType.NumIn())
		}

		// first type will be a task context
		firstInType := handlerType.In(0)
		if firstInType != taskContextType {
			return nil, fmt.Errorf("expected first type to be a %s, it was %s", taskContextType.Name(), firstInType.Name())
		}

		// first type will be static context type
		secondInType := handlerType.In(1)
		if staticCtxType == nil {
			if secondInType.Kind() != reflect.Ptr {
				return nil, errors.New("static context type must be pointer")
			} else if !secondInType.Implements(staticContextInterface) {
				return nil, errors.New("static context type must implement static context type")
			}

			staticCtxType = secondInType
		} else if staticCtxType != secondInType {
			return nil, fmt.Errorf("expected first type to be a %s, it was %s", staticCtxType.Name(), secondInType.Name())
		}

		// second type will be running context type
		thirdInType := handlerType.In(2)
		if runningCtxType == nil {
			if secondInType.Kind() != reflect.Ptr {
				return nil, errors.New("running context type must be pointer")
			}

			runningCtxType = thirdInType
		} else if runningCtxType != thirdInType {
			return nil, fmt.Errorf("expected second type to be a %s, it was %s", staticCtxType.Name(), thirdInType.Name())
		}

		if handlerType.NumOut() == 1 {
			// must return state
			firstOutType := handlerType.Out(0)
			if handlerType.Out(0) != stateType {
				return nil, fmt.Errorf("handler function must return %s or (%s, error), not %s", stateType.Name(), stateType.Name(), firstOutType.Name())
			}
		} else if handlerType.NumOut() == 2 {
			// must return (state, error)
			firstOutType := handlerType.Out(0)
			secondOutType := handlerType.Out(1)
			if firstOutType != stateType {
				return nil, fmt.Errorf("handler function must return %s or (%s, error), not (%s, %s) (first return type problematic)", stateType.Name(), stateType.Name(), firstOutType.Name(), secondOutType.Name())
			}

			if !handlerType.Out(1).Implements(errorInterface) {
				return nil, fmt.Errorf("handler function must return %s or (%s, error), not (%s, %s) (second return type problematic)", stateType.Name(), stateType.Name(), firstOutType.Name(), secondOutType.Name())
			}
		} else {
			return nil, fmt.Errorf("expected function to have one (state) or two (state, error) return types, it has %d", handlerType.NumOut())
		}

		handlerMap[state] = handlerValue
	}

	return &Handler{
		initialState:  initialState,
		staticCtx:     staticCtxType,
		runningCtx:    runningCtxType,
		stateHandlers: handlerMap,
	}, nil
}

func (h *Handler) callHandler(state interface{}, taskCtx *Context, staticCtx interface{}, runningCtx reflect.Value) ([]reflect.Value, error) {
	if handler, ok := h.stateHandlers[state]; ok {
		return handler.Call([]reflect.Value{reflect.ValueOf(taskCtx), reflect.ValueOf(staticCtx), runningCtx}), nil
	} else {
		return nil, fmt.Errorf("could not find handler for state %s", state)
	}
}

func (h *Handler) Run(taskCtx *Context, staticCtx interface{}) error {
	staticCtxType := reflect.TypeOf(staticCtx)
	if staticCtxType != h.staticCtx {
		return fmt.Errorf("expected state type %s, got type %s", h.staticCtx.Name(), staticCtxType.Name())
	} else if staticCtxInit, ok := staticCtx.(InitContext); ok {
		if err := staticCtxInit.Init(taskCtx); err != nil {
			return errors.New("failed to initialize static context")
		}
	}
	if staticCtxUninit, ok := staticCtx.(CleanupContext); ok {
		defer staticCtxUninit.Cleanup(taskCtx)
	}

	runningCtx := reflect.New(h.runningCtx.Elem())
	if runningCtxInitExt, ok := runningCtx.Interface().(InitContextExt); ok {
		if err := runningCtxInitExt.Init(taskCtx, staticCtx); err != nil {
			return errors.New("failed to initialize running context")
		}
	} else if runningCtxInit, ok := runningCtx.Interface().(InitContext); ok {
		if err := runningCtxInit.Init(taskCtx); err != nil {
			return errors.New("failed to initialize running context")
		}
	}
	if runningCtxUninit, ok := runningCtx.Interface().(CleanupContext); ok {
		defer runningCtxUninit.Cleanup(taskCtx)
	}

	currentState := h.initialState
	for {
		select {
		case <-taskCtx.Done():
			eventbus.Publish(&TaskCancelledEvent{TaskId: taskCtx.Identifier(), GroupId: taskCtx.GroupIdentifier()})
			return context.Canceled
		default:
			retVal, err := h.callHandler(currentState, taskCtx, staticCtx, runningCtx)
			if err != nil {
				eventbus.Publish(&TaskCompletedEvent{
					TaskId:  taskCtx.Identifier(),
					GroupId: taskCtx.GroupIdentifier(),
					Error:   err.Error(),
				})
				return err
			}

			currentState = retVal[0].Interface()

			if len(retVal) == 2 {
				retErrorInterface := retVal[1].Interface()
				if retErrorInterface != nil {
					if retErrorInterface == CompletedErr {
						eventbus.Publish(&TaskCompletedEvent{
							TaskId:  taskCtx.Identifier(),
							GroupId: taskCtx.GroupIdentifier(),
							Error:   "",
						})
						return nil
					} else {
						retError := retErrorInterface.(error)
						eventbus.Publish(&TaskCompletedEvent{
							TaskId:  taskCtx.Identifier(),
							GroupId: taskCtx.GroupIdentifier(),
							Error:   retError.Error(),
						})
						return retError
					}
				}
			}
		}
	}
}
