package bind

import (
	"reflect"
	"testing"
)

type iface interface{ Meth() }
type impl struct{}

func (i *impl) Meth() {}

func TestTypeOf(t *testing.T) {
	var x int

	e := reflect.TypeOf(x)
	g := typeOf[int]()

	if !reflect.DeepEqual(e, g) {
		t.Error("expected equal types")
	}

	g = typeOf[string]()

	if reflect.DeepEqual(e, g) {
		t.Error("didn't expected equal types")
	}
}

func TestTypeOfInterface(t *testing.T) {
	i := typeOf[iface]()
	j := typeOf[impl]()

	if j.AssignableTo(i) {
		t.Error("impl MUST NOT be assignable to iface")
	}

	j = typeOf[*impl]()

	if !j.AssignableTo(i) {
		t.Error("*impl MUST be assignable to iface")
	}
}

func TestAssignableTo(t *testing.T) {
	if assignableTo[iface, impl]() {
		t.Error("impl MUST NOT be assignable to iface")
	}

	if !assignableTo[iface, *impl]() {
		t.Error("*impl MUST be assignable to iface")
	}
}
