//go:build (darwin && arm64 && !ios) || (linux && amd64 && !android) || (linux && arm64 && !android)

/*
 * Copyright (c) 2024 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package net

import (
	"unsafe"

	"github.com/goplus/lib/c"
)

// Keep the existing untyped constants usable with Socket's c.Int parameter.
const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_RDM       = 4
	SOCK_SEQPACKET = 5
)

const (
	ALIGNSIZE = unsafe.Sizeof(c.LongLong(0))
	MAXSIZE   = 128
	PAD1_SIZE = ALIGNSIZE - unsafe.Sizeof(byte(0)) - unsafe.Sizeof(byte(0))
	PAD2_SIZE = MAXSIZE - unsafe.Sizeof(byte(0)) - unsafe.Sizeof(byte(0)) - PAD1_SIZE - ALIGNSIZE
)

//go:linkname Socket C.socket
func Socket(domain c.Int, typ c.Int, protocol c.Int) c.Int

//go:linkname Bind C.bind
func Bind(sockfd c.Int, addr *SockaddrIn, addrlen c.Uint) c.Int

//go:linkname Connect C.connect
func Connect(sockfd c.Int, addr *SockAddr, addrlen c.Uint) c.Int

//go:linkname Listen C.listen
func Listen(sockfd c.Int, backlog c.Int) c.Int

//go:linkname Accept C.accept
func Accept(sockfd c.Int, addr *SockaddrIn, addrlen *c.Uint) c.Int

//go:linkname GetHostByName C.gethostbyname
func GetHostByName(name *c.Char) *Hostent

//go:linkname InetNtop C.inet_ntop
func InetNtop(af c.Int, src c.Pointer, dst *c.Char, size c.Uint) *c.Char

//go:linkname InetAddr C.inet_addr
func InetAddr(s *c.Char) c.Uint

//go:linkname Send C.send
func Send(c.Int, c.Pointer, uintptr, c.Int) c.Long

//go:linkname Recv C.recv
func Recv(c.Int, c.Pointer, uintptr, c.Int) c.Long

//go:linkname SetSockOpt C.setsockopt
func SetSockOpt(socket c.Int, level c.Int, optionName c.Int, optionValue c.Pointer, sockLen c.Uint) c.Int

//go:linkname Ntohs C.ntohs
func Ntohs(x uint16) uint16

//go:linkname Htons C.htons
func Htons(x uint16) uint16

//go:linkname Ntohl C.ntohl
func Ntohl(x c.Uint) c.Uint

//go:linkname Htonl C.htonl
func Htonl(x c.Uint) c.Uint

//go:linkname Getaddrinfo C.getaddrinfo
func Getaddrinfo(host *c.Char, port *c.Char, addrInfo *AddrInfo, result **AddrInfo) c.Int

//go:linkname Freeaddrinfo C.freeaddrinfo
func Freeaddrinfo(addrInfo *AddrInfo)
