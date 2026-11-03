//go:build !windows

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

type DIR struct {
	Unused [0]byte
}

//go:linkname Opendir C.opendir
func Opendir(name *c.Char) *DIR
