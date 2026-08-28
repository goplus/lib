//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
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

import "github.com/goplus/lib/c"

const LLGoPackage = "link: -lws2_32"

type (
	SocketT  = uintptr
	SocklenT = c.Int
)

const (
	InvalidSocket SocketT = ^SocketT(0)
	SocketError           = -1
)

func MakeWord(low, high byte) uint16 {
	return uint16(low) | uint16(high)<<8
}

const (
	AF_UNSPEC    = 0
	AF_UNIX      = 1
	AF_LOCAL     = AF_UNIX
	AF_INET      = 2
	AF_IPX       = 6
	AF_APPLETALK = 16
	AF_NETBIOS   = 17
	AF_INET6     = 23
	AF_IRDA      = 26
	AF_BTH       = 32
	AF_MAX       = 34
)

const (
	SOCK_STREAM    = 1
	SOCK_DGRAM     = 2
	SOCK_RAW       = 3
	SOCK_RDM       = 4
	SOCK_SEQPACKET = 5
)

const (
	EAI_AGAIN    = 11002
	EAI_BADFLAGS = 10022
	EAI_FAIL     = 11003
	EAI_FAMILY   = 10047
	EAI_MEMORY   = 8
	EAI_NONAME   = 11001
	EAI_SERVICE  = 10109
	EAI_SOCKTYPE = 10044
)

const INET_ADDRSTRLEN = 16

// SockaddrIn matches the Winsock SOCKADDR_IN layout. Unlike BSD sockaddr_in,
// it has no leading length byte and its address family is 16 bits.
type SockaddrIn struct {
	Family uint16
	Port   uint16
	Addr   InAddr
	Zero   [8]c.Char
}

type SockaddrIn6 struct {
	Family   uint16
	Port     uint16
	Flowinfo c.Uint
	Addr     In6Addr
	ScopeId  c.Uint
}

type SockaddrStorage struct {
	Family uint16
	pad1   [6]c.Char
	align  c.LongLong
	pad2   [112]c.Char
}

type InAddr struct {
	Addr c.Uint
}

type In6Addr struct {
	U6Addr [16]uint8
}

type SockAddr struct {
	Family uint16
	Data   [14]c.Char
}

type Hostent struct {
	Name     *c.Char
	Aliases  **c.Char
	AddrType int16
	Length   int16
	AddrList **c.Char
}
