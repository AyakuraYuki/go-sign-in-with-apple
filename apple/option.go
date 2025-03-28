package apple

type Option func(*client)

// WithUpdatePubkeyFailedHandler allow you to do something when the updater failed to
// fetch the Apple's public key
func WithUpdatePubkeyFailedHandler(fn func()) Option {
	return func(c *client) {
		if fn != nil {
			c.onUpdatePubkeyFailed = fn
		}
	}
}
