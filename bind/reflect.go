package bind

import (
	"fmt"
	"reflect"
)

// typeOf returns the type of T.
func typeOf[T any]() reflect.Type {
	// HACK: Use a pointer here so that it works for interfaces too. 🤡
	//       We'll then "dereference" the pointer by using the element
	//       of the reflect.Type instance.
	var zero *T
	return reflect.TypeOf(zero).Elem()
}

// alloc an instance of u for type t.
func alloc(t reflect.Type) (v reflect.Value, err error) {
	switch t.Kind() {
	case reflect.Interface:
		err = fmt.Errorf("%w: %s", ErrUnsatisfiedInterface, t)
		return
	case reflect.Pointer:
		// In this case u is *Type so new(*Type) yields **Type but we want *Type.
		v = reflect.New(t.Elem())
	default:
		v = reflect.New(t)
	}

	return
}

// unboxValue v of type t.
//
// It is assumed that v is always a pointer to a value
// whereas the type t not necessarily.
func unboxValue[T any](t reflect.Type, v reflect.Value) T {
	switch t.Kind() {
	case reflect.Pointer,
		reflect.Interface:
		return v.Interface().(T)
	default:
		return *v.Interface().(*T)
	}
}

// assignableTo is true when B is assignable to A.
func assignableTo[A, B any]() (ok bool) {
	// TODO(joa): revisit once we have support for a constraint like [A any, B ~A]
	return typeOf[B]().AssignableTo(typeOf[A]())
}

// mustBeAssignable panics if B can't be assigned to A.
func mustBeAssignable[A, B any]() {
	if !assignableTo[A, B]() {
		panic(fmt.Errorf("can't assign %s to %s", typeOf[B](), typeOf[A]()))
	}
}
