//go:build ((freebsd || linux || darwin) && arm64) || (windows && (amd64 || arm64)) || 386

package ffi

const (
	DefaultAbi = 1
)
