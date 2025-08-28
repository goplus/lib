package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/syscall"
)

//go:linkname Fdopendir C.fdopendir$INODE64
func Fdopendir(fd c.Int) *DIR

//go:linkname Readdir C.readdir$INODE64
func Readdir(dir *DIR) *syscall.Dirent

//go:linkname Fstatat C.fstatat$INODE64
func Fstatat(dirfd c.Int, path *c.Char, buf *StatT, flags c.Int) c.Int
