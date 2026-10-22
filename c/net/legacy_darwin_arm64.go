//go:build darwin && arm64 && !ios

package net

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Ntohs C.ntohs
func Ntohs(x uint16) uint16

//go:linkname Htons C.htons
func Htons(x uint16) uint16

//go:linkname Ntohl C.ntohl
func Ntohl(x c.Uint) c.Uint

//go:linkname Htonl C.htonl
func Htonl(x c.Uint) c.Uint
