package slice

import (
	"testing"
)

func TestCloseAllChannels(t *testing.T) {
	CloseAllChannels[any](nil) // should not panic

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("channel wasn't closed")
		}
	}()

	ch := make(chan bool, 1)

	CloseAllChannels([]chan bool{ch})

	select {
	case ch <- true:
		t.Error("must not be able to send")
	default:
	}
}
