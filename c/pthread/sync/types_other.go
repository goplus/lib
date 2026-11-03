//go:build !windows && !(linux && amd64 && !android) && !(linux && arm64 && !android)

package sync

// #include <pthread.h>
import "C"

// MutexAttr is a mutex attribute object.
type MutexAttr C.pthread_mutexattr_t

// CondAttr is a condition variable attribute object.
type CondAttr C.pthread_condattr_t

// Mutex is a mutual exclusion lock.
type Mutex C.pthread_mutex_t

// Cond is a condition variable.
type Cond C.pthread_cond_t

// RWLock is a read-write lock.
type RWLock C.pthread_rwlock_t

// RWLockAttr is a read-write lock attribute object.
type RWLockAttr C.pthread_rwlockattr_t
