//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package net

import "github.com/goplus/lib/c"

// AddrInfo matches the Winsock ADDRINFOA layout. ai_addrlen is size_t on
// Windows, unlike the socklen_t field used by Unix implementations.
type AddrInfo struct {
	Flags     c.Int
	Family    c.Int
	SockType  c.Int
	Protocol  c.Int
	AddrLen   uintptr
	CanOnName *c.Char
	Addr      *SockAddr
	Next      *AddrInfo
}
