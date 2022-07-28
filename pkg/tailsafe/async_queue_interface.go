package tailsafe

// AsyncQueue is the interface that must be implemented by all async queues.
type AsyncQueue interface {

	// AddActionToQueue Add an action to the asynchronous queue
	AddActionToQueue(action string, call func() error) AsyncQueue

	// WaitActions Allows you to wait for the end of the specified actions
	WaitActions(actions ...string) error

	// WaitAll Allows you to wait for the end of all actions
	WaitAll()

	// SetWorkers sets the number of workers for the asynchronous queue
	SetWorkers(num int) AsyncQueue
}
