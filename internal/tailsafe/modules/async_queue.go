package modules

import (
	"fmt"
	"github.com/tailsafe/tailsafe/pkg/tailsafe"
	"golang.org/x/exp/slices"
	"runtime"
	"sync"
)

const (
	CLASSIC_NEXT      = iota
	NEXT_AFTER_FINISH = 1
)

type AsyncQueue struct {
	worker  int
	current int
	queue   []*Action

	stateAction map[string]*Action
	nextTick    chan int
	sync.Mutex
}

type Action struct {
	name string
	call func() error
	done chan bool
}

var asyncQueueInstance *AsyncQueue

func init() {
	asyncQueueInstance = new(AsyncQueue)
	asyncQueueInstance.worker = runtime.NumCPU()
	asyncQueueInstance.stateAction = make(map[string]*Action)
	asyncQueueInstance.nextTick = make(chan int)

	asyncQueueInstance._Run()
}

// GetAsyncQueue returns the async queue instance
func GetAsyncQueue() tailsafe.AsyncQueue {
	return asyncQueueInstance
}

// AddActionToQueue adds a new action to the queue
func (a *AsyncQueue) AddActionToQueue(action string, call func() error) tailsafe.AsyncQueue {
	pendingAction := &Action{action, call, make(chan bool, 1)}

	a.Lock()
	a.queue = append(a.queue, pendingAction)
	a.stateAction[pendingAction.name] = pendingAction
	a.Unlock()

	a._NextTick(CLASSIC_NEXT)

	return a
}

func (a *AsyncQueue) SetWorkers(num int) tailsafe.AsyncQueue {
	a.worker = num

	return a
}

// Run starts the queue
func (a *AsyncQueue) _Run() {
	go func() {
		for typeTick := range a.nextTick {
			if typeTick == NEXT_AFTER_FINISH {
				a.current--
			}

			a.Lock()
			total := len(a.queue)
			a.Unlock()

			if total == 0 {
				continue
			}

			if a.current >= a.worker {
				continue
			}

			task := a.queue[0]
			a.queue = slices.Delete[[]*Action](a.queue, 0, 1)
			a.current++

			go func() {
				_ = task.call()

				task.done <- true

				a._NextTick(NEXT_AFTER_FINISH)
			}()
		}
	}()
}

func (a *AsyncQueue) _NextTick(typeTick int) {
	a.nextTick <- typeTick
}

// WaitActions waits for an action to be done
func (a *AsyncQueue) WaitActions(actions ...string) (err error) {
	var isolateWait []*Action
	for _, action := range actions {
		a.Lock()
		task, ok := a.stateAction[action]
		a.Unlock()

		if !ok {
			return fmt.Errorf("action %s not found", action)
		}
		isolateWait = append(isolateWait, task)
	}

	for _, action := range isolateWait {
		<-action.done
	}

	return
}

func (a *AsyncQueue) WaitAll() {
	a.Lock()
	var isolateWait []*Action
	for _, action := range a.stateAction {
		isolateWait = append(isolateWait, action)
	}
	a.Unlock()

	for _, action := range isolateWait {
		<-action.done
	}
	return
}
