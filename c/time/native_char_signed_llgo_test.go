//go:build llgo && ((darwin && arm64 && !ios) || (linux && amd64 && !android))

package time_test

import "github.com/goplus/lib/c"

type nativeTimeChar = c.Char
