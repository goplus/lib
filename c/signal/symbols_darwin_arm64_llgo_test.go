//go:build llgo && darwin && arm64 && !ios

package signal

import (
	"github.com/goplus/lib/c"
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
	// Expected values are compiled from original macOS 14 SDK headers.
	var signalFunctionSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{
		{"Raise", c.Func(Raise)},
		{"BsdSignal", c.Func(BsdSignal)},
		{"Kill", c.Func(Kill)},
		{"Killpg", c.Func(Killpg)},
		{"Sigaction", c.Func(Sigaction)},
		{"Sigaddset", c.Func(Sigaddset)},
		{"*StackT.Sigaltstack", c.Func((*StackT).Sigaltstack)},
		{"Sigdelset", c.Func(Sigdelset)},
		{"Sigemptyset", c.Func(Sigemptyset)},
		{"Sigfillset", c.Func(Sigfillset)},
		{"Sighold", c.Func(Sighold)},
		{"Sigignore", c.Func(Sigignore)},
		{"Siginterrupt", c.Func(Siginterrupt)},
		{"Sigismember", c.Func(Sigismember)},
		{"Sigpause", c.Func(Sigpause)},
		{"Sigpending", c.Func(Sigpending)},
		{"Sigprocmask", c.Func(Sigprocmask)},
		{"Sigrelse", c.Func(Sigrelse)},
		{"Sigset", c.Func(Sigset)},
		{"Sigsuspend", c.Func(Sigsuspend)},
		{"Sigwait", c.Func(Sigwait)},
		{"Psignal", c.Func(Psignal)},
		{"Sigblock", c.Func(Sigblock)},
		{"Sigsetmask", c.Func(Sigsetmask)},
		{"Sigvec", c.Func(Sigvec)},
		{"NativeSignal", c.Func(NativeSignal)},
	}
	var signalFunctionPrototypes = []struct {
		name      string
		signature reflect.Type
		inputs    []signalABISlot
		output    *signalABISlot
		variadic  bool
	}{
		{"Raise", reflect.TypeOf(Raise), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"BsdSignal", reflect.TypeOf(BsdSignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, &signalABISlot{8, 8, 5}, false},
		{"Kill", reflect.TypeOf(Kill), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Killpg", reflect.TypeOf(Killpg), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigaction", reflect.TypeOf(Sigaction), []signalABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigaddset", reflect.TypeOf(Sigaddset), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"*StackT.Sigaltstack", reflect.TypeOf((*StackT).Sigaltstack), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigdelset", reflect.TypeOf(Sigdelset), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigemptyset", reflect.TypeOf(Sigemptyset), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigfillset", reflect.TypeOf(Sigfillset), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sighold", reflect.TypeOf(Sighold), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigignore", reflect.TypeOf(Sigignore), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Siginterrupt", reflect.TypeOf(Siginterrupt), []signalABISlot{{4, 4, 1}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigismember", reflect.TypeOf(Sigismember), []signalABISlot{{8, 8, 5}, {4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigpause", reflect.TypeOf(Sigpause), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigpending", reflect.TypeOf(Sigpending), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigprocmask", reflect.TypeOf(Sigprocmask), []signalABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigrelse", reflect.TypeOf(Sigrelse), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigset", reflect.TypeOf(Sigset), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, &signalABISlot{8, 8, 5}, false},
		{"Sigsuspend", reflect.TypeOf(Sigsuspend), []signalABISlot{{8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Sigwait", reflect.TypeOf(Sigwait), []signalABISlot{{8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"Psignal", reflect.TypeOf(Psignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, nil, false},
		{"Sigblock", reflect.TypeOf(Sigblock), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigsetmask", reflect.TypeOf(Sigsetmask), []signalABISlot{{4, 4, 1}}, &signalABISlot{4, 4, 1}, false},
		{"Sigvec", reflect.TypeOf(Sigvec), []signalABISlot{{4, 4, 1}, {8, 8, 5}, {8, 8, 5}}, &signalABISlot{4, 4, 1}, false},
		{"NativeSignal", reflect.TypeOf(NativeSignal), []signalABISlot{{4, 4, 1}, {8, 8, 5}}, &signalABISlot{8, 8, 5}, false},
	}
	var signalStaticFunctionPrototypes = []struct {
		name  string
		slots []struct {
			label     string
			goType    reflect.Type
			got, want signalABISlot
		}
	}{}
	var signalTypeLayouts = []struct {
		symbol, property string
		got, want        uintptr
	}{
		{"X__DarwinSigaltstack", "size", unsafe.Sizeof(*new(X__DarwinSigaltstack)), 24},
		{"X__DarwinSigaltstack", "alignment", unsafe.Alignof(*new(X__DarwinSigaltstack)), 8},
		{"X__DarwinSigaltstack", "SsSp offset", unsafe.Offsetof((*new(X__DarwinSigaltstack)).SsSp), 0},
		{"X__DarwinSigaltstack", "SsSp size", unsafe.Sizeof((*new(X__DarwinSigaltstack)).SsSp), 8},
		{"X__DarwinSigaltstack", "SsSp alignment", unsafe.Alignof((*new(X__DarwinSigaltstack)).SsSp), 8},
		{"X__DarwinSigaltstack", "SsSize offset", unsafe.Offsetof((*new(X__DarwinSigaltstack)).SsSize), 8},
		{"X__DarwinSigaltstack", "SsSize size", unsafe.Sizeof((*new(X__DarwinSigaltstack)).SsSize), 8},
		{"X__DarwinSigaltstack", "SsSize alignment", unsafe.Alignof((*new(X__DarwinSigaltstack)).SsSize), 8},
		{"X__DarwinSigaltstack", "SsFlags offset", unsafe.Offsetof((*new(X__DarwinSigaltstack)).SsFlags), 16},
		{"X__DarwinSigaltstack", "SsFlags size", unsafe.Sizeof((*new(X__DarwinSigaltstack)).SsFlags), 4},
		{"X__DarwinSigaltstack", "SsFlags alignment", unsafe.Alignof((*new(X__DarwinSigaltstack)).SsFlags), 4},
		{"StackT", "size", unsafe.Sizeof(*new(StackT)), 24},
		{"StackT", "alignment", unsafe.Alignof(*new(StackT)), 8},
		{"Sigval", "size", unsafe.Sizeof(*new(Sigval)), 8},
		{"Sigval", "alignment", unsafe.Alignof(*new(Sigval)), 8},
		{"X__Siginfo", "size", unsafe.Sizeof(*new(X__Siginfo)), 104},
		{"X__Siginfo", "alignment", unsafe.Alignof(*new(X__Siginfo)), 8},
		{"X__Siginfo", "SiSigno offset", unsafe.Offsetof((*new(X__Siginfo)).SiSigno), 0},
		{"X__Siginfo", "SiSigno size", unsafe.Sizeof((*new(X__Siginfo)).SiSigno), 4},
		{"X__Siginfo", "SiSigno alignment", unsafe.Alignof((*new(X__Siginfo)).SiSigno), 4},
		{"X__Siginfo", "SiErrno offset", unsafe.Offsetof((*new(X__Siginfo)).SiErrno), 4},
		{"X__Siginfo", "SiErrno size", unsafe.Sizeof((*new(X__Siginfo)).SiErrno), 4},
		{"X__Siginfo", "SiErrno alignment", unsafe.Alignof((*new(X__Siginfo)).SiErrno), 4},
		{"X__Siginfo", "SiCode offset", unsafe.Offsetof((*new(X__Siginfo)).SiCode), 8},
		{"X__Siginfo", "SiCode size", unsafe.Sizeof((*new(X__Siginfo)).SiCode), 4},
		{"X__Siginfo", "SiCode alignment", unsafe.Alignof((*new(X__Siginfo)).SiCode), 4},
		{"X__Siginfo", "SiPid offset", unsafe.Offsetof((*new(X__Siginfo)).SiPid), 12},
		{"X__Siginfo", "SiPid size", unsafe.Sizeof((*new(X__Siginfo)).SiPid), 4},
		{"X__Siginfo", "SiPid alignment", unsafe.Alignof((*new(X__Siginfo)).SiPid), 4},
		{"X__Siginfo", "SiUid offset", unsafe.Offsetof((*new(X__Siginfo)).SiUid), 16},
		{"X__Siginfo", "SiUid size", unsafe.Sizeof((*new(X__Siginfo)).SiUid), 4},
		{"X__Siginfo", "SiUid alignment", unsafe.Alignof((*new(X__Siginfo)).SiUid), 4},
		{"X__Siginfo", "SiStatus offset", unsafe.Offsetof((*new(X__Siginfo)).SiStatus), 20},
		{"X__Siginfo", "SiStatus size", unsafe.Sizeof((*new(X__Siginfo)).SiStatus), 4},
		{"X__Siginfo", "SiStatus alignment", unsafe.Alignof((*new(X__Siginfo)).SiStatus), 4},
		{"X__Siginfo", "SiAddr offset", unsafe.Offsetof((*new(X__Siginfo)).SiAddr), 24},
		{"X__Siginfo", "SiAddr size", unsafe.Sizeof((*new(X__Siginfo)).SiAddr), 8},
		{"X__Siginfo", "SiAddr alignment", unsafe.Alignof((*new(X__Siginfo)).SiAddr), 8},
		{"X__Siginfo", "SiValue offset", unsafe.Offsetof((*new(X__Siginfo)).SiValue), 32},
		{"X__Siginfo", "SiValue size", unsafe.Sizeof((*new(X__Siginfo)).SiValue), 8},
		{"X__Siginfo", "SiValue alignment", unsafe.Alignof((*new(X__Siginfo)).SiValue), 8},
		{"X__Siginfo", "SiBand offset", unsafe.Offsetof((*new(X__Siginfo)).SiBand), 40},
		{"X__Siginfo", "SiBand size", unsafe.Sizeof((*new(X__Siginfo)).SiBand), 8},
		{"X__Siginfo", "SiBand alignment", unsafe.Alignof((*new(X__Siginfo)).SiBand), 8},
		{"X__Siginfo", "X__Pad offset", unsafe.Offsetof((*new(X__Siginfo)).X__Pad), 48},
		{"X__Siginfo", "X__Pad size", unsafe.Sizeof((*new(X__Siginfo)).X__Pad), 56},
		{"X__Siginfo", "X__Pad alignment", unsafe.Alignof((*new(X__Siginfo)).X__Pad), 8},
		{"X__SigactionU", "size", unsafe.Sizeof(*new(X__SigactionU)), 8},
		{"X__SigactionU", "alignment", unsafe.Alignof(*new(X__SigactionU)), 8},
		{"SigactionT", "size", unsafe.Sizeof(*new(SigactionT)), 16},
		{"SigactionT", "alignment", unsafe.Alignof(*new(SigactionT)), 8},
		{"SigactionT", "X__SigactionU offset", unsafe.Offsetof((*new(SigactionT)).X__SigactionU), 0},
		{"SigactionT", "X__SigactionU size", unsafe.Sizeof((*new(SigactionT)).X__SigactionU), 8},
		{"SigactionT", "X__SigactionU alignment", unsafe.Alignof((*new(SigactionT)).X__SigactionU), 8},
		{"SigactionT", "SaMask offset", unsafe.Offsetof((*new(SigactionT)).SaMask), 8},
		{"SigactionT", "SaMask size", unsafe.Sizeof((*new(SigactionT)).SaMask), 4},
		{"SigactionT", "SaMask alignment", unsafe.Alignof((*new(SigactionT)).SaMask), 4},
		{"SigactionT", "SaFlags offset", unsafe.Offsetof((*new(SigactionT)).SaFlags), 12},
		{"SigactionT", "SaFlags size", unsafe.Sizeof((*new(SigactionT)).SaFlags), 4},
		{"SigactionT", "SaFlags alignment", unsafe.Alignof((*new(SigactionT)).SaFlags), 4},
		{"SigvecT", "size", unsafe.Sizeof(*new(SigvecT)), 16},
		{"SigvecT", "alignment", unsafe.Alignof(*new(SigvecT)), 8},
		{"SigvecT", "SvHandler offset", unsafe.Offsetof((*new(SigvecT)).SvHandler), 0},
		{"SigvecT", "SvHandler size", unsafe.Sizeof((*new(SigvecT)).SvHandler), 8},
		{"SigvecT", "SvHandler alignment", unsafe.Alignof((*new(SigvecT)).SvHandler), 8},
		{"SigvecT", "SvMask offset", unsafe.Offsetof((*new(SigvecT)).SvMask), 8},
		{"SigvecT", "SvMask size", unsafe.Sizeof((*new(SigvecT)).SvMask), 4},
		{"SigvecT", "SvMask alignment", unsafe.Alignof((*new(SigvecT)).SvMask), 4},
		{"SigvecT", "SvFlags offset", unsafe.Offsetof((*new(SigvecT)).SvFlags), 12},
		{"SigvecT", "SvFlags size", unsafe.Sizeof((*new(SigvecT)).SvFlags), 4},
		{"SigvecT", "SvFlags alignment", unsafe.Alignof((*new(SigvecT)).SvFlags), 4},
		{"SiginfoT", "size", unsafe.Sizeof(*new(SiginfoT)), 104},
		{"SiginfoT", "alignment", unsafe.Alignof(*new(SiginfoT)), 8},
		{"SigsetT", "size", unsafe.Sizeof(*new(SigsetT)), 4},
		{"SigsetT", "alignment", unsafe.Alignof(*new(SigsetT)), 4},
		{"SignalHandler", "size", unsafe.Sizeof(*new(SignalHandler)), 8},
		{"SignalHandler", "alignment", unsafe.Alignof(*new(SignalHandler)), 8},
		{"SysSigname", "global size", unsafe.Sizeof(SysSigname), 256},
		{"SysSigname", "global alignment", unsafe.Alignof(SysSigname), 8},
		{"SysSiglist", "global size", unsafe.Sizeof(SysSiglist), 256},
		{"SysSiglist", "global alignment", unsafe.Alignof(SysSiglist), 8},
	}
	var signalIntegerConstants = []struct {
		name      string
		got, want int64
	}{
		{"NSIG", int64(NSIG), 32},
		{"SIGHUP", int64(SIGHUP), 1},
		{"SIGINT", int64(SIGINT), 2},
		{"SIGQUIT", int64(SIGQUIT), 3},
		{"SIGILL", int64(SIGILL), 4},
		{"SIGTRAP", int64(SIGTRAP), 5},
		{"SIGABRT", int64(SIGABRT), 6},
		{"SIGIOT", int64(SIGIOT), 6},
		{"SIGEMT", int64(SIGEMT), 7},
		{"SIGFPE", int64(SIGFPE), 8},
		{"SIGKILL", int64(SIGKILL), 9},
		{"SIGBUS", int64(SIGBUS), 10},
		{"SIGSEGV", int64(SIGSEGV), 11},
		{"SIGSYS", int64(SIGSYS), 12},
		{"SIGPIPE", int64(SIGPIPE), 13},
		{"SIGALRM", int64(SIGALRM), 14},
		{"SIGTERM", int64(SIGTERM), 15},
		{"SIGURG", int64(SIGURG), 16},
		{"SIGSTOP", int64(SIGSTOP), 17},
		{"SIGTSTP", int64(SIGTSTP), 18},
		{"SIGCONT", int64(SIGCONT), 19},
		{"SIGCHLD", int64(SIGCHLD), 20},
		{"SIGTTIN", int64(SIGTTIN), 21},
		{"SIGTTOU", int64(SIGTTOU), 22},
		{"SIGIO", int64(SIGIO), 23},
		{"SIGXCPU", int64(SIGXCPU), 24},
		{"SIGXFSZ", int64(SIGXFSZ), 25},
		{"SIGVTALRM", int64(SIGVTALRM), 26},
		{"SIGPROF", int64(SIGPROF), 27},
		{"SIGWINCH", int64(SIGWINCH), 28},
		{"SIGINFO", int64(SIGINFO), 29},
		{"SIGUSR1", int64(SIGUSR1), 30},
		{"SIGUSR2", int64(SIGUSR2), 31},
		{"SIGEV_NONE", int64(SIGEV_NONE), 0},
		{"SIGEV_SIGNAL", int64(SIGEV_SIGNAL), 1},
		{"SIGEV_THREAD", int64(SIGEV_THREAD), 3},
		{"ILL_NOOP", int64(ILL_NOOP), 0},
		{"ILL_ILLOPC", int64(ILL_ILLOPC), 1},
		{"ILL_ILLTRP", int64(ILL_ILLTRP), 2},
		{"ILL_PRVOPC", int64(ILL_PRVOPC), 3},
		{"ILL_ILLOPN", int64(ILL_ILLOPN), 4},
		{"ILL_ILLADR", int64(ILL_ILLADR), 5},
		{"ILL_PRVREG", int64(ILL_PRVREG), 6},
		{"ILL_COPROC", int64(ILL_COPROC), 7},
		{"ILL_BADSTK", int64(ILL_BADSTK), 8},
		{"FPE_NOOP", int64(FPE_NOOP), 0},
		{"FPE_FLTDIV", int64(FPE_FLTDIV), 1},
		{"FPE_FLTOVF", int64(FPE_FLTOVF), 2},
		{"FPE_FLTUND", int64(FPE_FLTUND), 3},
		{"FPE_FLTRES", int64(FPE_FLTRES), 4},
		{"FPE_FLTINV", int64(FPE_FLTINV), 5},
		{"FPE_FLTSUB", int64(FPE_FLTSUB), 6},
		{"FPE_INTDIV", int64(FPE_INTDIV), 7},
		{"FPE_INTOVF", int64(FPE_INTOVF), 8},
		{"SEGV_NOOP", int64(SEGV_NOOP), 0},
		{"SEGV_MAPERR", int64(SEGV_MAPERR), 1},
		{"SEGV_ACCERR", int64(SEGV_ACCERR), 2},
		{"BUS_NOOP", int64(BUS_NOOP), 0},
		{"BUS_ADRALN", int64(BUS_ADRALN), 1},
		{"BUS_ADRERR", int64(BUS_ADRERR), 2},
		{"BUS_OBJERR", int64(BUS_OBJERR), 3},
		{"TRAP_BRKPT", int64(TRAP_BRKPT), 1},
		{"TRAP_TRACE", int64(TRAP_TRACE), 2},
		{"CLD_NOOP", int64(CLD_NOOP), 0},
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
		{"SA_ONSTACK", int64(SA_ONSTACK), 1},
		{"SA_RESTART", int64(SA_RESTART), 2},
		{"SA_RESETHAND", int64(SA_RESETHAND), 4},
		{"SA_NOCLDSTOP", int64(SA_NOCLDSTOP), 8},
		{"SA_NODEFER", int64(SA_NODEFER), 16},
		{"SA_NOCLDWAIT", int64(SA_NOCLDWAIT), 32},
		{"SA_SIGINFO", int64(SA_SIGINFO), 64},
		{"SA_USERTRAMP", int64(SA_USERTRAMP), 256},
		{"SA_64REGSET", int64(SA_64REGSET), 512},
		{"SA_USERSPACE_MASK", int64(SA_USERSPACE_MASK), 127},
		{"SIG_BLOCK", int64(SIG_BLOCK), 1},
		{"SIG_UNBLOCK", int64(SIG_UNBLOCK), 2},
		{"SIG_SETMASK", int64(SIG_SETMASK), 3},
		{"SI_USER", int64(SI_USER), 65537},
		{"SI_QUEUE", int64(SI_QUEUE), 65538},
		{"SI_TIMER", int64(SI_TIMER), 65539},
		{"SI_ASYNCIO", int64(SI_ASYNCIO), 65540},
		{"SI_MESGQ", int64(SI_MESGQ), 65541},
		{"SS_ONSTACK", int64(SS_ONSTACK), 1},
		{"SS_DISABLE", int64(SS_DISABLE), 4},
		{"MINSIGSTKSZ", int64(MINSIGSTKSZ), 32768},
		{"SIGSTKSZ", int64(SIGSTKSZ), 131072},
		{"SV_ONSTACK", int64(SV_ONSTACK), 1},
		{"SV_INTERRUPT", int64(SV_INTERRUPT), 2},
		{"SV_RESETHAND", int64(SV_RESETHAND), 4},
		{"SV_NODEFER", int64(SV_NODEFER), 16},
		{"SV_NOCLDSTOP", int64(SV_NOCLDSTOP), 8},
		{"SV_SIGINFO", int64(SV_SIGINFO), 64},
	}
	var signalVariableSymbols = []struct {
		name    string
		address unsafe.Pointer
	}{
		{"SysSigname", unsafe.Pointer(&SysSigname)},
		{"SysSiglist", unsafe.Pointer(&SysSiglist)},
	}
	var signalAccessorChecks = []struct {
		name                string
		check               func() (uintptr, uintptr, uintptr)
		offset, size, align uintptr
	}{
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
		{"*X__SigactionU.XGof_ref___sa_handler", func() (uintptr, uintptr, uintptr) {
			var value X__SigactionU
			view := value.XGof_ref___sa_handler()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
		{"*X__SigactionU.XGof_ref___sa_sigaction", func() (uintptr, uintptr, uintptr) {
			var value X__SigactionU
			view := value.XGof_ref___sa_sigaction()
			return uintptr(unsafe.Pointer(view)) - uintptr(unsafe.Pointer(&value)), unsafe.Sizeof(*view), unsafe.Alignof(*view)
		}, 0, 8, 8},
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
