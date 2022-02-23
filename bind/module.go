package bind

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	bindTag            = "bind"
	scopeEmptyDash     = "-"
	scopeEmptyWildcard = "*"
)

type typeBindings map[string]Binding

type moduleBindings map[reflect.Type]typeBindings

// Bindings within a context.
type Bindings interface {
	// Configure bindings in a context.
	//
	// Errors are returned when there are duplicate bindings.
	//
	// Example
	//
	//  bindings.Configure(
	//    bind.Instance[string]("username"),
	//    bind.Instance[string]("password"), // this will result in ErrDuplicate
	//  )
	//
	//  bindings.Configure(
	//    bind.Instance[string]("username").For("username"), // different scopes can be used to bind
	//    bind.Instance[string]("password").For("password"), // multiple values of the same type
	//  )
	Configure(bindings ...Binding) (err error)
}

// bindings is a set of type bindings within a context.
type bindings struct {
	mut      sync.RWMutex
	parent   *bindings
	bindings moduleBindings
}

// newBindings creates and returns an initialized bindings object.
func newBindings(parent *bindings) *bindings {
	return &bindings{
		parent:   parent,
		bindings: make(moduleBindings),
	}
}

func (m *bindings) Configure(bindings ...Binding) (err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	for _, b := range bindings {
		if err = m.configureBinding(b); err != nil {
			return
		}
	}

	return
}

// configureBinding - configure a single binding.
func (m *bindings) configureBinding(b Binding) (err error) {
	typ := b.typ()
	typeScope, loaded := m.bindings[typ]

	if !loaded {
		typeScope = make(typeBindings)
		m.bindings[typ] = typeScope
	}

	scope := b.scope()

	if _, loaded = typeScope[scope]; loaded {
		if scope == "" {
			err = fmt.Errorf("%w: %s", ErrDuplicate, typ)
		} else {
			err = fmt.Errorf(`%w: %s for "%s"`, ErrDuplicate, typ, scope)
		}

		return
	}

	typeScope[scope] = b

	return
}

// findBinding for type t and scope k in b and its parents.
func findBinding(b *bindings, t reflect.Type, k string) (Binding, bool) {
	for bb := b; bb != nil; bb = bb.parent {
		bb.mut.RLock()

		typeScope, loaded := bb.bindings[t]

		if !loaded {
			bb.mut.RUnlock()
			continue
		}

		b, loaded := typeScope[k]
		bb.mut.RUnlock()

		if !loaded {
			continue
		}

		return b, true
	}

	return nil, false
}

func (m *bindings) get(t reflect.Type, k string) (res reflect.Value, err error) {
	if k == scopeEmptyDash || k == scopeEmptyWildcard {
		k = ""
	}

	binding, ok := findBinding(m, t, k)

	if !ok {
		if k == "" {
			err = fmt.Errorf("%w: %s", ErrNoSuchBinding, t)
		} else {
			err = fmt.Errorf(`%w: %s for "%s"`, ErrNoSuchBinding, t, k)
		}

		return
	}

	// find a more concrete binding
	for {
		to, ok := binding.(bindingTo)

		if !ok {
			break
		}

		// TODO(joa): is it correct to search with an empty scope for a more concrete binding?
		better, ok := findBinding(m, to.typTo(), "")

		if binding == better || !ok {
			break
		}

		binding = better
	}

	res, init, err := binding.solve()

	if err != nil {
		return
	}

	if init {
		err = m.initialize(res.Type(), res)
	}

	return
}

func (m *bindings) initialize(typ reflect.Type, value reflect.Value) (err error) {
	switch typ.Kind() {
	case reflect.Pointer:
		if err = m.initialize(typ.Elem(), value.Elem()); err != nil {
			return
		}
	case reflect.Struct:
		numField := typ.NumField()

		for fieldIndex := 0; fieldIndex < numField; fieldIndex++ {
			field := value.Field(fieldIndex)

			if !field.CanSet() {
				continue
			}

			fieldType := typ.Field(fieldIndex)
			scope, inject := fieldType.Tag.Lookup(bindTag)

			if !inject {
				continue
			}

			v, err := m.get(fieldType.Type, scope)

			if err != nil {
				return err
			}

			field.Set(v)
		}
	}

	if init, ok := value.Interface().(Initializer); ok {
		err = init.InitAfter()
	}

	return
}
