package bind_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joa/goety/bind"
)

type Iface interface {
	Meth() string
}

type Impl struct {
	Field string
}

func (d *Impl) Meth() string { return "Impl-" + d.Field }

func TestGet(t *testing.T) {
	type T struct {
		X  string `bind:"x"` // bind strings with scope x
		Y  string `bind:""`  // bind strings without scope
		DB Iface  `bind:"-"` // bind:"-" is the same as bind:""
		Z  string // won't inject this as there is no bind annotation
		w  string `bind:"x"` // can't set this as it is private
	}

	// create bindings for context
	ctx, bindings := bind.WithBindings(context.Background())

	// configure bindings
	err := bindings.Configure(
		bind.Instance[string]("X").For("x"),  // bind string with scope x
		bind.Instance[string]("Y"),           // bind string without scope
		bind.Implementation[Iface, *Impl](),  // bind Iface to *Impl
		bind.Instance[*Impl](&Impl{"Field"}), // bind *Impl to concrete instance
		bind.Type[*T](),                      // bind *T (make it available)
		bind.Type[T](),                       // bind T which shouldn't be an error as it's a different type
	)

	if err != nil {
		t.Fatal(err)
		return
	}

	if res, err := bind.Get[*Impl](ctx); err != nil {
		// we should be able to get *Impl
		t.Fatal(err)
		return
	} else if res.Field != "Field" {
		t.Errorf("expected Field, got %s", res.Field)
	}

	if res, err := bind.Get[T](ctx); err != nil {
		// we should be able to get T
		t.Fatal(err)
		return
	} else if res.X != "X" {
		t.Errorf("expected X, got %s", res.X)
	} else if res.Y != "Y" {
		t.Errorf("expected Y, got %s", res.X)
	} else if res.Z != "" {
		t.Errorf("expected '', got %s", res.X)
	}

	// get a *T from the context
	res, err := bind.Get[*T](ctx)

	if err != nil {
		t.Fatal(err)
		return
	}

	if act := res.DB.Meth(); act != "Impl-Field" {
		t.Errorf("expected Impl-Field, got %s", res.DB.Meth())
	}

	if res.X != "X" {
		t.Errorf("expected X, got %s", res.X)
	}

	if res.Y != "Y" {
		t.Errorf("expected Y, got %s", res.X)
	}

	if res.Z != "" {
		t.Errorf("expected '', got %s", res.X)
	}
}

type IfaceA interface {
	A()
}

type IfaceB interface {
	IfaceA
	B()
}
type ImplAB struct{}

func (i *ImplAB) A() {}
func (i *ImplAB) B() {}

func TestGetIfaceToIfaceUnsatisfied(t *testing.T) {
	ctx, bindings := bind.WithBindings(context.Background())

	// configure bindings but leave IfaceB unsatisfied
	err := bindings.Configure(
		bind.Implementation[IfaceA, IfaceB](),
		//bind.Implementation[IfaceB, ???](),
	)

	if err != nil {
		t.Fatal(err)
		return
	}

	if _, err = bind.Get[IfaceA](ctx); !errors.Is(err, bind.ErrUnsatisfiedInterface) {
		t.Errorf("expected ErrUnsatisfiedInterface, got %s", err)
	}
}

func TestGetIfaceToIfaceSatisfied(t *testing.T) {
	ctx, bindings := bind.WithBindings(context.Background())

	err := bindings.Configure(
		bind.Implementation[IfaceA, IfaceB](),
		bind.Implementation[IfaceB, *ImplAB](),
	)

	if err != nil {
		t.Fatal(err)
		return
	}

	if _, err = bind.Get[IfaceA](ctx); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}
