package slice

// DeleteInPlaceNoOrder an element of a slice.
//
// Note: modifies s and does not preserve the order of elements in s.
func DeleteInPlaceNoOrder[T comparable](in []T, elem T) (res []T, modified bool) {
	var ts T

	for i, vv := range in {
		if vv == elem {
			// delete without preserving order
			// see https://github.com/golang/go/wiki/SliceTricks
			in[i] = in[len(in)-1]
			in[len(in)-1] = ts
			return in[:len(in)-1], true
		}
	}

	return in, false
}
