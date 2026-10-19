//go:build (darwin && arm64 && !ios) || (linux && amd64 && !android)

package os

import "github.com/goplus/lib/c/syscall"

// Dirent shares the representation already used by Readdir.
type Dirent = syscall.Dirent
