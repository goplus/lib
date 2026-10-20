//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package sync

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestGeneratedSynchronizationLayouts(t *testing.T) {
	var mutex Mutex
	var mutexAttr MutexAttr
	var cond Cond
	var condAttr CondAttr
	var rw RWLock
	var rwAttr RWLockAttr
	sizes := []uintptr{40, 4, 48, 4, 56, 8}
	aligns := []uintptr{8, 4, 8, 4, 8, 8}
	if runtime.GOOS == "darwin" {
		sizes = []uintptr{64, 16, 48, 16, 200, 24}
		aligns = []uintptr{8, 8, 8, 8, 8, 8}
	} else if runtime.GOARCH == "arm64" {
		sizes = []uintptr{48, 8, 48, 8, 56, 8}
		aligns = []uintptr{8, 8, 8, 4, 8, 8}
	}
	checks := []struct {
		name        string
		size, align uintptr
	}{{"mutex", unsafe.Sizeof(mutex), unsafe.Alignof(mutex)}, {"mutex attr", unsafe.Sizeof(mutexAttr), unsafe.Alignof(mutexAttr)}, {"cond", unsafe.Sizeof(cond), unsafe.Alignof(cond)}, {"cond attr", unsafe.Sizeof(condAttr), unsafe.Alignof(condAttr)}, {"rwlock", unsafe.Sizeof(rw), unsafe.Alignof(rw)}, {"rwlock attr", unsafe.Sizeof(rwAttr), unsafe.Alignof(rwAttr)}}
	for i, check := range checks {
		if check.size != sizes[i] || check.align != aligns[i] {
			t.Errorf("%s = %d/%d, want %d/%d", check.name, check.size, check.align, sizes[i], aligns[i])
		}
	}
}
func TestGeneratedMutex(t *testing.T) {
	t.Run("attributes", func(t *testing.T) {
		var attr MutexAttr
		if result := attr.Init(); result != 0 {
			t.Fatalf("mutexattr_init = %d", result)
		}
		defer func() {
			if result := attr.Destroy(); result != 0 {
				t.Errorf("mutexattr_destroy = %d", result)
			}
		}()
		if result := attr.SetType(MUTEX_NORMAL); result != 0 {
			t.Fatalf("mutexattr_settype = %d", result)
		}
		var mutex Mutex
		if result := mutex.Init(&attr); result != 0 {
			t.Fatalf("mutex_init with attributes = %d", result)
		}
		if result := mutex.Destroy(); result != 0 {
			t.Fatalf("mutex_destroy with attributes = %d", result)
		}
	})
	var mutex Mutex
	var _ interface {
		Lock()
		Unlock()
	} = &mutex
	if result := mutex.Init(nil); result != 0 {
		t.Fatalf("mutex_init = %d", result)
	}
	defer func() {
		if result := mutex.Destroy(); result != 0 {
			t.Errorf("mutex_destroy = %d", result)
		}
	}()
	if result := mutex.NativeLock(); result != 0 {
		t.Fatalf("mutex_lock = %d", result)
	}
	if result := mutex.NativeUnlock(); result != 0 {
		t.Fatalf("mutex_unlock = %d", result)
	}
}
func TestGeneratedCondition(t *testing.T) {
	var cond Cond
	if result := cond.Init(nil); result != 0 {
		t.Fatalf("cond_init = %d", result)
	}
	defer func() {
		if result := cond.Destroy(); result != 0 {
			t.Errorf("cond_destroy = %d", result)
		}
	}()
	if result := cond.Signal(); result != 0 {
		t.Fatalf("cond_signal = %d", result)
	}
	if result := cond.Broadcast(); result != 0 {
		t.Fatalf("cond_broadcast = %d", result)
	}
}
func TestGeneratedRWLock(t *testing.T) {
	var rw RWLock
	var _ interface {
		Lock()
		Unlock()
		RLock()
		RUnlock()
	} = &rw
	if result := rw.Init(nil); result != 0 {
		t.Fatalf("rwlock_init = %d", result)
	}
	defer func() {
		if result := rw.Destroy(); result != 0 {
			t.Errorf("rwlock_destroy = %d", result)
		}
	}()
	if result := rw.NativeRdlock(); result != 0 {
		t.Fatalf("rwlock_rdlock = %d", result)
	}
	if result := rw.NativeUnlock(); result != 0 {
		t.Fatalf("rwlock_unlock = %d", result)
	}
	if result := rw.NativeWrlock(); result != 0 {
		t.Fatalf("rwlock_wrlock = %d", result)
	}
	if result := rw.NativeUnlock(); result != 0 {
		t.Fatalf("rwlock_unlock = %d", result)
	}
}
