//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package os_test

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	llos "github.com/goplus/lib/c/os"
)

func TestWindowsStatLayout(t *testing.T) {
	var value llos.StatT
	if got := unsafe.Sizeof(value); got != 56 {
		t.Fatalf("sizeof(_stat64) = %d, want 56", got)
	}
	if got := unsafe.Alignof(value); got != 8 {
		t.Fatalf("alignof(_stat64) = %d, want 8", got)
	}
}

func TestWindowsGetcwd(t *testing.T) {
	buffer := make([]c.Char, llos.PATH_MAX+1)
	if result := llos.Getcwd(unsafe.Pointer(&buffer[0]), c.Int(len(buffer))); result == nil {
		t.Fatalf("_getcwd failed with errno %d", llos.Errno())
	} else if got := c.GoString(result); got == "" {
		t.Fatal("_getcwd returned an empty path")
	}
}

func TestWindowsFileIO(t *testing.T) {
	path := c.AllocaCStr(t.TempDir() + `\binding.txt`)
	fd := llos.Open(path, llos.O_CREAT|llos.O_TRUNC|llos.O_RDWR, c.Int(0o600))
	if fd < 0 {
		t.Fatalf("_open failed with errno %d", llos.Errno())
	}
	defer llos.Close(fd)
	var stat llos.StatT
	if result := llos.Fstat(fd, &stat); result != 0 {
		t.Fatalf("_fstat64 returned %d (errno %d)", result, llos.Errno())
	}
	if result := llos.Stat(path, &stat); result != 0 {
		t.Fatalf("_stat64 returned %d (errno %d)", result, llos.Errno())
	}

	const value = "native UCRT I/O"
	if result := llos.Write(fd, unsafe.Pointer(c.GoStringData(value)), c.Uint(len(value))); result != c.Int(len(value)) {
		t.Fatalf("_write returned %d, want %d (errno %d)", result, len(value), llos.Errno())
	}
	if result := llos.Lseek(fd, 0, 0); result != 0 {
		t.Fatalf("_lseeki64 returned %d, want 0 (errno %d)", result, llos.Errno())
	}

	buffer := make([]byte, len(value))
	if result := llos.Read(fd, unsafe.Pointer(&buffer[0]), c.Uint(len(buffer))); result != c.Int(len(buffer)) {
		t.Fatalf("_read returned %d, want %d (errno %d)", result, len(buffer), llos.Errno())
	}
	if got := string(buffer); got != value {
		t.Fatalf("file contents = %q, want %q", got, value)
	}
}
