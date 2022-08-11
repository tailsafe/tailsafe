package tailsafe

// ActionInterface is the interface that must be implemented by all actions.
type ActionInterface interface {
	// Configure the action
	Configure() (err ErrActionInterface)

	// Execute the action
	Execute() (err ErrActionInterface)

	// SetPayload sets the Payload for the action
	SetPayload(DataInterface)

	// GetConfig returns the configuration for the action
	GetConfig() any

	// GetResult returns the result for the action
	GetResult() any
}

type ActionGetter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ActionSetter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
