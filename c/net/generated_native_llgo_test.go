//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package net

import (
	"github.com/goplus/lib/c"
	"runtime"
	"syscall"
	"testing"
	"unsafe"
)

type networkABISlot struct {
	size, align uintptr
	class       int
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
		n := Send(pair[0], unsafe.Pointer(&input[offset]), c.SizeT(len(input)-offset), 0)
		if n <= 0 {
			t.Fatalf("send = %d", n)
		}
		offset += int(n)
	}
	for offset := 0; offset < len(output); {
		n := Recv(pair[1], unsafe.Pointer(&output[offset]), c.SizeT(len(output)-offset), 0)
		if n <= 0 {
			t.Fatalf("recv = %d", n)
		}
		offset += int(n)
	}
	if output != input {
		t.Fatalf("socket data = %v, want %v", output, input)
	}
	t.Run("byte order and IPv4 rendering", func(t *testing.T) {
		if Htons(0x1234) != 0x3412 || Ntohs(0x3412) != 0x1234 || Htonl(0x12345678) != 0x78563412 || Ntohl(0x78563412) != 0x12345678 {
			t.Fatal("native byte order conversion failed")
		}
		addr := InAddr{Addr: InetAddr(c.Str("127.0.0.1"))}
		if got := c.GoString((*c.Char)(unsafe.Pointer(addr.InetNtoa()))); got != "127.0.0.1" {
			t.Fatalf("inet_ntoa = %q", got)
		}
	})
	t.Run("IPv4 and IPv6 numeric addresses", func(t *testing.T) {
		for _, tc := range []struct {
			family c.Int
			text   string
		}{{AF_INET, "127.0.0.1"}, {AF_INET6, "::1"}} {
			t.Run(tc.text, func(t *testing.T) {
				var address [16]byte
				text := append([]byte(tc.text), 0)
				source := (*c.Char)(unsafe.Pointer(&text[0]))
				if result := InetPton(tc.family, (*nativeNetworkChar)(unsafe.Pointer(source)), unsafe.Pointer(&address[0])); result != 1 {
					t.Fatalf("inet_pton = %d", result)
				}
				var rendered [64]c.Char
				if InetNtop(tc.family, unsafe.Pointer(&address[0]), &rendered[0], SocklenT(len(rendered))) == nil {
					t.Fatal("inet_ntop failed")
				}
				if got := c.GoString(&rendered[0]); got != tc.text {
					t.Fatalf("inet_ntop = %q", got)
				}
				hints := AddrInfo{Flags: AI_NUMERICHOST | AI_NUMERICSERV, Family: tc.family, SockType: SOCK_STREAM}
				var info *AddrInfo
				if result := Getaddrinfo(source, c.Str("32123"), &hints, &info); result != 0 {
					t.Fatalf("getaddrinfo = %d", result)
				}
				if info == nil {
					t.Fatal("getaddrinfo returned nil")
				}
				defer Freeaddrinfo(info)
				if info.Family != tc.family || info.Addr == nil {
					t.Fatalf("addrinfo family/address = %d/%p", info.Family, info.Addr)
				}
				var host [64]nativeNetworkChar
				var service [16]nativeNetworkChar
				if result := info.Addr.Getnameinfo(info.AddrLen, &host[0], SocklenT(len(host)), &service[0], SocklenT(len(service)), NI_NUMERICHOST|NI_NUMERICSERV); result != 0 {
					t.Fatalf("getnameinfo = %d", result)
				}
				if got := c.GoString((*c.Char)(unsafe.Pointer(&host[0]))); got != tc.text {
					t.Fatalf("numeric host = %q", got)
				}
				if got := c.GoString((*c.Char)(unsafe.Pointer(&service[0]))); got != "32123" {
					t.Fatalf("numeric service = %q", got)
				}
			})
		}
		if In6addrAny.U6Addr != [16]uint8{} || In6addrLoopback.U6Addr != [16]uint8{15: 1} {
			t.Fatal("native IPv6 globals are incorrect")
		}
	})
	t.Run("interface index lifecycle", func(t *testing.T) {
		loopback := "lo"
		if runtime.GOOS == "darwin" {
			loopback = "lo0"
		}
		loopbackText := append([]byte(loopback), 0)
		index := IfNametoindex((*nativeNetworkChar)(unsafe.Pointer(&loopbackText[0])))
		if index == 0 {
			t.Fatal("loopback interface index is zero")
		}
		var name [256]nativeNetworkChar
		if IfIndextoname(index, &name[0]) == nil {
			t.Fatal("if_indextoname failed")
		}
		if got := c.GoString((*c.Char)(unsafe.Pointer(&name[0]))); got != loopback {
			t.Fatalf("interface name = %q", got)
		}
		interfaces := IfNameindex()
		if interfaces == nil {
			t.Fatal("if_nameindex failed")
		}
		defer interfaces.IfFreenameindex()
		found := false
		for i := 0; ; i++ {
			entry := (*IfNameindexT)(unsafe.Add(unsafe.Pointer(interfaces), uintptr(i)*unsafe.Sizeof(*interfaces)))
			if entry.IfIndex == 0 && entry.IfName == nil {
				break
			}
			if entry.IfIndex == index && c.GoString((*c.Char)(unsafe.Pointer(entry.IfName))) == loopback {
				found = true
			}
		}
		if !found {
			t.Fatal("loopback interface is absent from if_nameindex")
		}
	})
	t.Run("TCP loopback IPv4 compatibility", func(t *testing.T) {
		listener := Socket(AF_INET, SOCK_STREAM, 0)
		if listener < 0 {
			t.Fatal("listener socket failed")
		}
		defer syscall.Close(int(listener))
		address := SockaddrIn{Family: SaFamilyT(AF_INET), Addr: InAddr{Addr: Htonl(0x7f000001)}}
		if runtime.GOOS == "darwin" {
			*(*uint8)(unsafe.Pointer(&address)) = uint8(unsafe.Sizeof(address))
		}
		if result := Bind(listener, &address, SocklenT(unsafe.Sizeof(address))); result != 0 {
			t.Fatalf("bind = %d", result)
		}
		if result := Listen(listener, 1); result != 0 {
			t.Fatalf("listen = %d", result)
		}
		length := SocklenT(unsafe.Sizeof(address))
		if result := Getsockname(listener, (*SockAddr)(unsafe.Pointer(&address)), &length); result != 0 {
			t.Fatalf("getsockname = %d", result)
		}
		if length != SocklenT(unsafe.Sizeof(address)) || address.Port == 0 {
			t.Fatalf("bound address length/port = %d/%d", length, address.Port)
		}
		client := Socket(AF_INET, SOCK_STREAM, 0)
		if client < 0 {
			t.Fatal("client socket failed")
		}
		defer syscall.Close(int(client))
		if result := Connect(client, (*SockAddr)(unsafe.Pointer(&address)), length); result != 0 {
			t.Fatalf("connect = %d", result)
		}
		var peer SockaddrIn
		peerLength := SocklenT(unsafe.Sizeof(peer))
		accepted := Accept(listener, &peer, &peerLength)
		if accepted < 0 {
			t.Fatalf("accept = %d", accepted)
		}
		defer syscall.Close(int(accepted))
		if peer.Family != SaFamilyT(AF_INET) || peerLength != SocklenT(unsafe.Sizeof(address)) {
			t.Fatalf("peer family/length = %d/%d", peer.Family, peerLength)
		}
		input := byte(0xa5)
		if n := Send(client, unsafe.Pointer(&input), 1, 0); n != 1 {
			t.Fatalf("TCP send = %d", n)
		}
		var output byte
		if n := Recv(accepted, unsafe.Pointer(&output), 1, 0); n != 1 || output != input {
			t.Fatalf("TCP recv = %d/%x", n, output)
		}
		if result := Shutdown(client, SHUT_WR); result != 0 {
			t.Fatalf("shutdown = %d", result)
		}
	})

}
