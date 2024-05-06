package lang

func If[T any](cond bool, yes, no T) T {
	if cond {
		return yes
	} else {
		return no
	}
}

func IfLazy[T any](cond bool, yes, no func() T) T {
	if cond {
		return yes()
	} else {
		return no()
	}
}
