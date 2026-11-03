# Native signal bindings

The macOS arm64 and Linux amd64/arm64 bindings are classified from a unified libc generation using unmodified llcppg commit `da227ff4aaf07718d7ec0dea62e23dab03ac57fb`. The header baselines are macOS 14.0 and Bootlin stable 2017.05 with glibc 2.24. Other targets keep the existing implementation.

| Target | Build constraint | Direct Go dependencies |
| --- | --- | --- |
| macOS arm64 | `darwin && arm64 && !ios` | `c`, `unsafe` |
| Linux amd64 | `linux && amd64 && !android` | `c`, `c/time`, `unsafe` |
| Linux arm64 | `linux && arm64 && !android` | `c`, `c/time`, `unsafe` |

`Signal(sig c.Int, handler SignalHandler) c.Int` retains its existing contract: it installs the handler through `Sigaction` and returns that call's status. Its implementation uses the generated `SigactionT` and union accessors. `NativeSignal` binds C `signal` and returns the previous handler, so it remains a distinct API rather than replacing `Signal`. The generated `Kill` declaration keeps its established public name. No existing public function signature changes in this update.

Native signal records and their field accessors retain their original sizes, alignments and offsets. The previously approved opaque `c.Pointer` for Linux sigevent's optional pthread attributes avoids a package dependency cycle. Necessary implementation types and generated anonymous union storage remain present; their names do not introduce private C runtime entry points.

Private C implementation entry points are excluded: macOS `__sigbits` and its generated C wrapper, and Linux `__sigismember`, `__sigaddset`, `__sigdelset`, `__sysv_signal`, `__libc_current_sigrtmin` and `__libc_current_sigrtmax`. The private macOS wrapper source and its `LLGoFiles` metadata are removed. Linux also excludes `SysSiglist` with explicit approval: the old headers declare `sys_siglist`, but current glibc retains only versioned compatibility exports and cannot resolve an ordinary reference. macOS keeps `SysSiglist` and `SysSigname`. Legacy public C functions retain their native behavior; Linux amd64 `Sigreturn` is a glibc stub and has link/prototype coverage.

Each retained C function has a live symbol-address check and a prototype check against facts independently compiled from the original headers. Types, fields, constants, native globals and every generated union accessor are checked as well. The source inventory test rejects missing public symbol coverage. Go storage records whose sole field is `_xgo_union` are classified as C unions; their original size and alignment checks still apply.

Linux `TimerCreate` uses a compile-time function-signature assignment and per-slot ABI checks. Its link-address and complete timer lifecycle tests remain enabled. This avoids reflecting the full signature, which currently requires missing `c/time` method wrappers while that dependency uses `LLGoPackage="decl"`. The dependency and compiler are unchanged; this test accommodation does not repair that reflection limitation.

Runtime tests exercise empty/fill/add/delete/member signal-set operations, guarded `sigaction` output, handler installation and restoration, real delivery with `Raise`, and `Kill(self, 0)`. Linux additionally exercises sigevent layout and a non-notifying timer's create/query/disarm/delete lifecycle. Handler callbacks only store the received signal number and do not allocate. These tests run with GC and `nogc` on all three targets.

Validation composes the root `c` files from PR #45 at `99ddac92c7335a53a1c996174e1753bcd72c5050` and `c/time` from PR #51 at `855f0c294499b32a0d3380d1d608592255155e89`. Merge those prerequisites first; this PR changes only `c/signal`. Tests use LLGo v1.0.8 (`c171abfa130108bc6451e01fcf15a8a3474eb0a7`), with LLVM 22 and LLD 23.1.3 for native macOS linking. The header baselines do not establish a minimum supported executable runtime version.

`llcppg.cfg` contains only the prepared-header directory for local dependency discovery; it is not a standalone generation recipe. Actual generation uses the unified libc inputs followed by classification. `llcppg.pub` describes the classified public C type mappings. Inputs, scripts and header caches remain local and are not included in this PR.
