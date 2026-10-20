//go:build llgo && linux && (amd64 || arm64) && !android

package time_test

import (
	"testing"
	"unsafe"

	lltime "github.com/goplus/lib/c/time"
)

func TestNativeLinuxTimerLayouts(t *testing.T) {
	var timer lltime.TimerT
	var setting lltime.Itimerspec
	if unsafe.Sizeof(timer) != 8 || unsafe.Alignof(timer) != 8 {
		t.Fatalf("timer_t size/alignment = %d/%d, want 8/8", unsafe.Sizeof(timer), unsafe.Alignof(timer))
	}
	if unsafe.Sizeof(setting) != 32 || unsafe.Alignof(setting) != 8 || unsafe.Offsetof(setting.ItInterval) != 0 || unsafe.Offsetof(setting.ItValue) != 16 {
		t.Fatalf("itimerspec layout: size=%d align=%d interval=%d value=%d", unsafe.Sizeof(setting), unsafe.Alignof(setting), unsafe.Offsetof(setting.ItInterval), unsafe.Offsetof(setting.ItValue))
	}
}

func TestNativeLinuxProcessClock(t *testing.T) {
	var clock lltime.ClockidT
	if result := lltime.ClockGetcpuclockid(0, &clock); result != 0 {
		t.Fatalf("clock_getcpuclockid = %d", result)
	}
	var value lltime.Timespec
	if result := lltime.ClockGettime(clock, &value); result != 0 {
		t.Fatalf("process clock_gettime = %d", result)
	}
	if value.Sec < 0 || value.Nsec < 0 || value.Nsec >= 1_000_000_000 {
		t.Fatalf("invalid process CPU clock value: %+v", value)
	}
}
