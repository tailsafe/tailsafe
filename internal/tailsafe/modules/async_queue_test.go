package modules

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
	}, func(err error) {
		t.Error(err)
	}).AddActionToQueue("test2", func() error {
		time.Sleep(time.Second * 1)
		log.Print("test 1 second")
		return nil
	}, func(err error) {
		t.Error(err)
	})

	assert.Nil(t, asyncQueue.WaitActions("test", "test2"))

}

func TestAsyncQueue_Error(t *testing.T) {
	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test-error", func() error {
		return errors.New("error")
	}, func(err error) {
		assert.Error(t, err)
		log.Print(err)
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
		}, func(err error) {
			t.Error(err)
		}).
		AddActionToQueue("test2", func() error {
			time.Sleep(time.Second * 2)
			log.Print("test 1 second")
			return nil
		}, func(err error) {
			t.Error(err)
		})

	assert.Nil(t, GetAsyncQueue().WaitActions("test", "test2"))

}

func TestAsyncQueue_WaitActionsErr(t *testing.T) {
	asyncQueue := GetAsyncQueue()
	asyncQueue.AddActionToQueue("test-error", func() error {
		return errors.New("error")
	}, func(err error) {
		assert.Error(t, err)
		log.Print(err)
	})

	assert.NotNil(t, asyncQueue.WaitActions("test-error-not-exist"))
}
