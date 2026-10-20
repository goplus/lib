//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package signal

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
)

func TestNativeSignalDataLayouts(t *testing.T) {
	var set SigsetT
	var info SiginfoT
	var value Sigval
	if unsafe.Sizeof(set) != nativeSigsetSize || unsafe.Alignof(set) != nativeSigsetAlign {
		t.Fatalf("sigset_t = %d/%d, want %d/%d", unsafe.Sizeof(set), unsafe.Alignof(set), nativeSigsetSize, nativeSigsetAlign)
	}
	if unsafe.Sizeof(info) != nativeSiginfoSize || unsafe.Alignof(info) != 8 {
		t.Fatalf("siginfo_t = %d/%d, want %d/8", unsafe.Sizeof(info), unsafe.Alignof(info), nativeSiginfoSize)
	}
	if unsafe.Offsetof(info.SiSigno) != 0 || unsafe.Offsetof(info.SiErrno) != 4 || unsafe.Offsetof(info.SiCode) != 8 {
		t.Fatalf("siginfo_t header offsets = %d/%d/%d", unsafe.Offsetof(info.SiSigno), unsafe.Offsetof(info.SiErrno), unsafe.Offsetof(info.SiCode))
	}
	if unsafe.Sizeof(value) != 8 || unsafe.Alignof(value) != 8 {
		t.Fatalf("sigval = %d/%d, want 8/8", unsafe.Sizeof(value), unsafe.Alignof(value))
	}
	*value.XGof_ref_sival_int() = -7
	if got := *(*c.Int)(unsafe.Pointer(&value)); got != -7 {
		t.Fatalf("sigval integer storage = %d, want -7", got)
	}
	var marker uintptr
	*value.XGof_ref_sival_ptr() = unsafe.Pointer(&marker)
	if got := *(*unsafe.Pointer)(unsafe.Pointer(&value)); got != unsafe.Pointer(&marker) {
		t.Fatalf("sigval pointer storage = %p, want %p", got, &marker)
	}
}

func TestNativeSignalConstants(t *testing.T) {
	if SIGUSR1 != nativeUserSignal {
		t.Fatalf("SIGUSR1 = %d, want %d", SIGUSR1, nativeUserSignal)
	}
	if SIGEV_NONE != nativeNoNotification {
		t.Fatalf("SIGEV_NONE = %d, want %d", SIGEV_NONE, nativeNoNotification)
	}
}

func TestNativeSignalSetOperations(t *testing.T) {
	var set SigsetT
	if result := nativeEmptySet(&set); result != 0 {
		t.Fatalf("sigemptyset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 0 {
		t.Fatalf("empty membership = %d, want 0", got)
	}
	if result := nativeAddSet(&set, nativeUserSignal); result != 0 {
		t.Fatalf("sigaddset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 1 {
		t.Fatalf("added membership = %d, want 1", got)
	}
	if result := nativeDeleteSet(&set, nativeUserSignal); result != 0 {
		t.Fatalf("sigdelset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 0 {
		t.Fatalf("deleted membership = %d, want 0", got)
	}
	if result := nativeFillSet(&set); result != 0 {
		t.Fatalf("sigfillset = %d", result)
	}
	if got := nativeSetMember(&set, nativeUserSignal); got != 1 {
		t.Fatalf("filled membership = %d, want 1", got)
	}
}
