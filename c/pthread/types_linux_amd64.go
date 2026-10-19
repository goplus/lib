//go:build linux && amd64 && !android

// Code generated from the libc binding; DO NOT EDIT.

package pthread

import "github.com/goplus/lib/c"

// Attr represents POSIX thread attributes.
type Attr struct {
	_xgo_union [7]uint64
}

// AttrT exposes the native attribute mapping on validated targets.
type AttrT = Attr

type Key c.Uint
