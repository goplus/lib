//go:build darwin && arm64 && !ios

// Code generated from the libc binding; DO NOT EDIT.

package pthread

import "github.com/goplus/lib/c"

// Attr represents POSIX thread attributes.
type Attr struct {
	X__Sig    c.Long
	X__Opaque [56]c.Char
}

// AttrT exposes the native attribute mapping on validated targets.
type AttrT = Attr

type Key c.Ulong
