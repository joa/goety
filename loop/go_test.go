package loop

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var i int32

	f := func() {
		x := atomic.AddInt32(&i, 1)

		if x == 10 {
			cancel()
		}

		if x%3 == 0 {
			panic("panic")
		}
	}

	Go(ctx, f)

	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}

	if actual := atomic.LoadInt32(&i); actual != 10 {
		t.Errorf("expected %d, got %d", 10, actual)
	}
}

func TestGoErr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	errs := make(chan interface{})

	var i int32

	var s strings.Builder
	var mut sync.Mutex

	done := make(chan bool)

	go func() {
		for x := range errs {
			mut.Lock()
			s.WriteString(x.(string))
			mut.Unlock()
		}

		done <- true
	}()

	GoErr(ctx, errs, func() {
		x := atomic.AddInt32(&i, 1)

		if x == 10 {
			cancel()
		}

		if x%3 == 0 {
			panic(fmt.Sprintf("%d", x))
		}
	})

	select {
	case <-ctx.Done():
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}

	close(errs)

	if actual := atomic.LoadInt32(&i); actual != 10 {
		t.Errorf("expected %d, got %d", 10, actual)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for error channel")
	}

	mut.Lock()
	defer mut.Unlock()

	if actual := s.String(); actual != "369" {
		t.Errorf("expected '%s' got %s", "369", actual)
	}
}
