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
	var mutex Mutex
	if result := mutex.NativeInit(nil); result != 0 {
		t.Fatalf("mutex_init = %d", result)
	}
	defer func() {
		if result := mutex.NativeDestroy(); result != 0 {
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
	if result := cond.NativeInit(nil); result != 0 {
		t.Fatalf("cond_init = %d", result)
	}
	defer func() {
		if result := cond.NativeDestroy(); result != 0 {
			t.Errorf("cond_destroy = %d", result)
		}
	}()
	if result := cond.NativeSignal(); result != 0 {
		t.Fatalf("cond_signal = %d", result)
	}
	if result := cond.NativeBroadcast(); result != 0 {
		t.Fatalf("cond_broadcast = %d", result)
	}
}
func TestGeneratedRWLock(t *testing.T) {
	var rw RWLock
	if result := rw.NativeInit(nil); result != 0 {
		t.Fatalf("rwlock_init = %d", result)
	}
	defer func() {
		if result := rw.NativeDestroy(); result != 0 {
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
