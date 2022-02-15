package channel

// MaybeSend returns true if it could send msg to the channel; false otherwise.
//
// panics if the channel is closed.
func MaybeSend[T any](ch chan T, msg T) (ok bool) {
	select {
	case ch <- msg:
		return true
	default:
		return false
	}
}

// SafeMaybeSend returns true if it could send msg to the channel; false otherwise.
//
// Won't panic if the channel is closed.
func SafeMaybeSend[T any](ch chan T, msg T) (ok bool) {
	defer recoverToBool(&ok, false)
	return MaybeSend(ch, msg)
}
