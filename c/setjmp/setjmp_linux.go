//go:build linux

package setjmp

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

//go:linkname Sigsetjmp C.__sigsetjmp
func Sigsetjmp(env *SigjmpBuf, savemask c.Int) c.Int
