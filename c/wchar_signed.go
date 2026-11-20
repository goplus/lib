//go:build (linux || darwin || freebsd || netbsd || openbsd || solaris || js || wasip1) && !((linux || freebsd) && (arm || arm64))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c

type WcharT int32
