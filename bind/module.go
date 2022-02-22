package bind

import (
	"fmt"
	"reflect"
	"sync"
)

const bindTag = "bind"

type Binding interface {
	For(key string) Binding
	typ() reflect.Type
	scope() string
	solve() (res reflect.Value, init bool, err error)
}

type bindingTo interface {
	typTo() reflect.Type
}

type Bindings interface {
	Configure(bindings ...Binding) (err error)
}

type typeBindings map[string]Binding
type moduleBindings map[reflect.Type]typeBindings

type module struct {
	mut      sync.RWMutex
	parent   *module
	bindings moduleBindings
}

func newModule(parent *module) (res *module) {
	res = &module{
		parent:   parent,
		bindings: make(moduleBindings),
	}

	return
}

func (m *module) Configure(bindings ...Binding) (err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	for _, b := range bindings {
		if err = m.configureBinding(b); err != nil {
			return
		}
	}

	return
}

func (m *module) configureBinding(b Binding) (err error) {
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

func findBinding(m *module, t reflect.Type, k string) (Binding, bool) {
	for mm := m; mm != nil; mm = mm.parent {
		mm.mut.RLock()

		typeScope, loaded := mm.bindings[t]

		if !loaded {
			mm.mut.RUnlock()
			continue
		}

		b, loaded := typeScope[k]
		mm.mut.RUnlock()

		if !loaded {
			continue
		}

		return b, true
	}

	return nil, false
}

func (m *module) get(t reflect.Type, k string) (res reflect.Value, err error) {
	if k == "-" || k == "*" {
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

func (m *module) initialize(typ reflect.Type, value reflect.Value) (err error) {
	switch typ.Kind() {
	case reflect.Pointer:
		return m.initialize(typ.Elem(), value.Elem())
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

	return
}
