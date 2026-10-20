//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package c_test

import (
	"github.com/goplus/lib/c"
	"testing"
	"unsafe"
)

func TestGeneratedMemoryCalls(t *testing.T) {
	input := [6]byte{0, 1, 127, 128, 254, 255}
	var output [6]byte
	if p := c.Memcpy(unsafe.Pointer(&output[0]), unsafe.Pointer(&input[0]), 6); p != unsafe.Pointer(&output[0]) {
		t.Fatal("memcpy returned another buffer")
	}
	if output != input || c.Memcmp(unsafe.Pointer(&output[0]), unsafe.Pointer(&input[0]), 6) != 0 {
		t.Fatal("native memory copy differs")
	}
	c.Memset(unsafe.Pointer(&output[0]), 0, 6)
	if output != [6]byte{} {
		t.Fatalf("memset result %v", output)
	}
	value := c.Str("native libc")
	if got := c.Strlen(value); got != 11 {
		t.Fatalf("strlen = %d, want 11", got)
	}
}
func TestGeneratedStdioCalls(t *testing.T) {
	stream := c.Tmpfile()
	if stream == nil {
		t.Fatal("tmpfile returned nil")
	}
	defer func() {
		if result := c.Fclose(stream); result != 0 {
			t.Errorf("fclose = %d", result)
		}
	}()
	input := [6]byte{0, 1, 127, 128, 254, 255}
	var output [6]byte
	if got := c.Fwrite(unsafe.Pointer(&input[0]), 1, 6, stream); got != 6 {
		t.Fatalf("fwrite = %d", got)
	}
	if got := c.Fseek(stream, 0, 0); got != 0 {
		t.Fatalf("fseek = %d", got)
	}
	if got := c.Fread(unsafe.Pointer(&output[0]), 1, 6, stream); got != 6 {
		t.Fatalf("fread = %d", got)
	}
	if output != input {
		t.Fatalf("stdio result %v, want %v", output, input)
	}
}
