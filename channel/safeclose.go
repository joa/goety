package channel

// SafeClose the given channel.
//
// Won't panic if the given channel is already closed.
func SafeClose[T any](ch chan T) (ok bool) {
	if ch == nil {
		return
	}

	defer recoverToBool(&ok, false)

	ok = true
	close(ch)

	return
}
