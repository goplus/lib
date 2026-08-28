//go:build llgo && windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package net_test

import (
	"testing"

	llnet "github.com/goplus/lib/c/net"
)

func TestWinsockLifecycle(t *testing.T) {
	var data llnet.WSAData
	if result := llnet.WSAStartup(llnet.MakeWord(2, 2), &data); result != 0 {
		t.Fatalf("WSAStartup returned %d", result)
	}
	t.Cleanup(func() {
		if result := llnet.WSACleanup(); result != 0 {
			t.Errorf("WSACleanup returned %d", result)
		}
	})

	if data.Version != llnet.MakeWord(2, 2) {
		t.Fatalf("WSAStartup selected version %#x, want %#x", data.Version, llnet.MakeWord(2, 2))
	}

	socket := llnet.Socket(llnet.AF_INET, llnet.SOCK_STREAM, 0)
	if socket == llnet.InvalidSocket {
		t.Fatalf("socket failed with Winsock error %d", llnet.WSAGetLastError())
	}
	if result := llnet.Closesocket(socket); result != 0 {
		t.Fatalf("closesocket returned %d (Winsock error %d)", result, llnet.WSAGetLastError())
	}
}
