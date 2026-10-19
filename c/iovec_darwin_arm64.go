//go:build darwin && arm64 && !ios

// Code generated from the libc binding; DO NOT EDIT.

package c

// Iovec describes one buffer for scatter/gather I/O.
type Iovec struct {
	IovBase Pointer
	IovLen  Ulong
}
