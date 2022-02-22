package bind

import (
	"fmt"
	"reflect"
)

func mustBeAssignable[A, B any]() {
	if !assignableTo[A, B]() {
		panic(fmt.Errorf("can't assign %s to %s", typeOf[B](), typeOf[A]()))
	}
}

func Type[T any]() Binding {
	a := typeOf[T]()
	return &typeBind[T, T]{typA: a, typB: a}
}

func Implementation[Iface, Impl any]() Binding {
	//TODO(joa): compiler support not yet; need func Implementation[Iface any, Impl Iface]() Binding {
	mustBeAssignable[Iface, Impl]()
	return &typeBind[Iface, Impl]{typA: typeOf[Iface](), typB: typeOf[Impl]()}
}

func Instance[A, B any](inst B) Binding {
	mustBeAssignable[A, B]()
	return &instBind[A, B]{inst: reflect.ValueOf(inst)}
}

func Provider[A any](f func() (A, error)) Binding {
	return &providerBind[A]{f: f}
}

// typeBind represents a bind of a type A to a new instance of B
type typeBind[A, B any] struct {
	key  string
	typA reflect.Type
	typB reflect.Type
}

func (b *typeBind[A, B]) typ() reflect.Type   { return b.typA }
func (b *typeBind[A, B]) typTo() reflect.Type { return b.typB }
func (b *typeBind[A, B]) scope() string       { return b.key }

func (b *typeBind[A, B]) solve() (reflect.Value, bool, error) {
	var value reflect.Value

	switch b.typB.Kind() {
	case reflect.Interface:
		return value, false, fmt.Errorf("%w: %s", ErrUnsatisfiedInterface, b.typB)
	case reflect.Pointer:
		value = reflect.New(b.typB.Elem())
	default:
		value = reflect.New(b.typB)
	}

	return value, true, nil
}

func (b *typeBind[A, B]) For(k string) Binding {
	b.key = k
	return b
}

type instBind[A, B any] struct {
	key  string
	inst reflect.Value // of B
}

func (b *instBind[A, B]) typ() reflect.Type                   { return typeOf[A]() }
func (b *instBind[A, B]) scope() string                       { return b.key }
func (b *instBind[A, B]) solve() (reflect.Value, bool, error) { return b.inst, false, nil }

func (b *instBind[A, B]) For(k string) Binding {
	b.key = k
	return b
}

type providerBind[A any] struct {
	key string
	f   func() (A, error)
}

func (b *providerBind[A]) typ() reflect.Type { return typeOf[A]() }
func (b *providerBind[A]) scope() string     { return b.key }

func (b *providerBind[A]) solve() (reflect.Value, bool, error) {
	res, err := b.f()
	return reflect.ValueOf(res), true, err
}

func (b *providerBind[A]) For(k string) Binding {
	b.key = k
	return b
}
