//go:build llgo && darwin && arm64 && !ios

package pthread_test

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/pthread"
	"unsafe"
)

// nativeJoin follows the platform-specific method/global shape emitted for pthread_t.
func nativeJoin(thread pthread.Thread, result *unsafe.Pointer) c.Int {
	return thread.NativeJoin(result)
}
