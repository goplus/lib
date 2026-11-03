//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package sys

import (
	"testing"
	"unsafe"
)

func TestGeneratedGettimeofday(t *testing.T) {
	var timeout Timeval
	if result := Select(0, nil, nil, nil, &timeout); result != 0 {
		t.Fatalf("select with zero timeout = %d", result)
	}
	var value Timeval
	if unsafe.Sizeof(value) != 16 || unsafe.Alignof(value) != 8 {
		t.Fatalf("timeval layout = %d/%d", unsafe.Sizeof(value), unsafe.Alignof(value))
	}
	if result := Gettimeofday(&value, nil); result != 0 {
		t.Fatalf("gettimeofday = %d", result)
	}
	if value.Sec <= 0 || value.Usec < 0 || value.Usec >= 1000000 {
		t.Fatalf("invalid timeval %v", value)
	}
}
