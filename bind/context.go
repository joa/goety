package bind

import (
	"context"
	"errors"
	"fmt"
)

type contextKey string

const ctxKey contextKey = "github.com/joa/goety/bind"

// fromCtx looks for bindings in the current context.
func fromCtx(ctx context.Context) (b *bindings, loaded bool) {
	b, loaded = ctx.Value(ctxKey).(*bindings)
	return
}

// WithBindings enriches a context with bindings.
func WithBindings(ctx context.Context) (context.Context, Bindings) {
	parent, _ := fromCtx(ctx)
	b := newBindings(parent)
	return context.WithValue(ctx, ctxKey, b), b
}

// New creates and returns an instance of T for U.
//
// Note that New will create U if there is no binding present.
// The instance will be fully initialized.
//
// Therefor you can use New to create initialized instances
// that have no specific binding.
//
// This method will panic if any error is raised. Use TryNew instead
// to work with the error.
//
// Example
//
//  type Foo struct {
//    Foo string `bind:"-"`
//  }
//
//  bindings.Configure(
//    bind.String("foo"))
//
//  foo := bind.New[*Foo](ctx)
//  fmt.Println(foo.Foo) // "foo"
//
func New[T, U any](ctx context.Context) U {
	res, err := TryNew[T, U](ctx)
	if err != nil {
		panic(err)
	}
	return res
}

// TryNew creates and returns an instance of T for U.
//
// See New for more information. TryNew will return any error
// that is raised instead of a panic.
func TryNew[T, U any](ctx context.Context) (res U, err error) {
	b, loaded := fromCtx(ctx)
	t := typeOf[T]()
	u := typeOf[U]()

	if !loaded {
		err = fmt.Errorf("%w: for type %s", ErrContextWithoutBindings, t)
		return
	}

	v, err := b.get(t, "") // new will always search without a scope

	if errors.Is(err, ErrNoSuchBinding) {
		v, err = alloc(u)

		if err != nil {
			return
		}

		err = b.initialize(v.Type(), v)

		if err != nil {
			return
		}
	} else if err != nil {
		return
	}

	res = unboxValue[U](u, v)
	return
}

// Get an instance of V in the current context.
//
// This method panics if there's no instance of V available.
// Use TryGet if you expect bindings to fail. This can be
// the case when bindings are miss-configured or when a Provider
// is used that is expected to produce errors.
func Get[V any](ctx context.Context) V {
	res, err := TryGet[V](ctx)
	if err != nil {
		panic(err)
	}
	return res
}

// TryGet an instance of V or return an error.
func TryGet[V any](ctx context.Context) (res V, err error) {
	return TryFor[V](ctx, "")
}

// For - Get an instance of V for a specific key.
//
// This method panics if there's no instance of V available.
// Use TryFor if you expect bindings to fail. This can be
// the case when bindings are miss-configured or when a Provider
// is used that is expected to produce errors.
func For[V any](ctx context.Context, key string) V {
	res, err := TryFor[V](ctx, key)
	if err != nil {
		panic(err)
	}
	return res
}

// TryFor - TryGet an instance of V for a specific key.
func TryFor[V any](ctx context.Context, key string) (res V, err error) {
	b, loaded := fromCtx(ctx)
	t := typeOf[V]()

	if !loaded {
		err = fmt.Errorf("%w: for type %s", ErrContextWithoutBindings, t)
		return
	}

	v, err := b.get(t, key)

	if err != nil {
		return
	}

	res = unboxValue[V](t, v)

	return
}
