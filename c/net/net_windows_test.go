//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package net

import (
	"testing"
	"unsafe"
)

func TestWinsockLayouts(t *testing.T) {
	ptrSize := unsafe.Sizeof(uintptr(0))
	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{"SockaddrIn", unsafe.Sizeof(SockaddrIn{}), 16},
		{"SockaddrIn6", unsafe.Sizeof(SockaddrIn6{}), 28},
		{"SockaddrStorage", unsafe.Sizeof(SockaddrStorage{}), 128},
		{"SockAddr", unsafe.Sizeof(SockAddr{}), 16},
		{"Hostent.AddrType", unsafe.Offsetof(Hostent{}.AddrType), 2 * ptrSize},
		{"Hostent.Length", unsafe.Offsetof(Hostent{}.Length), 2*ptrSize + 2},
		{"Hostent.AddrList", unsafe.Offsetof(Hostent{}.AddrList), align(2*ptrSize+4, ptrSize)},
		{"AddrInfo.AddrLen", unsafe.Offsetof(AddrInfo{}.AddrLen), 16},
		{"AddrInfo.CanOnName", unsafe.Offsetof(AddrInfo{}.CanOnName), 16 + ptrSize},
		{"AddrInfo.Addr", unsafe.Offsetof(AddrInfo{}.Addr), 16 + 2*ptrSize},
		{"AddrInfo.Next", unsafe.Offsetof(AddrInfo{}.Next), 16 + 3*ptrSize},
		{"WSAData", unsafe.Sizeof(WSAData{}), wsaDataSize},
	}
	for _, test := range tests {
		if test.got != test.want {
			t.Errorf("%s = %d, want %d", test.name, test.got, test.want)
		}
	}
}

func align(value, alignment uintptr) uintptr {
	return (value + alignment - 1) &^ (alignment - 1)
}
