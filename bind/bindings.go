package bind

import (
	"fmt"
	"reflect"
)

// Initializer interface is used to let instances know about their creation.
//
// If a type implements the Initializer interface the InitAfter method
// is called after a type was initialized and bindings have completed.
type Initializer interface {
	// InitAfter bindings happened
	InitAfter() (err error)
}

// Binding represents a binding for a specific type.
type Binding interface {
	// For - Scope this binding for a specific key.
	For(key string) Binding

	// typ is the type of this binding.
	typ() reflect.Type

	// scope of this binding.
	scope() string

	// solve this binding.
	solve() (res reflect.Value, init bool, err error)
}

// bindingTo represents an edge to another type.
type bindingTo interface {
	// typTo is the type this binding points to.
	typTo() reflect.Type
}

// Type - Make T available for bindings.
//
// New instances of T are created when requested.
//
// Interfaces
//
// Note that type bindings are always leafs. They must be concrete types as
// the following is a duplicate binding and therefore an error:
//
//  bindings.Configure(
//   bind.Implementation[Iface, Impl](),
//   bind.Type[Iface]()) // Iface is already bound to Impl
//
// Since Type bindings can't be satisfied if the given type is an interface
// this method panics if T is an interface type.
func Type[T any]() Binding {
	t := typeOf[T]()
	if t.Kind() == reflect.Interface {
		panic(fmt.Errorf("can't satisfy %s", t))
	}
	return &typeBind[T, T]{typeFrom: t, typeTo: t}
}

// Implementation - Make Impl available for Iface.
//
// Impl must be assignable to Iface. This can't be checked at compile
// time at the moment and a runtime panic is thrown if that's not the
// case.
//
// Note that Impl may be an interface itself. Implementation bindings
// can be leafs or resolve further. If they are leafs instances
// of Impl are created when requested.
//
// Example
//
//  bindings.Configure(
//    bind.Implementation[Iface, *Impl]())  // in this case we create new instances of Impl
//
//  bindings.Configure(
//    bind.Implementation[Iface, *Impl](),  // here we resolve to another binding for Impl
//    bind.Instance[*Impl](&Impl{}))        // when Impl is requested for Iface we return this instance
//
//  bindings.Configure(
//    bind.Instance[Iface, *Impl](&Impl{})) // note the above can also be simply an instance bind
//
// TODO(joa): compiler support missing for Implementation[Iface any, Impl ~Iface]
func Implementation[Iface, Impl any]() Binding {
	mustBeAssignable[Iface, Impl]()
	return &typeBind[Iface, Impl]{
		typeFrom: typeOf[Iface](),
		typeTo:   typeOf[Impl](),
	}
}

// Instance - Make an instance of U available for T.
//
// Note this instance is created manually, prior to the existence of the bindings.
//
// Example
//
//  bindings.Configure(
//    bind.Instance[Iface, *Impl](&Impl{}), // bind an instance for another type
//    bind.Instance[string]("value"))       // bind an instance for the same type
func Instance[T, U any](inst U) Binding {
	mustBeAssignable[T, U]()
	return &instBind[T, U]{inst: reflect.ValueOf(inst)}
}

// String - Shortcut for bind.Instance[string]
func String(s string) Binding { return Instance[string](s) }

// Int - Shortcut for bind.Instance[int]
func Int(i int) Binding { return Instance[int](i) }

// Provider - Bind a function f to type T.
func Provider[T any](f func() (T, error)) Binding {
	return &providerBind[T]{f: f}
}

// typeBind represents a bind of a type From to type To
type typeBind[From, To any] struct {
	key      string
	typeFrom reflect.Type
	typeTo   reflect.Type
}

func (b *typeBind[From, To]) typ() reflect.Type   { return b.typeFrom }
func (b *typeBind[From, To]) typTo() reflect.Type { return b.typeTo }
func (b *typeBind[From, To]) scope() string       { return b.key }

func (b *typeBind[From, To]) solve() (value reflect.Value, init bool, err error) {
	value, err = alloc(b.typeTo)
	init = true
	return
}

func (b *typeBind[From, To]) For(k string) Binding {
	b.key = k
	return b
}

type instBind[T, U any] struct {
	key  string
	inst reflect.Value // of U
}

func (b *instBind[T, U]) typ() reflect.Type                   { return typeOf[T]() }
func (b *instBind[T, U]) scope() string                       { return b.key }
func (b *instBind[T, U]) solve() (reflect.Value, bool, error) { return b.inst, false, nil }

func (b *instBind[T, U]) For(k string) Binding {
	b.key = k
	return b
}

type providerBind[T any] struct {
	key string
	f   func() (T, error)
}

func (b *providerBind[T]) typ() reflect.Type { return typeOf[T]() }
func (b *providerBind[T]) scope() string     { return b.key }

func (b *providerBind[T]) solve() (reflect.Value, bool, error) {
	res, err := b.f()
	return reflect.ValueOf(res), true, err
}

func (b *providerBind[T]) For(k string) Binding {
	b.key = k
	return b
}
