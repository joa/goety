package slice

// CloseAllChannels in the given slice.
func CloseAllChannels[T any](s []chan T) {
	for _, v := range s {
		close(v)
	}
}
