package slice

import "testing"

func TestDeleteInPlaceNoOrder(t *testing.T) {
	var ts any

	if _, ok := DeleteInPlaceNoOrder[any](nil, ts); ok {
		t.Error("delete on empty slice must not modify")
	}

	if _, ok := DeleteInPlaceNoOrder[int]([]int{1, 2, 3, 4}, 5); ok {
		t.Error("delete slice \\wo element must not modify")
	}

	if s, ok := DeleteInPlaceNoOrder[int]([]int{1, 2, 5, 3, 4}, 5); !ok {
		t.Error("delete in slice with element must modify")

		if len(s) != 4 {
			t.Error("expected length 4")
		}
	}
}
