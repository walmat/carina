// Package example NOT COMPLETE
package example

import (
	"context"
	"sync"
	"sync/atomic"
)

var (
	watchers      = make(map[string]chan Product)
	watchersMutex = sync.RWMutex{}

	isRunning  int32
	cancelFunc context.CancelFunc
)

type Product struct {
	Name      string
	ProductId string
}

func QueueWatcher(identifier string) chan Product {
	watchersMutex.Lock()
	defer watchersMutex.Unlock()

	notifChan := make(chan Product)
	watchers[identifier] = notifChan

	if atomic.CompareAndSwapInt32(&isRunning, 0, 1) {
		ctx, ctxCancelFunc := context.WithCancel(context.Background())
		cancelFunc = ctxCancelFunc
		go handler(ctx)
	}

	return notifChan
}

func DequeueWatcher(identifier string) {
	watchersMutex.Lock()
	delete(watchers, identifier)
	watchersMutex.Unlock()

	if len(watchers) == 0 && cancelFunc != nil {
		cancelFunc()
	}

	atomic.SwapInt32(&isRunning, 0)
}

func handler(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := tick(ctx); err != nil {
			// TODO: handle error
		}
	}
}

func tick(ctx context.Context) error {
	return nil
}
