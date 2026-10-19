//go:build linux && amd64 && !android

// Code generated from the libc binding; DO NOT EDIT.

package c

// Iovec describes one buffer for scatter/gather I/O.
type Iovec struct {
	IovBase Pointer
	IovLen  SizeT
}
