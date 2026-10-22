//go:build llgo && linux && (amd64 || arm64) && !android

package signal

import (
	"os"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	lltime "github.com/goplus/lib/c/time"
)

const (
	nativeActionSize        = 152
	nativeActionFlagsOffset = 136
	nativeSigsetSize        = 128
	nativeSigsetAlign       = 8
	nativeSiginfoSize       = 128
	nativeUserSignal        = 10
	nativeNoNotification    = 1
)

func nativeEmptySet(set *SigsetT) c.Int                { return set.Sigemptyset() }
func nativeFillSet(set *SigsetT) c.Int                 { return set.Sigfillset() }
func nativeAddSet(set *SigsetT, number c.Int) c.Int    { return set.Sigaddset(number) }
func nativeDeleteSet(set *SigsetT, number c.Int) c.Int { return set.Sigdelset(number) }
func nativeSetMember(set *SigsetT, number c.Int) c.Int { return set.Sigismember(number) }

func TestNativeSigeventLayout(t *testing.T) {
	var event Sigevent
	if unsafe.Sizeof(event) != 64 || unsafe.Alignof(event) != 8 {
		t.Fatalf("sigevent = %d/%d, want 64/8", unsafe.Sizeof(event), unsafe.Alignof(event))
	}
	thread := event.X_sigevUn.XGof_ref__sigev_thread()
	offset := unsafe.Offsetof(event.X_sigevUn) + unsafe.Offsetof(thread.X_attribute)
	if unsafe.Sizeof(*thread) != 16 || unsafe.Sizeof(thread.X_function) != 8 || offset != 24 {
		t.Fatalf("sigevent thread layout: size=%d function=%d attribute offset=%d", unsafe.Sizeof(*thread), unsafe.Sizeof(thread.X_function), offset)
	}
	var marker uintptr
	thread.X_attribute = c.Pointer(&marker)
	if got := *(*c.Pointer)(unsafe.Add(unsafe.Pointer(&event), offset)); got != c.Pointer(&marker) {
		t.Fatalf("attribute pointer = %p, want %p", got, &marker)
	}
}

func TestNativeSignalTimerLifecycle(t *testing.T) {
	var event Sigevent
	event.SigevNotify = SIGEV_NONE
	var timer lltime.TimerT
	if result := TimerCreate(lltime.CLOCK_MONOTONIC, &event, &timer); result != 0 {
		t.Fatalf("timer_create = %d", result)
	}
	defer func() {
		if result := timer.TimerDelete(); result != 0 {
			t.Errorf("timer_delete = %d", result)
		}
	}()
	state := lltime.Itimerspec{ItInterval: lltime.Timespec{Sec: 123, Nsec: 456}, ItValue: lltime.Timespec{Sec: 789, Nsec: 123}}
	if result := timer.TimerGettime(&state); result != 0 {
		t.Fatalf("timer_gettime = %d", result)
	}
	if state.ItInterval.Sec != 0 || state.ItInterval.Nsec != 0 || state.ItValue.Sec != 0 || state.ItValue.Nsec != 0 {
		t.Fatalf("new timer is not disarmed: %+v", state)
	}
	disabled := lltime.Itimerspec{}
	old := lltime.Itimerspec{ItValue: lltime.Timespec{Sec: 123, Nsec: 456}}
	if result := timer.TimerSettime(0, &disabled, &old); result != 0 {
		t.Fatalf("timer_settime(disarm) = %d", result)
	}
	if old.ItValue.Sec != 0 || old.ItValue.Nsec != 0 {
		t.Fatalf("old timer value = %+v, want zero", old)
	}
}

func TestNativeSignalActionLayout(t *testing.T) {
	var action SigactionT
	checks := []struct {
		name      string
		got, want uintptr
	}{
		{"size", unsafe.Sizeof(action), nativeActionSize},
		{"alignment", unsafe.Alignof(action), 8},
		{"handler offset", unsafe.Offsetof(action.X__SigactionHandler), 0},
		{"handler width", unsafe.Sizeof(action.X__SigactionHandler), 8},
		{"mask offset", unsafe.Offsetof(action.SaMask), 8},
		{"mask width", unsafe.Sizeof(action.SaMask), nativeSigsetSize},
		{"flags offset", unsafe.Offsetof(action.SaFlags), nativeActionFlagsOffset},
		{"restorer offset", unsafe.Offsetof(action.SaRestorer), 144},
		{"restorer width", unsafe.Sizeof(action.SaRestorer), 8},
	}
	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d, want %d", check.name, check.got, check.want)
		}
	}
	if t.Failed() {
		return // The current smaller layout must not be passed to a C output parameter.
	}
	var query struct {
		before [8]byte
		action SigactionT
		after  [8]byte
	}
	for i := range query.before {
		query.before[i] = 0xa5
		query.after[i] = 0xa5
	}
	if result := Sigaction(nativeUserSignal, nil, &query.action); result != 0 {
		t.Fatalf("sigaction query = %d", result)
	}
	for i := range query.before {
		if query.before[i] != 0xa5 || query.after[i] != 0xa5 {
			t.Fatalf("sigaction query overwrote a boundary guard at %d", i)
		}
	}
}

func TestNativeSignalWrapper(t *testing.T) {
	var original SigactionT
	if result := Sigaction(nativeUserSignal, nil, &original); result != 0 {
		t.Fatalf("read original signal action = %d", result)
	}
	defer func() {
		if result := Sigaction(nativeUserSignal, &original, nil); result != 0 {
			t.Errorf("restore original signal action = %d", result)
		}
	}()
	deliveredSignal = 0
	handler := SignalHandler(func(number c.Int) { deliveredSignal = number })
	if result := Signal(nativeUserSignal, handler); result != 0 {
		t.Fatalf("Signal = %d", result)
	}
	var actual SigactionT
	if result := Sigaction(nativeUserSignal, nil, &actual); result != 0 {
		t.Fatalf("read installed signal action = %d", result)
	}
	got := *(*unsafe.Pointer)(unsafe.Pointer(actual.X__SigactionHandler.XGof_ref_sa_handler()))
	want := *(*unsafe.Pointer)(unsafe.Pointer(&handler))
	if got != want {
		t.Fatalf("handler = %p, want %p", got, want)
	}
	if member := actual.SaMask.Sigismember(nativeUserSignal); member != 0 {
		t.Fatalf("signal mask membership = %d, want 0", member)
	}
	if result := Raise(nativeUserSignal); result != 0 {
		t.Fatalf("raise = %d", result)
	}
	if deliveredSignal != nativeUserSignal {
		t.Fatalf("delivered signal = %d, want %d", deliveredSignal, nativeUserSignal)
	}
	if result := Kill(c.Int(os.Getpid()), 0); result != 0 {
		t.Fatalf("kill(self, 0) = %d", result)
	}

}
