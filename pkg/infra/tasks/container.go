package tasks

import (
	"context"
	"errors"
	"fmt"
	"nebula/pkg/infra/eventbus"
	"reflect"
	"sync/atomic"
)

var (
	AlreadyRunningErr = errors.New("already running")
	NotRunningErr     = errors.New("task not running")
)

// Container is used to manage all information for a task. It contains it's static context, handler,
// context.CancelFunc, and manages the running state of the task.
type Container struct {
	taskId string

	StaticCtx interface{} `json:"staticCtx"`

	handler    *Handler
	isRunning  int32
	cancelFunc context.CancelFunc
}

func newContainer(taskId string, handler *Handler) (*Container, error) {
	return &Container{
		taskId:    taskId,
		StaticCtx: reflect.New(handler.staticCtx.Elem()).Interface(),
		handler:   handler,
	}, nil
}

func newContainerWithContext(taskId string, handler *Handler, staticCtx interface{}) (*Container, error) {
	ctxType := reflect.TypeOf(staticCtx)
	if ctxType != handler.staticCtx {
		return nil, fmt.Errorf("expected static context type of %s, got %s", handler.staticCtx.Elem().Name(), ctxType.Elem().Name())
	}

	return &Container{
		taskId:    taskId,
		StaticCtx: staticCtx,
		handler:   handler,
	}, nil
}

// Run handles the running of a task, it will return an AlreadyRunningErr if it is already running.
func (c *Container) Run(ctx context.Context) error {
	if atomic.CompareAndSwapInt32(&c.isRunning, 0, 1) {
		defer atomic.SwapInt32(&c.isRunning, 0)

		ctx, c.cancelFunc = context.WithCancel(ctx)
		defer func() {
			// to prevent context leaking.
			if c.cancelFunc != nil {
				c.cancelFunc()
				c.cancelFunc = nil
			}
		}()

		groupId, err := GetTaskGroup(c.taskId)
		if err != nil {
			return err
		}

		taskCtx := &Context{
			Container: c,
			groupId:   groupId,
			Context:   ctx,
		}
		eventbus.Publish(&TaskStartEvent{TaskId: c.taskId, GroupId: taskCtx.groupId})
		return c.handler.Run(taskCtx, c.StaticCtx)
	} else {
		return AlreadyRunningErr
	}
}

// IsRunning returns a bool in relation to whether or not the task is running.
func (c *Container) IsRunning() bool {
	return atomic.LoadInt32(&c.isRunning) == 1
}

func (c *Container) RunAsync(ctx context.Context) chan error {
	retChan := make(chan error)
	go func() {
		retChan <- c.Run(ctx)
	}()
	return retChan
}

// Stop calls the context.CancelFunc for the task if it is running.
// It will return a NotRunningErr if the task is not running.
func (c *Container) Stop() error {
	if c.cancelFunc == nil {
		return NotRunningErr
	}

	c.cancelFunc()
	return nil
}
