package bind

import "errors"

var (
	ErrDuplicate              = errors.New("duplicate binding")     // this exact same binding already exists
	ErrNoSuchBinding          = errors.New("no such binding")       // this binding doesn't exist (when resolving)
	ErrContextWithoutBindings = errors.New("no bindings")           // there are no bindings for the context
	ErrUnsatisfiedInterface   = errors.New("interface unsatisfied") // the interface isn't bound to a concrete instance
)
