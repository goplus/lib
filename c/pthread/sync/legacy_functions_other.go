//go:build !windows && !((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package sync

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/time"
	_ "unsafe"
)

// llgo:link (*Once).Do C.pthread_once
func (o *Once) Do(f func()) c.Int { return 0 }

// llgo:link (*MutexAttr).Init C.pthread_mutexattr_init
func (a *MutexAttr) Init(attr *MutexAttr) c.Int { return 0 }

// llgo:link (*MutexAttr).Destroy C.pthread_mutexattr_destroy
func (a *MutexAttr) Destroy() {}

// llgo:link (*MutexAttr).SetType C.pthread_mutexattr_settype
func (a *MutexAttr) SetType(typ MutexType) c.Int { return 0 }

// llgo:link (*Mutex).Init C.pthread_mutex_init
func (m *Mutex) Init(attr *MutexAttr) c.Int { return 0 }

// llgo:link (*Mutex).Destroy C.pthread_mutex_destroy
func (m *Mutex) Destroy() {}

// llgo:link (*Mutex).TryLock C.pthread_mutex_trylock
func (m *Mutex) TryLock() c.Int { return 0 }

// llgo:link (*RWLockAttr).Init C.pthread_rwlockattr_init
func (a *RWLockAttr) Init(attr *RWLockAttr) c.Int { return 0 }

// llgo:link (*RWLockAttr).Destroy C.pthread_rwlockattr_destroy
func (a *RWLockAttr) Destroy() {}

// llgo:link (*RWLockAttr).SetPShared C.pthread_rwlockattr_setpshared
func (a *RWLockAttr) SetPShared(pshared c.Int) c.Int { return 0 }

// llgo:link (*RWLockAttr).GetPShared C.pthread_rwlockattr_getpshared
func (a *RWLockAttr) GetPShared(pshared *c.Int) c.Int { return 0 }

// llgo:link (*RWLock).Init C.pthread_rwlock_init
func (rw *RWLock) Init(attr *RWLockAttr) c.Int { return 0 }

// llgo:link (*RWLock).Destroy C.pthread_rwlock_destroy
func (rw *RWLock) Destroy() {}

// llgo:link (*RWLock).TryRLock C.pthread_rwlock_tryrdlock
func (rw *RWLock) TryRLock() c.Int { return 0 }

// llgo:link (*RWLock).TryLock C.pthread_rwlock_trywrlock
func (rw *RWLock) TryLock() c.Int { return 0 }

// llgo:link (*CondAttr).Init C.pthread_condattr_init
func (a *CondAttr) Init(attr *CondAttr) c.Int { return 0 }

// llgo:link (*CondAttr).Destroy C.pthread_condattr_destroy
func (a *CondAttr) Destroy() {}

// llgo:link (*Cond).Init C.pthread_cond_init
func (c *Cond) Init(attr *CondAttr) c.Int { return 0 }

// llgo:link (*Cond).Destroy C.pthread_cond_destroy
func (c *Cond) Destroy() {}

// llgo:link (*Cond).Signal C.pthread_cond_signal
func (c *Cond) Signal() c.Int { return 0 }

// llgo:link (*Cond).Broadcast C.pthread_cond_broadcast
func (c *Cond) Broadcast() c.Int { return 0 }

// llgo:link (*Cond).Wait C.pthread_cond_wait
func (c *Cond) Wait(m *Mutex) c.Int { return 0 }

// llgo:link (*Cond).TimedWait C.pthread_cond_timedwait
func (c *Cond) TimedWait(m *Mutex, abstime *time.Timespec) c.Int { return 0 }
