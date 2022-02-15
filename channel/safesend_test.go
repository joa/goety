package channel

import "testing"

func TestSafeSend(t *testing.T) {
	if SafeSend(nil, "foo") {
		t.Error("SafeSend must not send on nil")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error("SafeSend must not panic")
		}
	}()

	ch := make(chan string, 1)

	if !SafeSend(ch, "foo") {
		t.Error("SafeSend must send")
	}

	if <-ch != "foo" {
		t.Error("SafeSend must send 'foo'")
	}

	close(ch)

	if SafeSend(ch, "bar") {
		t.Error("SafeSend must not be able to send")
	}
}
