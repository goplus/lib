//go:build !windows && !((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package pthread

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

// The pthread_exit() function terminates the calling thread and
// returns a value via retval that (if the thread is joinable) is
// available to another thread in the same process that calls
// pthread_join(3).
//
// See https://man7.org/linux/man-pages/man3/pthread_exit.3.html
//
//go:linkname Exit C.pthread_exit
func Exit(retval c.Pointer)

// The pthread_cancel() function sends a cancelation request to the
// thread thread.
//
// See https://man7.org/linux/man-pages/man3/pthread_cancel.3.html
//
//go:linkname Cancel C.pthread_cancel
func Cancel(thread Thread) c.Int

// llgo:link (*Attr).Init C.pthread_attr_init
func (attr *Attr) Init() c.Int { return 0 }

// llgo:link (*Attr).Destroy C.pthread_attr_destroy
func (attr *Attr) Destroy() c.Int { return 0 }

// llgo:link (*Attr).GetDetached C.pthread_attr_getdetachstate
func (attr *Attr) GetDetached(detached *c.Int) c.Int { return 0 }

// llgo:link (*Attr).SetDetached C.pthread_attr_setdetachstate
func (attr *Attr) SetDetached(detached c.Int) c.Int { return 0 }

// llgo:link (*Attr).GetStackSize C.pthread_attr_getstacksize
func (attr *Attr) GetStackSize(stackSize *uintptr) c.Int { return 0 }

// llgo:link (*Attr).SetStackSize C.pthread_attr_setstacksize
func (attr *Attr) SetStackSize(stackSize uintptr) c.Int { return 0 }

// llgo:link (*Attr).GetStackAddr C.pthread_attr_getstackaddr
func (attr *Attr) GetStackAddr(stackAddr *c.Pointer) c.Int { return 0 }

// llgo:link (*Attr).SetStackAddr C.pthread_attr_setstackaddr
func (attr *Attr) SetStackAddr(stackAddr c.Pointer) c.Int { return 0 }

// llgo:link (*Key).Create C.pthread_key_create
func (key *Key) Create(destructor func(c.Pointer)) c.Int { return 0 }

// llgo:link Key.Delete C.pthread_key_delete
func (key Key) Delete() c.Int { return 0 }

// llgo:link Key.Get C.pthread_getspecific
func (key Key) Get() c.Pointer { return nil }

// llgo:link Key.Set C.pthread_setspecific
func (key Key) Set(value c.Pointer) c.Int { return __noop__() }
