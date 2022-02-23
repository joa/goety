package bind

import (
	"testing"
)

func TestModule_Get(t *testing.T) {
	type T struct {
		Foo string `bind:"foo"`
		Bar string `bind:""`
		Baz string
	}

	m := newBindings(nil)

	err := m.configure([]Binding{
		Type[*T](),
		Instance[string]("foo").For("foo"),
		Instance[string]("bar"),
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = m.get(typeOf[*T](), "")

	if err != nil {
		t.Fatal(err)
	}
}
