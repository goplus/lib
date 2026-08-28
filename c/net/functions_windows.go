//go:build windows

// Copyright (c) 2026 The GoPlus Authors. Licensed under the Apache License 2.0.

package net

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

// Winsock uses stdcall on Windows/386. On Windows/AMD64 and Windows/ARM64,
// stdcall maps to the platform's unified native C calling convention.

//go:linkname WSAStartup stdcall.WSAStartup
func WSAStartup(version uint16, data *WSAData) c.Int

//go:linkname WSACleanup stdcall.WSACleanup
func WSACleanup() c.Int

//go:linkname Getaddrinfo stdcall.getaddrinfo
func Getaddrinfo(host *c.Char, port *c.Char, addrInfo *AddrInfo, result **AddrInfo) c.Int

//go:linkname Freeaddrinfo stdcall.freeaddrinfo
func Freeaddrinfo(addrInfo *AddrInfo)

//go:linkname Socket stdcall.socket
func Socket(domain c.Int, typ c.Int, protocol c.Int) SocketT

//go:linkname Bind stdcall.bind
func Bind(sockfd SocketT, addr *SockAddr, addrlen SocklenT) c.Int

//go:linkname Connect stdcall.connect
func Connect(sockfd SocketT, addr *SockAddr, addrlen SocklenT) c.Int

//go:linkname Listen stdcall.listen
func Listen(sockfd SocketT, backlog c.Int) c.Int

//go:linkname Accept stdcall.accept
func Accept(sockfd SocketT, addr *SockAddr, addrlen *SocklenT) SocketT

//go:linkname Closesocket stdcall.closesocket
func Closesocket(sockfd SocketT) c.Int

//go:linkname GetHostByName stdcall.gethostbyname
func GetHostByName(name *c.Char) *Hostent

//go:linkname InetNtop stdcall.inet_ntop
func InetNtop(af c.Int, src c.Pointer, dst *c.Char, size uintptr) *c.Char

//go:linkname InetAddr stdcall.inet_addr
func InetAddr(value *c.Char) c.Uint

//go:linkname Send stdcall.send
func Send(sockfd SocketT, buffer c.Pointer, length c.Int, flags c.Int) c.Int

//go:linkname Recv stdcall.recv
func Recv(sockfd SocketT, buffer c.Pointer, length c.Int, flags c.Int) c.Int

//go:linkname SetSockOpt stdcall.setsockopt
func SetSockOpt(sockfd SocketT, level c.Int, optionName c.Int, optionValue c.Pointer, optionLength SocklenT) c.Int

//go:linkname Ntohs stdcall.ntohs
func Ntohs(value uint16) uint16

//go:linkname Htons stdcall.htons
func Htons(value uint16) uint16

//go:linkname Ntohl stdcall.ntohl
func Ntohl(value c.Uint) c.Uint

//go:linkname Htonl stdcall.htonl
func Htonl(value c.Uint) c.Uint

//go:linkname WSAGetLastError stdcall.WSAGetLastError
func WSAGetLastError() c.Int
