//go:build llgo && linux && (amd64 || arm64) && !android

package setjmp_test

import (
	"github.com/goplus/lib/c/setjmp"
	"testing"
	"unsafe"
)

func TestGeneratedSetjmpWritesCallerBuffer(t *testing.T) {
	var buffer setjmp.X__JmpBufTag
	if result := buffer.NativeSetjmp(); result != 0 {
		t.Fatalf("initial setjmp = %d", result)
	}
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&buffer)), unsafe.Sizeof(buffer))
	changed := false
	for _, b := range bytes {
		changed = changed || b != 0
	}
	if !changed {
		t.Fatal("setjmp did not write the caller-owned buffer")
	}
}
