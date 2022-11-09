package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	rpcFuncs      = make(map[string]reflect.Value)
	rpcFuncsMutex sync.RWMutex

	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

func RegisterRPC(name string, fn interface{}) error {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()
	if fnType.Kind() != reflect.Func {
		return errors.New("rpc: handler must be a function")
	}

	numOut := fnType.NumOut()
	if numOut == 2 {
		if !fnType.Out(1).Implements(errorInterface) {
			return errors.New("rpc: handler can only return an error and/or response type")
		}
	} else if numOut > 2 {
		return errors.New("rpc: handler can only return an error and/or response type")
	}

	rpcFuncsMutex.Lock()
	rpcFuncs[name] = fnValue
	rpcFuncsMutex.Unlock()
	return nil
}

func handleRPC(name string, rawArgs []interface{}) (interface{}, error) {
	rpcFuncsMutex.RLock()
	method, ok := rpcFuncs[name]
	if !ok {
		return nil, fmt.Errorf("rpc: call to unregistered handler `%s`", name)
	}

	isVariadic := method.Type().IsVariadic()
	numIn := method.Type().NumIn()

	if (isVariadic && len(rawArgs) < numIn-1) || (!isVariadic && len(rawArgs) != numIn) {
		return nil, errors.New("function arguments mismatch")
	}
	var args []reflect.Value
	for i := range rawArgs {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(method.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(method.Type().In(i))
		}

		// handle undefined values from js
		argValue := reflect.ValueOf(rawArgs[i])
		if !argValue.IsValid() {
			arg.Elem().Set(reflect.Zero(arg.Type()))
		} else {
			if arg.Elem().Type().AssignableTo(argValue.Type()) {
				arg.Elem().Set(argValue)
			} else {
				jsonData, err := json.Marshal(argValue.Interface())
				if err != nil {
					return nil, err
				}

				if err = json.Unmarshal(jsonData, arg.Interface()); err != nil {
					return nil, err
				}
			}
		}
		args = append(args, arg.Elem())
	}
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := method.Call(args)
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil
	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil
	case 2:
		// Two results: first one is value, second is error
		if !res[1].Type().Implements(errorType) {
			return nil, errors.New("second return value must be an error")
		}
		if res[1].Interface() == nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)
	default:
		return nil, errors.New("unexpected number of return values")
	}
}
