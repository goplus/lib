//go:build llgo && darwin && arm64 && !ios

package signal

import (
	"os"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
)

const (
	nativeActionSize        = 16
	nativeActionFlagsOffset = 12
	nativeSigsetSize        = 4
	nativeSigsetAlign       = 4
	nativeSiginfoSize       = 104
	nativeUserSignal        = 30
	nativeNoNotification    = 0
)

func nativeEmptySet(set *SigsetT) c.Int                { return Sigemptyset(set) }
func nativeFillSet(set *SigsetT) c.Int                 { return Sigfillset(set) }
func nativeAddSet(set *SigsetT, number c.Int) c.Int    { return Sigaddset(set, number) }
func nativeDeleteSet(set *SigsetT, number c.Int) c.Int { return Sigdelset(set, number) }
func nativeSetMember(set *SigsetT, number c.Int) c.Int { return Sigismember(set, number) }

func TestNativeSignalActionLayout(t *testing.T) {
	var action SigactionT
	checks := []struct {
		name      string
		got, want uintptr
	}{
		{"size", unsafe.Sizeof(action), nativeActionSize},
		{"alignment", unsafe.Alignof(action), 8},
		{"handler offset", unsafe.Offsetof(action.X__SigactionU), 0},
		{"handler width", unsafe.Sizeof(*action.X__SigactionU.XGof_ref___sa_handler()), 8},
		{"mask offset", unsafe.Offsetof(action.SaMask), 8},
		{"mask width", unsafe.Sizeof(action.SaMask), nativeSigsetSize},
		{"flags offset", unsafe.Offsetof(action.SaFlags), nativeActionFlagsOffset},
	}
	for _, check := range checks {
		if check.got != check.want {
			t.Errorf("%s = %d, want %d", check.name, check.got, check.want)
		}
	}
	if t.Failed() {
		return
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
		t.Fatalf("read original action = %d", result)
	}
	defer func() {
		if result := Sigaction(nativeUserSignal, &original, nil); result != 0 {
			t.Errorf("restore original action = %d", result)
		}
	}()
	deliveredSignal = 0
	handler := SignalHandler(func(number c.Int) { deliveredSignal = number })
	if result := Signal(nativeUserSignal, handler); result != 0 {
		t.Fatalf("Signal = %d", result)
	}
	var actual SigactionT
	if result := Sigaction(nativeUserSignal, nil, &actual); result != 0 {
		t.Fatalf("read installed action = %d", result)
	}
	want := *(*unsafe.Pointer)(unsafe.Pointer(&handler))
	got := *(*unsafe.Pointer)(unsafe.Pointer(actual.X__SigactionU.XGof_ref___sa_handler()))
	if got != want {
		t.Fatalf("installed handler = %p, want %p", got, want)
	}
	if actual.SaMask != 0 || actual.SaFlags != 0 {
		t.Fatalf("installed mask/flags = %d/%d, want 0/0", actual.SaMask, actual.SaFlags)
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
