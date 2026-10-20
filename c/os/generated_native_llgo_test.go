//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package os

import (
	"github.com/goplus/lib/c"
	"testing"
	"unsafe"
)

func TestGeneratedPipeIO(t *testing.T) {
	var pair [2]c.Int
	if result := NativePipe(&pair[0]); result != 0 {
		t.Fatalf("pipe = %d", result)
	}
	defer func() {
		for _, fd := range pair {
			if result := NativeClose(fd); result != 0 {
				t.Errorf("close = %d", result)
			}
		}
	}()
	input := [6]byte{0, 1, 127, 128, 254, 255}
	var output [6]byte
	for offset := 0; offset < len(input); {
		n := NativeWrite(pair[1], unsafe.Pointer(&input[offset]), c.SizeT(len(input)-offset))
		if n <= 0 {
			t.Fatalf("write = %d", n)
		}
		offset += int(n)
	}
	for offset := 0; offset < len(output); {
		n := NativeRead(pair[0], unsafe.Pointer(&output[offset]), c.SizeT(len(output)-offset))
		if n <= 0 {
			t.Fatalf("read = %d", n)
		}
		offset += int(n)
	}
	if output != input {
		t.Fatalf("pipe data = %v, want %v", output, input)
	}
}
