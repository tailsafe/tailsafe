package tailsafe

import (
	"github.com/tailsafe/tailsafe/internal/tailsafe"
	tailsafeInterface "github.com/tailsafe/tailsafe/pkg/tailsafe"
)

// New creates a new instance of the engine
func New() tailsafeInterface.EngineInterface {
	return tailsafe.New()
}
