//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package pthread_test

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/pthread"
	"runtime"
	"testing"
	"unsafe"
)

func TestGeneratedThreadAttributeLayout(t *testing.T) {
	size := uintptr(64)
	keySize := uintptr(4)
	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		size = 56
	}
	if runtime.GOOS == "darwin" {
		keySize = 8
	}
	var attr pthread.Attr
	var key pthread.Key
	if unsafe.Sizeof(attr) != size || unsafe.Alignof(attr) != 8 {
		t.Fatalf("pthread_attr_t = %d/%d, want %d/8", unsafe.Sizeof(attr), unsafe.Alignof(attr), size)
	}
	if unsafe.Sizeof(key) != keySize {
		t.Fatalf("pthread_key_t = %d, want %d", unsafe.Sizeof(key), keySize)
	}
}
func TestGeneratedThreadAttributes(t *testing.T) {
	var probe struct {
		before [8]byte
		attr   pthread.Attr
		after  [8]byte
	}
	for i := range probe.before {
		probe.before[i] = 0xa5
		probe.after[i] = 0xa5
	}
	if result := probe.attr.NativeInit(); result != 0 {
		t.Fatalf("attr_init = %d", result)
	}
	defer func() {
		if result := probe.attr.NativeDestroy(); result != 0 {
			t.Errorf("attr_destroy = %d", result)
		}
	}()
	var size c.SizeT
	if result := probe.attr.NativeGetstacksize(&size); result != 0 || size == 0 {
		t.Fatalf("attr_getstacksize = %d, value %d", result, size)
	}
	const requested = 1 << 20
	if result := probe.attr.NativeSetstacksize(requested); result != 0 {
		t.Fatalf("attr_setstacksize = %d", result)
	}
	if result := probe.attr.NativeGetstacksize(&size); result != 0 || size != requested {
		t.Fatalf("stack size = %d, result %d", size, result)
	}
	for i := range probe.before {
		if probe.before[i] != 0xa5 || probe.after[i] != 0xa5 {
			t.Fatal("attribute operation overwrote a buffer boundary")
		}
	}
}

//llgo:type C
type nativeThreadStart func(unsafe.Pointer) unsafe.Pointer

func TestGeneratedNativeThreadRoundtrip(t *testing.T) {
	var thread pthread.Thread
	var marker uintptr = 123
	start := nativeThreadStart(func(argument unsafe.Pointer) unsafe.Pointer { return argument })
	if result := pthread.NativeCreate(&thread, nil, start, unsafe.Pointer(&marker)); result != 0 {
		t.Fatalf("pthread_create = %d", result)
	}
	var returned unsafe.Pointer
	if result := nativeJoin(thread, &returned); result != 0 {
		t.Fatalf("pthread_join = %d", result)
	}
	if returned != unsafe.Pointer(&marker) {
		t.Fatalf("thread result = %p, want %p", returned, &marker)
	}
}
