//go:build windows

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package rand

// The Universal CRT provides only the ISO C rand and srand functions declared
// in rand.go. BSD/POSIX extensions remain in rand_default.go.
const LLGoPackage = "decl"
