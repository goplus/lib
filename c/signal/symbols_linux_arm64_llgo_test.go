//go:build llgo && linux && arm64 && !android

package signal

import (
	"github.com/goplus/lib/c"
	lltime "github.com/goplus/lib/c/time"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

func TestNativeSignalDataLayouts(t *testing.T) {
	// Expected values are compiled from original Bootlin glibc 2.24 headers.
	var signalFunctionSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{
		{"NativeSignal", c.Func(NativeSignal)},
		{"Kill", c.Func(Kill)},
		{"Killpg", c.Func(Killpg)},
		{"Raise", c.Func(Raise)},
		{"Ssignal", c.Func(Ssignal)},
		{"Gsignal", c.Func(Gsignal)},
		{"Psignal", c.Func(Psignal)},
		{"*SiginfoT.Psiginfo", c.Func((*SiginfoT).Psiginfo)},
		{"Sigblock", c.Func(Sigblock)},
		{"Sigsetmask", c.Func(Sigsetmask)},
		{"Siggetmask", c.Func(Siggetmask)},
		{"*SigsetT.Sigemptyset", c.Func((*SigsetT).Sigemptyset)},
		{"*SigsetT.Sigfillset", c.Func((*SigsetT).Sigfillset)},
		{"*SigsetT.Sigaddset", c.Func((*SigsetT).Sigaddset)},
		{"*SigsetT.Sigdelset", c.Func((*SigsetT).Sigdelset)},
		{"*SigsetT.Sigismember", c.Func((*SigsetT).Sigismember)},
		{"Sigprocmask", c.Func(Sigprocmask)},
		{"*SigsetT.Sigsuspend", c.Func((*SigsetT).Sigsuspend)},
		{"Sigaction", c.Func(Sigaction)},
		{"*SigsetT.Sigpending", c.Func((*SigsetT).Sigpending)},
		{"*SigsetT.Sigwait", c.Func((*SigsetT).Sigwait)},
		{"*SigsetT.Sigwaitinfo", c.Func((*SigsetT).Sigwaitinfo)},
		{"*SigsetT.Sigtimedwait", c.Func((*SigsetT).Sigtimedwait)},
		{"Sigqueue", c.Func(Sigqueue)},
		{"Siginterrupt", c.Func(Siginterrupt)},
		{"Sigstack", c.Func(Sigstack)},
		{"*Sigaltstack.Sigaltstack", c.Func((*Sigaltstack).Sigaltstack)},
		{"TimerCreate", c.Func(TimerCreate)},
	}
	var signalFunctionPrototypes = []struct {
		name      string
		signature reflect.Type
		inputs    []signalABISlot
		output    *signalABISlot
		variadic  bool
	}{
		{"NativeSignal", reflect.TypeOf(NativeSignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, &signalABISlot{8, 8, 5}, false},
		{"Kill", reflect.TypeOf(Kill), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Killpg", reflect.TypeOf(Killpg), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Raise", reflect.TypeOf(Raise), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Ssignal", reflect.TypeOf(Ssignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, &signalABISlot{8, 8, 5}, false},
		{"Gsignal", reflect.TypeOf(Gsignal), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Psignal", reflect.TypeOf(Psignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, nil, false},
		{"*SiginfoT.Psiginfo", reflect.TypeOf((*SiginfoT).Psiginfo), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, nil, false},
		{"Sigblock", reflect.TypeOf(Sigblock), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigsetmask", reflect.TypeOf(Sigsetmask), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Siggetmask", reflect.TypeOf(Siggetmask), []signalABISlot{}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigemptyset", reflect.TypeOf((*SigsetT).Sigemptyset), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigfillset", reflect.TypeOf((*SigsetT).Sigfillset), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigaddset", reflect.TypeOf((*SigsetT).Sigaddset), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigdelset", reflect.TypeOf((*SigsetT).Sigdelset), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigismember", reflect.TypeOf((*SigsetT).Sigismember), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigprocmask", reflect.TypeOf(Sigprocmask), []signalABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigsuspend", reflect.TypeOf((*SigsetT).Sigsuspend), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigaction", reflect.TypeOf(Sigaction), []signalABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigpending", reflect.TypeOf((*SigsetT).Sigpending), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigwait", reflect.TypeOf((*SigsetT).Sigwait), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigwaitinfo", reflect.TypeOf((*SigsetT).Sigwaitinfo), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*SigsetT.Sigtimedwait", reflect.TypeOf((*SigsetT).Sigtimedwait), []signalABISlot{{8, 8, 5}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigqueue", reflect.TypeOf(Sigqueue), []signalABISlot{{4, 4, 1}, {4, 4, 1}, {8, 8, 13}}, &signalABISlot{4, 4, 1}, false},
		{"Siginterrupt", reflect.TypeOf(Siginterrupt), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigstack", reflect.TypeOf(Sigstack), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"*Sigaltstack.Sigaltstack", reflect.TypeOf((*Sigaltstack).Sigaltstack), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
	}
	// The static assignment checks arity, result, and exact Go parameter types.
	// Reflecting the full signature would require missing c/time decl-mode method wrappers.
	var _ func(lltime.ClockidT, *Sigevent, *lltime.TimerT) c.Int = TimerCreate
	var signalStaticFunctionPrototypes = []struct {
		name  string
		slots []struct {
			label     string
			goType    reflect.Type
			got, want signalABISlot
		}
	}{
		{"TimerCreate", []struct {
			label     string
			goType    reflect.Type
			got, want signalABISlot
		}{
			{"clock", reflect.TypeOf(lltime.ClockidT(0)), signalABISlot{}, signalABISlot{4, 4, 1}},
			{"event", nil, signalABISlot{unsafe.Sizeof((*Sigevent)(nil)), unsafe.Alignof((*Sigevent)(nil)), 5}, signalABISlot{8, 8, 5}},
			{"timer", nil, signalABISlot{unsafe.Sizeof((*lltime.TimerT)(nil)), unsafe.Alignof((*lltime.TimerT)(nil)), 5}, signalABISlot{8, 8, 5}},
			{"result", reflect.TypeOf(c.Int(0)), signalABISlot{}, signalABISlot{4, 4, 1}},
		}},
	}
	var signalTypeLayouts = []struct {
		symbol, property string
		got, want        uintptr
	}{
		{"SigactionT", "size", unsafe.Sizeof(*new(SigactionT)), 152},
		{"SigactionT", "alignment", unsafe.Alignof(*new(SigactionT)), 8},
		{"SigactionT", "X__SigactionHandler offset", unsafe.Offsetof((*new(SigactionT)).X__SigactionHandler), 0},
		{"SigactionT", "X__SigactionHandler size", unsafe.Sizeof((*new(SigactionT)).X__SigactionHandler), 8},
		{"SigactionT", "X__SigactionHandler alignment", unsafe.Alignof((*new(SigactionT)).X__SigactionHandler), 8},
		{"SigactionT", "SaMask offset", unsafe.Offsetof((*new(SigactionT)).SaMask), 8},
		{"SigactionT", "SaMask size", unsafe.Sizeof((*new(SigactionT)).SaMask), 128},
		{"SigactionT", "SaMask alignment", unsafe.Alignof((*new(SigactionT)).SaMask), 8},
		{"SigactionT", "SaFlags offset", unsafe.Offsetof((*new(SigactionT)).SaFlags), 136},
		{"SigactionT", "SaFlags size", unsafe.Sizeof((*new(SigactionT)).SaFlags), 4},
		{"SigactionT", "SaFlags alignment", unsafe.Alignof((*new(SigactionT)).SaFlags), 4},
		{"SigactionT", "SaRestorer offset", unsafe.Offsetof((*new(SigactionT)).SaRestorer), 144},
		{"SigactionT", "SaRestorer size", unsafe.Sizeof((*new(SigactionT)).SaRestorer), 8},
		{"SigactionT", "SaRestorer alignment", unsafe.Alignof((*new(SigactionT)).SaRestorer), 8},
		{"Sigval", "size", unsafe.Sizeof(*new(Sigval)), 8},
		{"Sigval", "alignment", unsafe.Alignof(*new(Sigval)), 8},
		{"SiginfoT", "size", unsafe.Sizeof(*new(SiginfoT)), 128},
		{"SiginfoT", "alignment", unsafe.Alignof(*new(SiginfoT)), 8},
		{"SiginfoT", "SiSigno offset", unsafe.Offsetof((*new(SiginfoT)).SiSigno), 0},
		{"SiginfoT", "SiSigno size", unsafe.Sizeof((*new(SiginfoT)).SiSigno), 4},
		{"SiginfoT", "SiSigno alignment", unsafe.Alignof((*new(SiginfoT)).SiSigno), 4},
		{"SiginfoT", "SiErrno offset", unsafe.Offsetof((*new(SiginfoT)).SiErrno), 4},
		{"SiginfoT", "SiErrno size", unsafe.Sizeof((*new(SiginfoT)).SiErrno), 4},
		{"SiginfoT", "SiErrno alignment", unsafe.Alignof((*new(SiginfoT)).SiErrno), 4},
		{"SiginfoT", "SiCode offset", unsafe.Offsetof((*new(SiginfoT)).SiCode), 8},
		{"SiginfoT", "SiCode size", unsafe.Sizeof((*new(SiginfoT)).SiCode), 4},
		{"SiginfoT", "SiCode alignment", unsafe.Alignof((*new(SiginfoT)).SiCode), 4},
		{"SiginfoT", "X_sifields offset", unsafe.Offsetof((*new(SiginfoT)).X_sifields), 16},
		{"SiginfoT", "X_sifields size", unsafe.Sizeof((*new(SiginfoT)).X_sifields), 112},
		{"SiginfoT", "X_sifields alignment", unsafe.Alignof((*new(SiginfoT)).X_sifields), 8},
		{"Sigevent", "size", unsafe.Sizeof(*new(Sigevent)), 64},
		{"Sigevent", "alignment", unsafe.Alignof(*new(Sigevent)), 8},
		{"Sigevent", "SigevValue offset", unsafe.Offsetof((*new(Sigevent)).SigevValue), 0},
		{"Sigevent", "SigevValue size", unsafe.Sizeof((*new(Sigevent)).SigevValue), 8},
		{"Sigevent", "SigevValue alignment", unsafe.Alignof((*new(Sigevent)).SigevValue), 8},
		{"Sigevent", "SigevSigno offset", unsafe.Offsetof((*new(Sigevent)).SigevSigno), 8},
		{"Sigevent", "SigevSigno size", unsafe.Sizeof((*new(Sigevent)).SigevSigno), 4},
		{"Sigevent", "SigevSigno alignment", unsafe.Alignof((*new(Sigevent)).SigevSigno), 4},
		{"Sigevent", "SigevNotify offset", unsafe.Offsetof((*new(Sigevent)).SigevNotify), 12},
		{"Sigevent", "SigevNotify size", unsafe.Sizeof((*new(Sigevent)).SigevNotify), 4},
		{"Sigevent", "SigevNotify alignment", unsafe.Alignof((*new(Sigevent)).SigevNotify), 4},
		{"Sigevent", "X_sigevUn offset", unsafe.Offsetof((*new(Sigevent)).X_sigevUn), 16},
		{"Sigevent", "X_sigevUn size", unsafe.Sizeof((*new(Sigevent)).X_sigevUn), 48},
		{"Sigevent", "X_sigevUn alignment", unsafe.Alignof((*new(Sigevent)).X_sigevUn), 8},
		{"SigvalT", "size", unsafe.Sizeof(*new(SigvalT)), 8},
		{"SigvalT", "alignment", unsafe.Alignof(*new(SigvalT)), 8},
		{"SigstackT", "size", unsafe.Sizeof(*new(SigstackT)), 16},
		{"SigstackT", "alignment", unsafe.Alignof(*new(SigstackT)), 8},
		{"SigstackT", "SsSp offset", unsafe.Offsetof((*new(SigstackT)).SsSp), 0},
		{"SigstackT", "SsSp size", unsafe.Sizeof((*new(SigstackT)).SsSp), 8},
		{"SigstackT", "SsSp alignment", unsafe.Alignof((*new(SigstackT)).SsSp), 8},
		{"SigstackT", "SsOnstack offset", unsafe.Offsetof((*new(SigstackT)).SsOnstack), 8},
		{"SigstackT", "SsOnstack size", unsafe.Sizeof((*new(SigstackT)).SsOnstack), 4},
		{"SigstackT", "SsOnstack alignment", unsafe.Alignof((*new(SigstackT)).SsOnstack), 4},
		{"Sigaltstack", "size", unsafe.Sizeof(*new(Sigaltstack)), 24},
		{"Sigaltstack", "alignment", unsafe.Alignof(*new(Sigaltstack)), 8},
		{"Sigaltstack", "SsSp offset", unsafe.Offsetof((*new(Sigaltstack)).SsSp), 0},
		{"Sigaltstack", "SsSp size", unsafe.Sizeof((*new(Sigaltstack)).SsSp), 8},
		{"Sigaltstack", "SsSp alignment", unsafe.Alignof((*new(Sigaltstack)).SsSp), 8},
		{"Sigaltstack", "SsFlags offset", unsafe.Offsetof((*new(Sigaltstack)).SsFlags), 8},
		{"Sigaltstack", "SsFlags size", unsafe.Sizeof((*new(Sigaltstack)).SsFlags), 4},
		{"Sigaltstack", "SsFlags alignment", unsafe.Alignof((*new(Sigaltstack)).SsFlags), 4},
		{"Sigaltstack", "SsSize offset", unsafe.Offsetof((*new(Sigaltstack)).SsSize), 16},
		{"Sigaltstack", "SsSize size", unsafe.Sizeof((*new(Sigaltstack)).SsSize), 8},
		{"Sigaltstack", "SsSize alignment", unsafe.Alignof((*new(Sigaltstack)).SsSize), 8},
		{"X__SighandlerT", "size", unsafe.Sizeof(*new(X__SighandlerT)), 8},
		{"X__SighandlerT", "alignment", unsafe.Alignof(*new(X__SighandlerT)), 8},
		{"SigsetT", "size", unsafe.Sizeof(*new(SigsetT)), 128},
		{"SigsetT", "alignment", unsafe.Alignof(*new(SigsetT)), 8},
		{"X__SigsetT", "size", unsafe.Sizeof(*new(X__SigsetT)), 128},
		{"X__SigsetT", "alignment", unsafe.Alignof(*new(X__SigsetT)), 8},
		{"X__SigsetT", "X__Val offset", unsafe.Offsetof((*new(X__SigsetT)).X__Val), 0},
		{"X__SigsetT", "X__Val size", unsafe.Sizeof((*new(X__SigsetT)).X__Val), 128},
		{"X__SigsetT", "X__Val alignment", unsafe.Alignof((*new(X__SigsetT)).X__Val), 8},
		{"SignalHandler", "size", unsafe.Sizeof(*new(SignalHandler)), 8},
		{"SignalHandler", "alignment", unsafe.Alignof(*new(SignalHandler)), 8},
		{"SigeventT", "size", unsafe.Sizeof(*new(SigeventT)), 64},
		{"SigeventT", "alignment", unsafe.Alignof(*new(SigeventT)), 8},
		{"_llcppg_anon_7", "size", unsafe.Sizeof(*new(_llcppg_anon_7)), 8},
		{"_llcppg_anon_7", "alignment", unsafe.Alignof(*new(_llcppg_anon_7)), 8},
		{"_llcppg_anon_5", "size", unsafe.Sizeof(*new(_llcppg_anon_5)), 112},
		{"_llcppg_anon_5", "alignment", unsafe.Alignof(*new(_llcppg_anon_5)), 8},
		{"_llcppg_anon_6", "size", unsafe.Sizeof(*new(_llcppg_anon_6)), 48},
		{"_llcppg_anon_6", "alignment", unsafe.Alignof(*new(_llcppg_anon_6)), 8},
		{"_llcppg_anon_23", "size", unsafe.Sizeof(*new(_llcppg_anon_23)), 8},
		{"_llcppg_anon_23", "alignment", unsafe.Alignof(*new(_llcppg_anon_23)), 4},
		{"_llcppg_anon_23", "SiPid offset", unsafe.Offsetof((*new(_llcppg_anon_23)).SiPid), 0},
		{"_llcppg_anon_23", "SiPid size", unsafe.Sizeof((*new(_llcppg_anon_23)).SiPid), 4},
		{"_llcppg_anon_23", "SiPid alignment", unsafe.Alignof((*new(_llcppg_anon_23)).SiPid), 4},
		{"_llcppg_anon_23", "SiUid offset", unsafe.Offsetof((*new(_llcppg_anon_23)).SiUid), 4},
		{"_llcppg_anon_23", "SiUid size", unsafe.Sizeof((*new(_llcppg_anon_23)).SiUid), 4},
		{"_llcppg_anon_23", "SiUid alignment", unsafe.Alignof((*new(_llcppg_anon_23)).SiUid), 4},
		{"_llcppg_anon_24", "size", unsafe.Sizeof(*new(_llcppg_anon_24)), 16},
		{"_llcppg_anon_24", "alignment", unsafe.Alignof(*new(_llcppg_anon_24)), 8},
		{"_llcppg_anon_24", "SiTid offset", unsafe.Offsetof((*new(_llcppg_anon_24)).SiTid), 0},
		{"_llcppg_anon_24", "SiTid size", unsafe.Sizeof((*new(_llcppg_anon_24)).SiTid), 4},
		{"_llcppg_anon_24", "SiTid alignment", unsafe.Alignof((*new(_llcppg_anon_24)).SiTid), 4},
		{"_llcppg_anon_24", "SiOverrun offset", unsafe.Offsetof((*new(_llcppg_anon_24)).SiOverrun), 4},
		{"_llcppg_anon_24", "SiOverrun size", unsafe.Sizeof((*new(_llcppg_anon_24)).SiOverrun), 4},
		{"_llcppg_anon_24", "SiOverrun alignment", unsafe.Alignof((*new(_llcppg_anon_24)).SiOverrun), 4},
		{"_llcppg_anon_24", "SiSigval offset", unsafe.Offsetof((*new(_llcppg_anon_24)).SiSigval), 8},
		{"_llcppg_anon_24", "SiSigval size", unsafe.Sizeof((*new(_llcppg_anon_24)).SiSigval), 8},
		{"_llcppg_anon_24", "SiSigval alignment", unsafe.Alignof((*new(_llcppg_anon_24)).SiSigval), 8},
		{"_llcppg_anon_25", "size", unsafe.Sizeof(*new(_llcppg_anon_25)), 16},
		{"_llcppg_anon_25", "alignment", unsafe.Alignof(*new(_llcppg_anon_25)), 8},
		{"_llcppg_anon_25", "SiPid offset", unsafe.Offsetof((*new(_llcppg_anon_25)).SiPid), 0},
		{"_llcppg_anon_25", "SiPid size", unsafe.Sizeof((*new(_llcppg_anon_25)).SiPid), 4},
		{"_llcppg_anon_25", "SiPid alignment", unsafe.Alignof((*new(_llcppg_anon_25)).SiPid), 4},
		{"_llcppg_anon_25", "SiUid offset", unsafe.Offsetof((*new(_llcppg_anon_25)).SiUid), 4},
		{"_llcppg_anon_25", "SiUid size", unsafe.Sizeof((*new(_llcppg_anon_25)).SiUid), 4},
		{"_llcppg_anon_25", "SiUid alignment", unsafe.Alignof((*new(_llcppg_anon_25)).SiUid), 4},
		{"_llcppg_anon_25", "SiSigval offset", unsafe.Offsetof((*new(_llcppg_anon_25)).SiSigval), 8},
		{"_llcppg_anon_25", "SiSigval size", unsafe.Sizeof((*new(_llcppg_anon_25)).SiSigval), 8},
		{"_llcppg_anon_25", "SiSigval alignment", unsafe.Alignof((*new(_llcppg_anon_25)).SiSigval), 8},
		{"_llcppg_anon_26", "size", unsafe.Sizeof(*new(_llcppg_anon_26)), 32},
		{"_llcppg_anon_26", "alignment", unsafe.Alignof(*new(_llcppg_anon_26)), 8},
		{"_llcppg_anon_26", "SiPid offset", unsafe.Offsetof((*new(_llcppg_anon_26)).SiPid), 0},
		{"_llcppg_anon_26", "SiPid size", unsafe.Sizeof((*new(_llcppg_anon_26)).SiPid), 4},
		{"_llcppg_anon_26", "SiPid alignment", unsafe.Alignof((*new(_llcppg_anon_26)).SiPid), 4},
		{"_llcppg_anon_26", "SiUid offset", unsafe.Offsetof((*new(_llcppg_anon_26)).SiUid), 4},
		{"_llcppg_anon_26", "SiUid size", unsafe.Sizeof((*new(_llcppg_anon_26)).SiUid), 4},
		{"_llcppg_anon_26", "SiUid alignment", unsafe.Alignof((*new(_llcppg_anon_26)).SiUid), 4},
		{"_llcppg_anon_26", "SiStatus offset", unsafe.Offsetof((*new(_llcppg_anon_26)).SiStatus), 8},
		{"_llcppg_anon_26", "SiStatus size", unsafe.Sizeof((*new(_llcppg_anon_26)).SiStatus), 4},
		{"_llcppg_anon_26", "SiStatus alignment", unsafe.Alignof((*new(_llcppg_anon_26)).SiStatus), 4},
		{"_llcppg_anon_26", "SiUtime offset", unsafe.Offsetof((*new(_llcppg_anon_26)).SiUtime), 16},
		{"_llcppg_anon_26", "SiUtime size", unsafe.Sizeof((*new(_llcppg_anon_26)).SiUtime), 8},
		{"_llcppg_anon_26", "SiUtime alignment", unsafe.Alignof((*new(_llcppg_anon_26)).SiUtime), 8},
		{"_llcppg_anon_26", "SiStime offset", unsafe.Offsetof((*new(_llcppg_anon_26)).SiStime), 24},
		{"_llcppg_anon_26", "SiStime size", unsafe.Sizeof((*new(_llcppg_anon_26)).SiStime), 8},
		{"_llcppg_anon_26", "SiStime alignment", unsafe.Alignof((*new(_llcppg_anon_26)).SiStime), 8},
		{"_llcppg_anon_27", "size", unsafe.Sizeof(*new(_llcppg_anon_27)), 16},
		{"_llcppg_anon_27", "alignment", unsafe.Alignof(*new(_llcppg_anon_27)), 8},
		{"_llcppg_anon_27", "SiAddr offset", unsafe.Offsetof((*new(_llcppg_anon_27)).SiAddr), 0},
		{"_llcppg_anon_27", "SiAddr size", unsafe.Sizeof((*new(_llcppg_anon_27)).SiAddr), 8},
		{"_llcppg_anon_27", "SiAddr alignment", unsafe.Alignof((*new(_llcppg_anon_27)).SiAddr), 8},
		{"_llcppg_anon_27", "SiAddrLsb offset", unsafe.Offsetof((*new(_llcppg_anon_27)).SiAddrLsb), 8},
		{"_llcppg_anon_27", "SiAddrLsb size", unsafe.Sizeof((*new(_llcppg_anon_27)).SiAddrLsb), 2},
		{"_llcppg_anon_27", "SiAddrLsb alignment", unsafe.Alignof((*new(_llcppg_anon_27)).SiAddrLsb), 2},
		{"_llcppg_anon_28", "size", unsafe.Sizeof(*new(_llcppg_anon_28)), 16},
		{"_llcppg_anon_28", "alignment", unsafe.Alignof(*new(_llcppg_anon_28)), 8},
		{"_llcppg_anon_28", "SiBand offset", unsafe.Offsetof((*new(_llcppg_anon_28)).SiBand), 0},
		{"_llcppg_anon_28", "SiBand size", unsafe.Sizeof((*new(_llcppg_anon_28)).SiBand), 8},
		{"_llcppg_anon_28", "SiBand alignment", unsafe.Alignof((*new(_llcppg_anon_28)).SiBand), 8},
		{"_llcppg_anon_28", "SiFd offset", unsafe.Offsetof((*new(_llcppg_anon_28)).SiFd), 8},
		{"_llcppg_anon_28", "SiFd size", unsafe.Sizeof((*new(_llcppg_anon_28)).SiFd), 4},
		{"_llcppg_anon_28", "SiFd alignment", unsafe.Alignof((*new(_llcppg_anon_28)).SiFd), 4},
		{"_llcppg_anon_29", "size", unsafe.Sizeof(*new(_llcppg_anon_29)), 16},
		{"_llcppg_anon_29", "alignment", unsafe.Alignof(*new(_llcppg_anon_29)), 8},
		{"_llcppg_anon_29", "X_callAddr offset", unsafe.Offsetof((*new(_llcppg_anon_29)).X_callAddr), 0},
		{"_llcppg_anon_29", "X_callAddr size", unsafe.Sizeof((*new(_llcppg_anon_29)).X_callAddr), 8},
		{"_llcppg_anon_29", "X_callAddr alignment", unsafe.Alignof((*new(_llcppg_anon_29)).X_callAddr), 8},
		{"_llcppg_anon_29", "X_syscall offset", unsafe.Offsetof((*new(_llcppg_anon_29)).X_syscall), 8},
		{"_llcppg_anon_29", "X_syscall size", unsafe.Sizeof((*new(_llcppg_anon_29)).X_syscall), 4},
		{"_llcppg_anon_29", "X_syscall alignment", unsafe.Alignof((*new(_llcppg_anon_29)).X_syscall), 4},
		{"_llcppg_anon_29", "X_arch offset", unsafe.Offsetof((*new(_llcppg_anon_29)).X_arch), 12},
		{"_llcppg_anon_29", "X_arch size", unsafe.Sizeof((*new(_llcppg_anon_29)).X_arch), 4},
		{"_llcppg_anon_29", "X_arch alignment", unsafe.Alignof((*new(_llcppg_anon_29)).X_arch), 4},
		{"_llcppg_anon_30", "size", unsafe.Sizeof(*new(_llcppg_anon_30)), 16},
		{"_llcppg_anon_30", "alignment", unsafe.Alignof(*new(_llcppg_anon_30)), 8},
		{"_llcppg_anon_30", "X_function offset", unsafe.Offsetof((*new(_llcppg_anon_30)).X_function), 0},
		{"_llcppg_anon_30", "X_function size", unsafe.Sizeof((*new(_llcppg_anon_30)).X_function), 8},
		{"_llcppg_anon_30", "X_function alignment", unsafe.Alignof((*new(_llcppg_anon_30)).X_function), 8},
		{"_llcppg_anon_30", "X_attribute offset", unsafe.Offsetof((*new(_llcppg_anon_30)).X_attribute), 8},
		{"_llcppg_anon_30", "X_attribute size", unsafe.Sizeof((*new(_llcppg_anon_30)).X_attribute), 8},
		{"_llcppg_anon_30", "X_attribute alignment", unsafe.Alignof((*new(_llcppg_anon_30)).X_attribute), 8},
	}
	var signalIntegerConstants = []struct {
		name      string
		got, want int64
	}{
		{"SA_NOCLDSTOP", int64(SA_NOCLDSTOP), 1},
		{"SA_NOCLDWAIT", int64(SA_NOCLDWAIT), 2},
		{"SA_SIGINFO", int64(SA_SIGINFO), 4},
		{"SA_ONSTACK", int64(SA_ONSTACK), 134217728},
		{"SA_RESTART", int64(SA_RESTART), 268435456},
		{"SA_NODEFER", int64(SA_NODEFER), 1073741824},
		{"SA_RESETHAND", int64(SA_RESETHAND), 2147483648},
		{"SA_INTERRUPT", int64(SA_INTERRUPT), 536870912},
		{"SA_NOMASK", int64(SA_NOMASK), 1073741824},
		{"SA_ONESHOT", int64(SA_ONESHOT), 2147483648},
		{"SA_STACK", int64(SA_STACK), 134217728},
		{"SIG_BLOCK", int64(SIG_BLOCK), 0},
		{"SIG_UNBLOCK", int64(SIG_UNBLOCK), 1},
		{"SIG_SETMASK", int64(SIG_SETMASK), 2},
		{"SI_ASYNCNL", int64(SI_ASYNCNL), -60},
		{"SI_TKILL", int64(SI_TKILL), -6},
		{"SI_SIGIO", int64(SI_SIGIO), -5},
		{"SI_ASYNCIO", int64(SI_ASYNCIO), -4},
		{"SI_MESGQ", int64(SI_MESGQ), -3},
		{"SI_TIMER", int64(SI_TIMER), -2},
		{"SI_QUEUE", int64(SI_QUEUE), -1},
		{"SI_USER", int64(SI_USER), 0},
		{"SI_KERNEL", int64(SI_KERNEL), 128},
		{"ILL_ILLOPC", int64(ILL_ILLOPC), 1},
		{"ILL_ILLOPN", int64(ILL_ILLOPN), 2},
		{"ILL_ILLADR", int64(ILL_ILLADR), 3},
		{"ILL_ILLTRP", int64(ILL_ILLTRP), 4},
		{"ILL_PRVOPC", int64(ILL_PRVOPC), 5},
		{"ILL_PRVREG", int64(ILL_PRVREG), 6},
		{"ILL_COPROC", int64(ILL_COPROC), 7},
		{"ILL_BADSTK", int64(ILL_BADSTK), 8},
		{"FPE_INTDIV", int64(FPE_INTDIV), 1},
		{"FPE_INTOVF", int64(FPE_INTOVF), 2},
		{"FPE_FLTDIV", int64(FPE_FLTDIV), 3},
		{"FPE_FLTOVF", int64(FPE_FLTOVF), 4},
		{"FPE_FLTUND", int64(FPE_FLTUND), 5},
		{"FPE_FLTRES", int64(FPE_FLTRES), 6},
		{"FPE_FLTINV", int64(FPE_FLTINV), 7},
		{"FPE_FLTSUB", int64(FPE_FLTSUB), 8},
		{"SEGV_MAPERR", int64(SEGV_MAPERR), 1},
		{"SEGV_ACCERR", int64(SEGV_ACCERR), 2},
		{"BUS_ADRALN", int64(BUS_ADRALN), 1},
		{"BUS_ADRERR", int64(BUS_ADRERR), 2},
		{"BUS_OBJERR", int64(BUS_OBJERR), 3},
		{"BUS_MCEERR_AR", int64(BUS_MCEERR_AR), 4},
		{"BUS_MCEERR_AO", int64(BUS_MCEERR_AO), 5},
		{"CLD_EXITED", int64(CLD_EXITED), 1},
		{"CLD_KILLED", int64(CLD_KILLED), 2},
		{"CLD_DUMPED", int64(CLD_DUMPED), 3},
		{"CLD_TRAPPED", int64(CLD_TRAPPED), 4},
		{"CLD_STOPPED", int64(CLD_STOPPED), 5},
		{"CLD_CONTINUED", int64(CLD_CONTINUED), 6},
		{"POLL_IN", int64(POLL_IN), 1},
		{"POLL_OUT", int64(POLL_OUT), 2},
		{"POLL_MSG", int64(POLL_MSG), 3},
		{"POLL_ERR", int64(POLL_ERR), 4},
		{"POLL_PRI", int64(POLL_PRI), 5},
		{"POLL_HUP", int64(POLL_HUP), 6},
		{"SIGEV_SIGNAL", int64(SIGEV_SIGNAL), 0},
		{"SIGEV_NONE", int64(SIGEV_NONE), 1},
		{"SIGEV_THREAD", int64(SIGEV_THREAD), 2},
		{"SIGEV_THREAD_ID", int64(SIGEV_THREAD_ID), 4},
		{"SIGHUP", int64(SIGHUP), 1},
		{"SIGINT", int64(SIGINT), 2},
		{"SIGQUIT", int64(SIGQUIT), 3},
		{"SIGILL", int64(SIGILL), 4},
		{"SIGTRAP", int64(SIGTRAP), 5},
		{"SIGABRT", int64(SIGABRT), 6},
		{"SIGIOT", int64(SIGIOT), 6},
		{"SIGBUS", int64(SIGBUS), 7},
		{"SIGFPE", int64(SIGFPE), 8},
		{"SIGKILL", int64(SIGKILL), 9},
		{"SIGUSR1", int64(SIGUSR1), 10},
		{"SIGSEGV", int64(SIGSEGV), 11},
		{"SIGUSR2", int64(SIGUSR2), 12},
		{"SIGPIPE", int64(SIGPIPE), 13},
		{"SIGALRM", int64(SIGALRM), 14},
		{"SIGTERM", int64(SIGTERM), 15},
		{"SIGSTKFLT", int64(SIGSTKFLT), 16},
		{"SIGCHLD", int64(SIGCHLD), 17},
		{"SIGCONT", int64(SIGCONT), 18},
		{"SIGSTOP", int64(SIGSTOP), 19},
		{"SIGTSTP", int64(SIGTSTP), 20},
		{"SIGTTIN", int64(SIGTTIN), 21},
		{"SIGTTOU", int64(SIGTTOU), 22},
		{"SIGURG", int64(SIGURG), 23},
		{"SIGXCPU", int64(SIGXCPU), 24},
		{"SIGXFSZ", int64(SIGXFSZ), 25},
		{"SIGVTALRM", int64(SIGVTALRM), 26},
		{"SIGPROF", int64(SIGPROF), 27},
		{"SIGWINCH", int64(SIGWINCH), 28},
		{"SIGIO", int64(SIGIO), 29},
		{"SIGPWR", int64(SIGPWR), 30},
		{"SIGSYS", int64(SIGSYS), 31},
		{"SIGUNUSED", int64(SIGUNUSED), 31},
		{"NSIG", int64(NSIG), 65},
	}
	var signalVariableSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{}
	var signalAccessorChecks = []struct {
		name                string
		check               func() (uintptr, uintptr, uintptr)
		offset, size, align uintptr
	}{
		{"*_llcppg_anon_7.XGof_ref_sa_handler", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_7
			view := value.XGof_ref_sa_handler()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*_llcppg_anon_7.XGof_ref_sa_sigaction", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_7
			view := value.XGof_ref_sa_sigaction()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*Sigval.XGof_ref_sival_int", func() (uintptr, uintptr, uintptr) {
			var value Sigval
			view := value.XGof_ref_sival_int()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 4, 4},
		{"*Sigval.XGof_ref_sival_ptr", func() (uintptr, uintptr, uintptr) {
			var value Sigval
			view := value.XGof_ref_sival_ptr()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*_llcppg_anon_5.XGof_ref__pad", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__pad()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 112, 4},
		{"*_llcppg_anon_5.XGof_ref__kill", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__kill()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 4},
		{"*_llcppg_anon_5.XGof_ref__timer", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__timer()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
		{"*_llcppg_anon_5.XGof_ref__rt", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__rt()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
		{"*_llcppg_anon_5.XGof_ref__sigchld", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__sigchld()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 32, 8},
		{"*_llcppg_anon_5.XGof_ref__sigfault", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__sigfault()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
		{"*_llcppg_anon_5.XGof_ref__sigpoll", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__sigpoll()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
		{"*_llcppg_anon_5.XGof_ref__sigsys", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_5
			view := value.XGof_ref__sigsys()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
		{"*_llcppg_anon_6.XGof_ref__pad", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_6
			view := value.XGof_ref__pad()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 48, 4},
		{"*_llcppg_anon_6.XGof_ref__tid", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_6
			view := value.XGof_ref__tid()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 4, 4},
		{"*_llcppg_anon_6.XGof_ref__sigev_thread", func() (uintptr, uintptr, uintptr) {
			var value _llcppg_anon_6
			view := value.XGof_ref__sigev_thread()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 16, 8},
	}

	t.Run("function symbols", func(t *testing.T) {
		for _, symbol := range signalFunctionSymbols {
			t.Run(symbol.name, func(t *testing.T) {
				if symbol.address == nil {
					t.Fatal("C function address is nil")
				}
				t.Logf("C function address: %p", symbol.address)
			})
		}
	})
	t.Run("native prototypes", func(t *testing.T) {
		check := func(name string, got reflect.Type, want signalABISlot) {
			if got.Size() != want.size || uintptr(got.Align()) != want.align {
				t.Errorf("%s size/alignment = %d/%d, C = %d/%d", name, got.Size(), got.Align(), want.size, want.align)
			}
			class := 0
			switch got.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				class = 1
			case reflect.Ptr, reflect.UnsafePointer, reflect.Func:
				class = 5
			case reflect.Float32, reflect.Float64:
				class = 8
			case reflect.Struct:
				class = 12
				if got.NumField() == 1 && got.Field(0).Name == "_xgo_union" {
					class = 13
				}
			}
			if class != want.class {
				t.Errorf("%s ABI class = %d (%v), C = %d", name, class, got.Kind(), want.class)
			}
		}
		for _, tc := range signalFunctionPrototypes {
			t.Run(tc.name, func(t *testing.T) {
				inputs := len(tc.inputs)
				if tc.variadic {
					inputs++
				}
				if tc.signature.Kind() != reflect.Func || tc.signature.NumIn() != inputs || tc.signature.IsVariadic() != tc.variadic {
					t.Fatalf("Go signature = %v, C fixed parameters = %d, variadic = %v", tc.signature, len(tc.inputs), tc.variadic)
				}
				for i, slot := range tc.inputs {
					check(tc.name+" parameter", tc.signature.In(i), slot)
				}
				outputs := 0
				if tc.output != nil {
					outputs = 1
				}
				if tc.signature.NumOut() != outputs {
					t.Fatalf("Go result count = %d, C = %d", tc.signature.NumOut(), outputs)
				}
				if tc.output != nil {
					check(tc.name+" result", tc.signature.Out(0), *tc.output)
				}
			})
		}
		for _, tc := range signalStaticFunctionPrototypes {
			t.Run(tc.name, func(t *testing.T) {
				for _, slot := range tc.slots {
					if slot.goType != nil {
						check(tc.name+" "+slot.label, slot.goType, slot.want)
					} else if slot.got != slot.want {
						t.Errorf("%s %s ABI = %+v, C = %+v", tc.name, slot.label, slot.got, slot.want)
					}
				}
			})
		}

	})
	t.Run("native types and globals", func(t *testing.T) {
		for _, tc := range signalTypeLayouts {
			t.Run(tc.symbol+"/"+tc.property, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go layout = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
		for _, symbol := range signalVariableSymbols {
			if symbol.address == nil {
				t.Errorf("%s C variable address is nil", symbol.name)
			}
			t.Logf("%s C variable address: %p", symbol.name, symbol.address)
		}
	})
	t.Run("native constants", func(t *testing.T) {
		for _, tc := range signalIntegerConstants {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.want {
					t.Fatalf("Go constant = %d, C headers = %d", tc.got, tc.want)
				}
			})
		}
	})
	t.Run("union accessors", func(t *testing.T) {
		for _, tc := range signalAccessorChecks {
			t.Run(tc.name, func(t *testing.T) {
				offset, size, align := tc.check()
				if offset != tc.offset || size != tc.size || align != tc.align {
					t.Fatalf("Go field offset/size/alignment = %d/%d/%d, C = %d/%d/%d", offset, size, align, tc.offset, tc.size, tc.align)
				}
			})
		}
	})
	t.Run("symbol inventory", func(t *testing.T) {
		functions, prototypes, types, constants, variables := map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}
		for _, tc := range signalFunctionPrototypes {
			prototypes[tc.name] = true
		}
		for _, tc := range signalStaticFunctionPrototypes {
			prototypes[tc.name] = true
		}

		for _, symbol := range signalFunctionSymbols {
			functions[symbol.name] = true
			if !prototypes[symbol.name] {
				t.Errorf("%s has no native prototype test", symbol.name)
			}
		}
		for _, tc := range signalAccessorChecks {
			functions[tc.name] = true
		}
		for _, tc := range signalTypeLayouts {
			types[tc.symbol] = true
		}
		for _, tc := range signalIntegerConstants {
			constants[tc.name] = true
		}
		for _, symbol := range signalVariableSymbols {
			variables[symbol.name] = true
		}
		constants["LLGoPackage"] = true
		functions["Signal"] = true
		if LLGoPackage != "link" {
			t.Errorf("LLGoPackage = %q, want link", LLGoPackage)
		}
		files := []string{"libc_" + runtime.GOOS + "_" + runtime.GOARCH + ".go"}

		for _, filename := range files {
			source, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			for _, declaration := range source.Decls {
				switch decl := declaration.(type) {
				case *ast.FuncDecl:
					name := decl.Name.Name
					if !ast.IsExported(name) {
						continue
					}
					if decl.Recv != nil {
						receiver, prefix := decl.Recv.List[0].Type, ""
						if pointer, ok := receiver.(*ast.StarExpr); ok {
							receiver, prefix = pointer.X, "*"
						}
						name = prefix + receiver.(*ast.Ident).Name + "." + name
					}
					if !functions[name] {
						t.Errorf("%s: function %s has no test", filename, name)
					}
					if strings.HasPrefix(decl.Name.Name, "Native") && decl.Name.Name != "NativeSignal" {
						t.Errorf("%s retains duplicate %s", filename, name)
					}
					if decl.Doc != nil {
						for _, comment := range decl.Doc.List {
							for _, private := range []string{"C._llcppg_", "C.__sig", "C.__sysv_signal", "C.__libc_current_sigrt"} {
								if strings.Contains(comment.Text, private) {
									t.Errorf("%s still exposes private C entry point %s", filename, private)
								}
							}
						}
					}
				case *ast.GenDecl:
					for _, specification := range decl.Specs {
						switch spec := specification.(type) {
						case *ast.TypeSpec:
							if ast.IsExported(spec.Name.Name) && !types[spec.Name.Name] {
								t.Errorf("%s: type %s has no ABI test", filename, spec.Name.Name)
							}
						case *ast.ValueSpec:
							for _, name := range spec.Names {
								if !ast.IsExported(name.Name) {
									continue
								}
								if decl.Tok == token.CONST && !constants[name.Name] {
									t.Errorf("%s: constant %s has no test", filename, name.Name)
								}
								if decl.Tok == token.VAR && !variables[name.Name] {
									t.Errorf("%s: variable %s has no test", filename, name.Name)
								}
							}
						}
					}
				}
			}
		}
		t.Logf("coverage: %d C functions, %d Go union accessors, %d types/global layouts, %d constants, %d globals", len(signalFunctionSymbols), len(signalAccessorChecks), len(signalTypeLayouts), len(constants), len(signalVariableSymbols))
	})

	var set SigsetT
	var info SiginfoT
	var value Sigval
	if unsafe.Sizeof(set) != nativeSigsetSize || unsafe.Alignof(set) != nativeSigsetAlign {
		t.Fatalf("sigset_t = %d/%d, want %d/%d", unsafe.Sizeof(set), unsafe.Alignof(set), nativeSigsetSize, nativeSigsetAlign)
	}
	if unsafe.Sizeof(info) != nativeSiginfoSize || unsafe.Alignof(info) != 8 {
		t.Fatalf("siginfo_t = %d/%d, want %d/8", unsafe.Sizeof(info), unsafe.Alignof(info), nativeSiginfoSize)
	}
	if unsafe.Offsetof(info.SiSigno) != 0 || unsafe.Offsetof(info.SiErrno) != 4 || unsafe.Offsetof(info.SiCode) != 8 {
		t.Fatalf("siginfo_t header offsets = %d/%d/%d", unsafe.Offsetof(info.SiSigno), unsafe.Offsetof(info.SiErrno), unsafe.Offsetof(info.SiCode))
	}
	if unsafe.Sizeof(value) != 8 || unsafe.Alignof(value) != 8 {
		t.Fatalf("sigval = %d/%d, want 8/8", unsafe.Sizeof(value), unsafe.Alignof(value))
	}
	*value.XGof_ref_sival_int() = -7
	if got := *(*c.Int)(unsafe.Pointer(&value)); got != -7 {
		t.Fatalf("sigval integer storage = %d, want -7", got)
	}
	var marker uintptr
	*value.XGof_ref_sival_ptr() = unsafe.Pointer(&marker)
	if got := *(*unsafe.Pointer)(unsafe.Pointer(&value)); got != unsafe.Pointer(&marker) {
		t.Fatalf("sigval pointer storage = %p, want %p", got, &marker)
	}

}
