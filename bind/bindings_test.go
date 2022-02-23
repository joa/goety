package bind_test

import (
	"context"
	"testing"

	"github.com/joa/goety/bind"
)

var counter int

type onceInst struct{}

func (oi *onceInst) InitAfter() (err error) {
	counter += 1
	return
}

func TestOnce(t *testing.T) {
	ctx, err := bind.Configure(context.Background(), bind.Once[*onceInst]())

	if err != nil {
		t.Fatal(err)
	}

	if counter != 1 {
		t.Error("onceInst wasn't eagerly initialized")
	}

	if bind.Get[*onceInst](ctx) != bind.New[*onceInst](ctx) {
		t.Error("two different instances of once")
	}

	if bind.For[*onceInst](ctx, "") != bind.New[*onceInst](ctx) {
		t.Error("two different instances of once")
	}

	if counter != 1 {
		t.Error("onceInst initializer called more than once")
	}
}
