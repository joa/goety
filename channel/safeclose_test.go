package channel

import "testing"

func TestSafeClose(t *testing.T) {
	if SafeClose[any](nil) {
		t.Errorf("SafeClose must not succeed on nil")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("channel wasn't closed")
		}
	}()

	ch := make(chan bool, 1)

	if !SafeClose(ch) {
		t.Errorf("SafeClose must succeed")
	}

	select {
	case ch <- true:
		t.Error("must not be able to send")
	default:
	}

	if SafeClose(ch) {
		t.Errorf("SafeClose must not succeed twice")
	}
}

func TestSafeCloseTwice(t *testing.T) {
	ch := make(chan bool, 1)

	SafeClose(ch)

	if SafeClose(ch) {
		t.Errorf("SafeClose must not succeed twice")
	}
}

func TestSafeCloseNoPanic(t *testing.T) {
	ch := make(chan bool)

	close(ch)

	SafeClose(ch) // must not panic
}
