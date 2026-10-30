//go:build !windows && !(darwin && arm64 && !ios) && !(linux && amd64 && !android) && !(linux && arm64 && !android)

package pthread

import "github.com/goplus/lib/c"

// Attr represents POSIX thread attributes.
type Attr struct {
	Detached byte
	SsSp     *c.Char
	SsSize   uintptr
}

type Key c.Uint
