//go:build !windows && !((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

package time

import (
	"github.com/goplus/lib/c"
	_ "unsafe"
)

//go:linkname Time C.time
func Time(timer *TimeT) TimeT

//go:linkname Mktime C.mktime
func Mktime(timer *Tm) TimeT

//go:linkname Difftime C.difftime
func Difftime(end, start TimeT) float64

//go:linkname Gmtime C.gmtime
func Gmtime(timer *TimeT) *Tm

//go:linkname Localtime C.localtime
func Localtime(timer *TimeT) *Tm

//go:linkname Clock C.clock
func Clock() ClockT

//go:linkname ClockGettime C.clock_gettime
func ClockGettime(clkId ClockidT, tp *Timespec) c.Int

//go:linkname ClockSettime C.clock_settime
func ClockSettime(clkId ClockidT, tp *Timespec) c.Int

//go:linkname ClockGetres C.clock_getres
func ClockGetres(clkId ClockidT, res *Timespec) c.Int
