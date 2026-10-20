//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package net

import (
	"github.com/goplus/lib/c"
	"runtime"
	"syscall"
	"testing"
	"unsafe"
)

func TestGeneratedNetworkLayouts(t *testing.T) {
	familyOffset := uintptr(0)
	addrInfoAddrOffset := uintptr(24)
	canOnNameOffset := uintptr(32)
	controlSize := uintptr(16)
	controlAlign := uintptr(8)
	if runtime.GOOS == "darwin" {
		familyOffset = 1
		addrInfoAddrOffset = 32
		canOnNameOffset = 24
		controlSize = 12
		controlAlign = 4
	}
	var address SockAddr
	var in6 In6Addr
	var storage SockaddrStorage
	var info AddrInfo
	var control Cmsghdr
	checks := []struct {
		name      string
		got, want uintptr
	}{
		{"sockaddr size", unsafe.Sizeof(address), 16}, {"sockaddr family", unsafe.Offsetof(address.Family), familyOffset},
		{"in6_addr size", unsafe.Sizeof(in6), 16}, {"in6_addr alignment", unsafe.Alignof(in6), 4},
		{"sockaddr_storage size", unsafe.Sizeof(storage), 128}, {"sockaddr_storage alignment", unsafe.Alignof(storage), 8},
		{"addrinfo size", unsafe.Sizeof(info), 48}, {"addrinfo addr", unsafe.Offsetof(info.Addr), addrInfoAddrOffset}, {"addrinfo canonical name", unsafe.Offsetof(info.CanOnName), canOnNameOffset},
		{"cmsghdr size", unsafe.Sizeof(control), controlSize}, {"cmsghdr alignment", unsafe.Alignof(control), controlAlign}}
	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d, want %d", check.name, check.got, check.want)
		}
	}
}
func TestGeneratedSocketIO(t *testing.T) {
	var pair [2]c.Int
	if result := Socketpair(c.Int(AF_UNIX), c.Int(SOCK_STREAM), 0, &pair[0]); result != 0 {
		t.Fatalf("socketpair = %d", result)
	}
	defer syscall.Close(int(pair[0]))
	defer syscall.Close(int(pair[1]))
	var kind c.Int
	length := SocklenT(unsafe.Sizeof(kind))
	if result := Getsockopt(pair[0], SOL_SOCKET, SO_TYPE, unsafe.Pointer(&kind), &length); result != 0 {
		t.Fatalf("getsockopt = %d", result)
	}
	if kind != SOCK_STREAM || length != SocklenT(unsafe.Sizeof(kind)) {
		t.Fatalf("socket type/length = %d/%d", kind, length)
	}
	input := [6]byte{0, 1, 127, 128, 254, 255}
	var output [6]byte
	for offset := 0; offset < len(input); {
		n := NativeSend(pair[0], unsafe.Pointer(&input[offset]), SizeT(len(input)-offset), 0)
		if n <= 0 {
			t.Fatalf("send = %d", n)
		}
		offset += int(n)
	}
	for offset := 0; offset < len(output); {
		n := NativeRecv(pair[1], unsafe.Pointer(&output[offset]), SizeT(len(output)-offset), 0)
		if n <= 0 {
			t.Fatalf("recv = %d", n)
		}
		offset += int(n)
	}
	if output != input {
		t.Fatalf("socket data = %v, want %v", output, input)
	}
}
