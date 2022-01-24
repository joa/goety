package dispatcher

type Dispatcher[T any] interface {
	// NewRecv creates and registers a receiver for this dispatcher.
	NewRecv() (recv <-chan T, err error)

	// DeleteRecv deletes and unregisters a receiver of this dispatcher.
	DeleteRecv(recv <-chan T) (err error)

	// Dispatch a value that is observed by all receivers.
	Dispatch(v T)
}
