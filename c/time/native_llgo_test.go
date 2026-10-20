//go:build llgo && ((darwin && arm64 && !ios) || (linux && (amd64 || arm64) && !android))

// Copyright (c) 2026 The XGo Authors. Licensed under the Apache License 2.0.

package time_test

import (
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	lltime "github.com/goplus/lib/c/time"
)

func TestNativeTimeLayouts(t *testing.T) {
	var stamp lltime.TimeT
	var clock lltime.ClockT
	var clockID lltime.ClockidT
	var calendar lltime.Tm
	var timespec lltime.Timespec
	cases := []struct {
		name                             string
		size, align, wantSize, wantAlign uintptr
	}{
		{"time_t", unsafe.Sizeof(stamp), unsafe.Alignof(stamp), 8, 8},
		{"clock_t", unsafe.Sizeof(clock), unsafe.Alignof(clock), 8, 8},
		{"clockid_t", unsafe.Sizeof(clockID), unsafe.Alignof(clockID), 4, 4},
		{"tm", unsafe.Sizeof(calendar), unsafe.Alignof(calendar), 56, 8},
		{"timespec", unsafe.Sizeof(timespec), unsafe.Alignof(timespec), 16, 8},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.size != tc.wantSize || tc.align != tc.wantAlign {
				t.Fatalf("size/alignment = %d/%d, want %d/%d", tc.size, tc.align, tc.wantSize, tc.wantAlign)
			}
		})
	}
	if unsafe.Offsetof(calendar.Sec) != 0 || unsafe.Offsetof(calendar.Year) != 20 || unsafe.Offsetof(calendar.Gmtoff) != 40 || unsafe.Offsetof(calendar.Zone) != 48 {
		t.Fatalf("tm offsets: sec=%d year=%d gmtoff=%d zone=%d", unsafe.Offsetof(calendar.Sec), unsafe.Offsetof(calendar.Year), unsafe.Offsetof(calendar.Gmtoff), unsafe.Offsetof(calendar.Zone))
	}
	if unsafe.Offsetof(timespec.Sec) != 0 || unsafe.Offsetof(timespec.Nsec) != 8 {
		t.Fatalf("timespec offsets = %d/%d, want 0/8", unsafe.Offsetof(timespec.Sec), unsafe.Offsetof(timespec.Nsec))
	}
}

func TestNativeCalendarConversions(t *testing.T) {
	stamp := lltime.TimeT(946782245) // 2000-01-02 03:04:05 UTC.
	var utc lltime.Tm
	if result := lltime.GmtimeR(&stamp, &utc); result != &utc {
		t.Fatalf("gmtime_r returned %p, want caller buffer %p", result, &utc)
	}
	if utc.Year != 100 || utc.Mon != 0 || utc.Mday != 2 || utc.Hour != 3 || utc.Min != 4 || utc.Sec != 5 || utc.Wday != 0 || utc.Yday != 1 || utc.Isdst != 0 {
		t.Fatalf("gmtime_r fields = %+v", utc)
	}
	if got := utc.Timegm(); got != stamp {
		t.Fatalf("timegm = %d, want %d", got, stamp)
	}
	var local lltime.Tm
	if result := lltime.LocaltimeR(&stamp, &local); result != &local {
		t.Fatalf("localtime_r returned %p, want caller buffer %p", result, &local)
	}
	if got := local.Timelocal(); got != stamp {
		t.Fatalf("timelocal(localtime_r(t)) = %d, want %d", got, stamp)
	}
}

func TestNativeTimeOutputBuffers(t *testing.T) {
	stamp := lltime.TimeT(946782245)
	var utc lltime.Tm
	if lltime.GmtimeR(&stamp, &utc) != &utc {
		t.Fatal("gmtime_r did not return the caller buffer")
	}
	var buffer [64]byte
	for i := range buffer {
		buffer[i] = 0xa5
	}
	pointer := (*nativeTimeChar)(unsafe.Pointer(&buffer[0]))
	if result := utc.AsctimeR(pointer); result != pointer {
		t.Fatalf("asctime_r returned %p, want %p", result, pointer)
	}
	if got := c.GoString((*c.Char)(unsafe.Pointer(pointer))); got != "Sun Jan  2 03:04:05 2000\n" {
		t.Fatalf("asctime_r = %q", got)
	}
	if buffer[26] != 0xa5 {
		t.Fatalf("asctime_r changed the byte after its 26-byte output: %x", buffer[26])
	}
	expected := lltime.Ctime(&stamp)
	if expected == nil {
		t.Fatal("ctime returned nil")
	}
	want := c.GoString(expected)
	for i := range buffer {
		buffer[i] = 0xa5
	}
	if result := lltime.CtimeR(&stamp, pointer); result != pointer {
		t.Fatalf("ctime_r returned %p, want %p", result, pointer)
	}
	if got := c.GoString((*c.Char)(unsafe.Pointer(pointer))); got != want {
		t.Fatalf("ctime_r = %q, want %q", got, want)
	}
	if buffer[26] != 0xa5 {
		t.Fatalf("ctime_r changed the byte after its 26-byte output: %x", buffer[26])
	}
}

func TestNativeClockQueries(t *testing.T) {
	var before, after, resolution lltime.Timespec
	if result := lltime.ClockGettime(lltime.CLOCK_MONOTONIC, &before); result != 0 {
		t.Fatalf("clock_gettime = %d", result)
	}
	if result := lltime.ClockGettime(lltime.CLOCK_MONOTONIC, &after); result != 0 {
		t.Fatalf("clock_gettime = %d", result)
	}
	if before.Sec < 0 || before.Nsec < 0 || before.Nsec >= 1_000_000_000 || after.Sec < before.Sec || (after.Sec == before.Sec && after.Nsec < before.Nsec) || after.Nsec < 0 || after.Nsec >= 1_000_000_000 {
		t.Fatalf("invalid monotonic clock samples: before=%+v after=%+v", before, after)
	}
	if result := lltime.ClockGetres(lltime.CLOCK_MONOTONIC, &resolution); result != 0 {
		t.Fatalf("clock_getres = %d", result)
	}
	if resolution.Sec < 0 || resolution.Nsec < 0 || resolution.Nsec >= 1_000_000_000 || (resolution.Sec == 0 && resolution.Nsec == 0) {
		t.Fatalf("invalid clock resolution: %+v", resolution)
	}
	var realtime lltime.Timespec
	if result := realtime.Get(lltime.TIME_UTC); result != lltime.TIME_UTC {
		t.Fatalf("timespec_get = %d, want TIME_UTC", result)
	}
	if realtime.Sec <= 0 || realtime.Nsec < 0 || realtime.Nsec >= 1_000_000_000 {
		t.Fatalf("invalid timespec_get output: %+v", realtime)
	}
}
