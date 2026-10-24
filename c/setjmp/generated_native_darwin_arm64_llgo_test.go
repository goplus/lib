//go:build llgo && darwin && arm64 && !ios

package setjmp_test

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/setjmp"
	"testing"
)

func TestGeneratedSetjmpWritesCallerBuffer(t *testing.T) {
	var buffer [48]c.Int
	if result := setjmp.NativeSetjmp(&buffer[0]); result != 0 {
		t.Fatalf("initial setjmp = %d", result)
	}
	changed := false
	for _, word := range buffer {
		changed = changed || word != 0
	}
	if !changed {
		t.Fatal("setjmp did not write the caller-owned buffer")
	}
}
