# C networking bindings

This package binds `arpa/inet.h`, `net/if.h`, `netdb.h`, `netinet/in.h`, `sys/socket.h`, and `sys/un.h`, together with their required C types. Each target is generated as one common libc translation unit, then classified into the existing C packages, with unmodified llcppg commit `da227ff4aaf07718d7ec0dea62e23dab03ac57fb`.

| Target | Header source | Build constraint |
| --- | --- | --- |
| macOS arm64 | macOS 14.0 SDK | `darwin && arm64 && !ios` |
| Linux amd64 | Bootlin stable 2017.05, glibc 2.24 | `linux && amd64 && !android` |
| Linux arm64 | Bootlin stable 2017.05, glibc 2.24 | `linux && arm64 && !android` |

Linux support assumes glibc. Existing implementations remain selected for other targets.

Generated declarations use the existing public function names, including `Bind`, `Accept`, `GetHostByName`, `InetAddr`, `InetNtop`, `Getaddrinfo`, `Send`, and `Recv`. Duplicate `Native` declarations have been removed. `Bind` and `Accept` now accept the generic `*SockAddr` rather than restricting their address argument to `*SockaddrIn`. `Send` and `Recv` retain `c.SizeT` (`uintptr`) for their lengths. Existing string entry points retain `*c.Char`, including on Linux arm64. `legacy_native.go` retains the five untyped socket constants; `legacy_darwin_arm64.go` retains the four native byte-order declarations whose macOS headers expose macros. The root package `c.Char` and `c.Printf` remain unchanged.

Address types retain the existing field names for corresponding C fields. Linux address families use a 16-bit field without the BSD length byte, and Linux `AddrInfo` places `Addr` before `CanOnName`. `In6Addr.U6Addr` remains a byte array with native IPv6 alignment. `Iovec` aliases the shared `c.Iovec` definition from the root C package. The compatibility padding constants follow the native `SockaddrStorage`: `PAD1_SIZE`/`PAD2_SIZE` are 6/112 on macOS and 118/0 on Linux.

The directory `llcppg.cfg` contains only `Dir`, used by the existing local directory-dependency caller to locate prepared headers. It is dependency metadata, not an independent generation recipe. Generation uses the unified libc configurations followed by classification. Each target is kept in one `libc_<os>_<arch>.go` file. Generation inputs, scripts, and header caches remain local; `_generate` is not included. The `llcppg.pub` manifest maps types present in the classified outputs.

The kernel-only macOS `pfctlinput` entry point and the Linux implementation helpers `__h_errno_location` and `__cmsg_nxthdr` are excluded. Types needed to express native layouts, including anonymous union storage and its generated field accessors, remain available. Existing approved generation-error and unsupported-layout exclusions remain recorded in the local generation inputs. `Cmsghdr` contains only its fixed C header fields: 16 bytes with 8-byte alignment on Linux and 12 bytes with 4-byte alignment on macOS. The interface enumeration function is `IfNameindex`; its record is `IfNameindexT`. Linux resolver declarations link `libresolv`. This remains a selection of common native APIs rather than the entire libc API.

The tests compare every retained C function's parameter/result ABI against independently compiled native prototypes and take its address to require symbol resolution. They also check native type sizes, alignments and field offsets, constant values, global symbols and layouts, and generated union accessors. A source inventory check rejects missing coverage. Runtime tests exercise binary socket I/O, IPv4/IPv6 numeric conversion and resolution, interface enumeration/freeing, and TCP loopback through the generic `Bind`/`Accept` signatures.

The independent native C link check passes for macOS. Linux excludes `bindresvport6` with explicit approval: glibc 2.24's header declares it, but neither its original libc nor the validation environment's libc exports it. macOS arm64 and Linux amd64/arm64 pass `llgo test -count=1 -v ./c/net` with GC and `nogc` using LLGo v1.0.8 (`c171abfa130108bc6451e01fcf15a8a3474eb0a7`). The macOS runs use LLD 23.1.3 with the compiler on LLVM 22. Tests use the root C bindings from dependent PR #45, which supplies `c.Iovec`. The glibc 2.24 and macOS 14 SDK header baselines do not establish minimum executable runtime versions.
