//go:build llgo

package py_test

import (
	"testing"

	"github.com/goplus/lib/py"
)

func TestNativeSetConstructor(t *testing.T) {
	py.Initialize()
	set := py.NewSet(nil)
	if set == nil {
		t.Fatal("PySet_New returned nil for an empty set")
	}
	defer set.DecRef()
	if got := set.SetLen(); got != 0 {
		t.Fatalf("empty set length = %d, want 0", got)
	}
	value := py.Long(42)
	if value == nil {
		t.Fatal("PyLong_FromLong returned nil")
	}
	defer value.DecRef()
	if set.SetAdd(value) != 0 || set.SetLen() != 1 || set.SetContains(value) != 1 {
		t.Fatal("native set did not retain its inserted value")
	}
}
