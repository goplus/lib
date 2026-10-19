//go:build linux && amd64 && !android

// Code generated from the libc binding; DO NOT EDIT.

package sync

// MutexAttr is a mutex attribute object.
type MutexAttr struct {
	_xgo_union [1]uint32
}

// CondAttr is a condition variable attribute object.
type CondAttr struct {
	_xgo_union [1]uint32
}

// Mutex is a mutual exclusion lock.
type Mutex struct {
	_xgo_union [5]uint64
}

// Cond is a condition variable.
type Cond struct {
	_xgo_union [6]uint64
}

// RWLock is a read-write lock.
type RWLock struct {
	_xgo_union [7]uint64
}

// RWLockAttr is a read-write lock attribute object.
type RWLockAttr struct {
	_xgo_union [1]uint64
}
