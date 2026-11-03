//go:build !windows && !((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package os

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Closedir C.closedir
func Closedir(dir *DIR) c.Int

//go:linkname Close C.close
func Close(fd c.Int) c.Int

//go:linkname Fchdir C.fchdir
func Fchdir(dirfd c.Int) c.Int

//go:linkname Dup C.dup
func Dup(fd c.Int) c.Int

//go:linkname Dup2 C.dup2
func Dup2(oldfd c.Int, newfd c.Int) c.Int

//go:linkname Fsync C.fsync
func Fsync(fd c.Int) c.Int

//go:linkname Isatty C.isatty
func Isatty(fd c.Int) c.Int

// In Mach, all system calls take normal arguments and always return an
// error code (zero for success).
//
//go:linkname Syscall C.syscall
func Syscall(sysno c.Long, __llgo_va_list ...any) c.Long

//go:linkname Fcntl C.fcntl
func Fcntl(a c.Int, b c.Int, __llgo_va_list ...any) c.Int
