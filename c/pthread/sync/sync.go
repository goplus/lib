//go:build !windows

/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sync

// #include <pthread.h>
import "C"

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/time"
)

const (
	LLGoFiles   = "_wrap/pthd.c"
	LLGoPackage = "link"
)

// -----------------------------------------------------------------------------

// Once is an object that will perform exactly one action.
type Once C.pthread_once_t

//go:linkname OnceInit llgoSyncOnceInitVal
var OnceInit Once

// -----------------------------------------------------------------------------

type MutexType c.Int

const (
	MUTEX_NORMAL     MutexType = C.PTHREAD_MUTEX_NORMAL
	MUTEX_ERRORCHECK MutexType = C.PTHREAD_MUTEX_ERRORCHECK
	MUTEX_RECURSIVE  MutexType = C.PTHREAD_MUTEX_RECURSIVE
	MUTEX_DEFAULT    MutexType = C.PTHREAD_MUTEX_DEFAULT
)

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------

// llgo:link (*Mutex).Lock C.wrap_pthread_mutex_lock
func (m *Mutex) Lock() {}

// llgo:link (*Mutex).Unlock C.wrap_pthread_mutex_unlock
func (m *Mutex) Unlock() {}

// llgo:link (*RWLock).RLock C.pthread_rwlock_rdlock
func (rw *RWLock) RLock() {}

// llgo:link (*RWLock).RUnlock C.pthread_rwlock_unlock
func (rw *RWLock) RUnlock() {}

// llgo:link (*RWLock).Lock C.pthread_rwlock_wrlock
func (rw *RWLock) Lock() {}

// llgo:link (*RWLock).Unlock C.pthread_rwlock_unlock
func (rw *RWLock) Unlock() {}

// llgo:link (*CondAttr).SetClock C.pthread_condattr_setclock
func (a *CondAttr) SetClock(clock time.ClockidT) c.Int { return 0 }

// llgo:link (*CondAttr).GetClock C.pthread_condattr_getclock
func (a *CondAttr) GetClock(clock *time.ClockidT) c.Int { return 0 }
