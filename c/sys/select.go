//go:build !windows

package sys

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/syscall"
)

const (
	LLGoFiles   = "_wrap/fddef.c"
	LLGoPackage = "link"
)

// -----------------------------------------------------------------------------

type (
	FdSet   = syscall.FdSet
	Timeval = syscall.Timeval
)

//go:linkname FD_ZERO C.llgo_FD_ZERO
func FD_ZERO(fds *FdSet)

//go:linkname FD_SET C.llgo_FD_SET
func FD_SET(fd c.Int, fds *FdSet)

//go:linkname FD_ISSET C.llgo_FD_ISSET
func FD_ISSET(fd c.Int, fds *FdSet) c.Int

// -----------------------------------------------------------------------------

// Check the first NFDS descriptors each in READFDS (if not NULL) for read
// readiness, in WRITEFDS (if not NULL) for write readiness, and in EXCEPTFDS
// (if not NULL) for exceptional conditions.  If TIMEOUT is not NULL, time out
// after waiting the interval specified therein.  Returns the number of ready
// descriptors, or -1 for errors.
//
// This function is a cancellation point and therefore not marked with
// __THROW.
//
//go:linkname Select C.select
func Select(__nfds c.Int, __readfds *FdSet, __writefds *FdSet, __exceptfds *FdSet, __timeout *Timeval) c.Int
