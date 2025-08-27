//go:build !darwin
// +build !darwin

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

//go:linkname Stat C.stat
func Stat(path *c.Char, buf *StatT) c.Int

//go:linkname Lstat C.lstat
func Lstat(path *c.Char, buf *StatT) c.Int

//go:linkname Fstat C.fstat
func Fstat(fd c.Int, buf *StatT) c.Int
