package channel

import "testing"

func TestMaybeSend(t *testing.T) {
	ch := make(chan bool, 1)

	if MaybeSend(ch, true) != true {
		t.Error("must be able to send to ch")
	}

	if MaybeSend(ch, true) == true {
		t.Error("must not be able to send to ch")
	}

	<-ch

	if MaybeSend(ch, true) != true {
		t.Error("must be able to send to ch again")
	}

	if MaybeSend[bool](nil, true) == true {
		t.Error("must be able to send to nil chan")
	}
}

func TestMaybeSendPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MaybeSend should panic when the channel is closed")
		}
	}()

	ch := make(chan bool)

	close(ch)

	MaybeSend(ch, true)
}

func TestSafeMaybeSend(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("SafeMaybeSend should NOT panic when the channel is closed")
		}
	}()

	ch := make(chan bool)

	close(ch)

	if SafeMaybeSend(ch, true) {
		t.Error("SafeMaybeSend must not be able to send on a closed channel")
	}
}
