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

func (bs *bindings) configure(bindings []Binding) (err error) {
	bs.mut.Lock()
	defer bs.mut.Unlock()

	for _, b := range bindings {
		if err = bs.configureBinding(b); err != nil {
			return
		}
	}

	// initialize all eager bindings
	for _, b := range bindings {
		if !b.eager() {
			continue
		}

		_, err = bs.solve(b)

		if err != nil {
			return
		}
	}

	return
}

// configureBinding - configure a single binding.
func (bs *bindings) configureBinding(b Binding) (err error) {
	typ := b.typ()
	typeScope, loaded := bs.bindings[typ]

	if !loaded {
		typeScope = make(typeBindings)
		bs.bindings[typ] = typeScope
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

func (bs *bindings) get(t reflect.Type, k string) (res reflect.Value, err error) {
	if k == scopeEmptyDash || k == scopeEmptyWildcard {
		k = ""
	}

	binding, ok := findBinding(bs, t, k)

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
		better, ok := findBinding(bs, to.typTo(), "")

		if binding == better || !ok {
			break
		}

		binding = better
	}

	res, err = bs.solve(binding)

	return
}

func (bs *bindings) solve(b Binding) (res reflect.Value, err error) {
	var init bool

	res, init, err = b.solve()

	if err != nil {
		return
	}

	if init {
		err = bs.initialize(res.Type(), res)
	}

	return
}

func (bs *bindings) initialize(typ reflect.Type, value reflect.Value) (err error) {
	switch typ.Kind() {
	case reflect.Pointer:
		if err = bs.initialize(typ.Elem(), value.Elem()); err != nil {
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

			v, err := bs.get(fieldType.Type, scope)

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
