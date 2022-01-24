package channel

// SafeSend msg to ch.
//
// Won't panic if the given channel is already closed.
func SafeSend[T any](ch chan T, msg T) (ok bool) {
	if ch == nil {
		return
	}

	defer recoverToBool(&ok, false)

	ok = true
	ch <- msg

	return
}
