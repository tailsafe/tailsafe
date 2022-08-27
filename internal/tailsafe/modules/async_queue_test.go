package modules

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tailsafe/tailsafe/internal/tailsafe/data"
	"log"
	"testing"
	"time"
)

func TestAsyncQueue(t *testing.T) {
	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test", func() error {
		time.Sleep(time.Second * 2)
		log.Print("test 2 second")
		return nil
	}).AddActionToQueue("test2", func() error {
		time.Sleep(time.Second * 1)
		log.Print("test 1 second")
		return nil
	})

	assert.Nil(t, asyncQueue.WaitActions("test", "test2"))
}

func TestAsyncQueue_WaitAll(t *testing.T) {
	payload := data.NewPayload()
	fn1 := func(p *data.Payload) func() error {
		return func() error {
			p.Set("fn1", "ok")
			return nil
		}
	}
	fn2 := func(p *data.Payload) func() error {
		return func() error {
			p.Set("fn2", "ok")
			return nil
		}
	}

	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test", fn1(payload)).AddActionToQueue("test2", fn2(payload))

	asyncQueue.WaitAll()

	assert.Equal(t, "ok", payload.Get("fn1"))
	assert.Equal(t, "ok", payload.Get("fn2"))
}

func TestAsyncQueue_Error(t *testing.T) {
	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test-error", func() error {
		return errors.New("error")
	})

	assert.Nil(t, asyncQueue.WaitActions("test-error"))
}

func TestAsyncQueue_SetWorkers(t *testing.T) {
	GetAsyncQueue().
		SetWorkers(1).
		AddActionToQueue("test", func() error {
			time.Sleep(time.Second * 1)
			log.Print("test 2 second")
			return nil
		}).
		AddActionToQueue("test2", func() error {
			time.Sleep(time.Second * 2)
			log.Print("test 1 second")
			return nil
		})

	assert.Nil(t, GetAsyncQueue().WaitActions("test", "test2"))

}

func TestAsyncQueue_WaitActionsErr(t *testing.T) {
	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test-error", func() error {
		return errors.New("error")
	})

	assert.NotNil(t, asyncQueue.WaitActions("test-error-not-exist"))
}
