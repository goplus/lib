//go:build (linux || freebsd) && (arm || arm64)

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package c

// The Linux and FreeBSD ARM ABIs use unsigned int for wchar_t.
type WcharT uint32
