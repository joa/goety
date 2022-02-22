package bind

import (
	"context"
	"fmt"
	"reflect"
)

type contextKey string

const ctxKey contextKey = "github.com/joa/goety/bind"

func fromCtx(ctx context.Context) (m *module, loaded bool) {
	m, loaded = ctx.Value(ctxKey).(*module)
	return
}

func WithBindings(ctx context.Context) (context.Context, Bindings) {
	parent, _ := fromCtx(ctx)
	m := newModule(parent)
	return context.WithValue(ctx, ctxKey, m), m
}

func Must[V any](ctx context.Context) V {
	res, err := Get[V](ctx)
	if err != nil {
		panic(err)
	}
	return res
}

func Get[V any](ctx context.Context) (res V, err error) {
	return For[V](ctx, "")
}

func MustFor[V any](ctx context.Context, key string) V {
	res, err := For[V](ctx, key)
	if err != nil {
		panic(err)
	}
	return res
}

func For[V any](ctx context.Context, key string) (res V, err error) {
	m, loaded := fromCtx(ctx)
	t := typeOf[V]()

	if !loaded {
		err = fmt.Errorf("%w: for type %s", ErrContextWithoutBindings, t)
		return
	}

	v, err := m.get(t, key)

	if err != nil {
		return
	}

	switch t.Kind() {
	case reflect.Pointer, reflect.Interface:
		return v.Interface().(V), nil
	default:
		return *v.Interface().(*V), nil
	}
}
