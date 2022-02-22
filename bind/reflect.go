package bind

import "reflect"

// typeOf returns the type of T.
func typeOf[T any]() reflect.Type {
	// use a pointer here so that it works for interfaces too 🤡
	var zero *T
	return reflect.TypeOf(zero).Elem()
}

// assignableTo is true when B is assignable to A.
func assignableTo[A, B any]() (ok bool) {
	// TODO(joa): revisit once we have support for a constraint like [A any, B ~A]
	var zeroA *A
	var zeroB *B
	return reflect.TypeOf(zeroB).Elem().AssignableTo(reflect.TypeOf(zeroA).Elem())
}
