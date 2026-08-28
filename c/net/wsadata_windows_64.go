//go:build windows && !386

package net

import "github.com/goplus/lib/c"

const wsaDataSize = 408

// WSAData matches the 64-bit Winsock WSADATA layout.
type WSAData struct {
	Version      uint16
	HighVersion  uint16
	MaxSockets   uint16
	MaxUdpDg     uint16
	VendorInfo   *c.Char
	Description  [257]c.Char
	SystemStatus [129]c.Char
}
