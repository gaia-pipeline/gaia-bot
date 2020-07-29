package providers

import (
	"context"
)

// Listener listens for PR events which the bot subscribed to. If the person is authorized to run a command
// it will execute the responding command.
type Listener interface {
	// Listen listens for PR events.
	Listen(ctx context.Context, handle string, commentURL string, pullURL string) error
}
