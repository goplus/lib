# C networking bindings

This package binds `arpa/inet.h`, `net/if.h`, `netdb.h`, `netinet/in.h`, `sys/socket.h`, and `sys/un.h`, together with the C types those headers require. Each target is generated as one common libc translation unit, then classified into the existing C packages, with unmodified llcppg commit `da227ff4aaf07718d7ec0dea62e23dab03ac57fb`.

| Target | Header source | Build constraint |
| --- | --- | --- |
| macOS arm64 | macOS 14.0 SDK | `darwin && arm64 && !ios` |
| Linux amd64 | Bootlin stable 2017.05, glibc 2.24 | `linux && amd64 && !android` |
| Linux arm64 | Bootlin stable 2017.05, glibc 2.24 | `linux && arm64 && !android` |

Linux support assumes glibc. Existing implementations remain selected for other targets.

`Socket`, `Connect`, `Listen`, `SetSockOpt`, and `Freeaddrinfo` use generated declarations under their existing public names. `legacy_native.go` retains the five untyped socket constants and compatibility declarations where the public Go parameter or return types differ from the generated API, including `Bind`, `Accept`, `GetHostByName`, `Send`, and `Recv`. The distinct raw declarations retain their `Native` names. The root package `c.Char` and `c.Printf` remain unchanged. Address types retain the existing field names for corresponding C fields. Linux address families use a 16-bit field without the BSD length byte, and Linux `AddrInfo` places `Addr` before `CanOnName`. `In6Addr.U6Addr` remains a byte array with native IPv6 alignment. `Iovec` aliases the shared `c.Iovec` definition.

The directory llcppg.cfg contains only Dir, used by the existing local directory-dependency caller to locate prepared headers. It is dependency metadata, not an independent generation recipe. Obsolete generation rules and platform configuration variants have been removed. Actual generation uses the unified libc configurations followed by classification. Each target is kept in one libc_<os>_<arch>.go file. Generation inputs, scripts, and ignored header caches remain local; _generate is not included in these changes. The llcppg.pub manifest maps types present in the classified outputs.

Generation retains the common native C functions, including APIs that also have compatibility declarations. Existing generation-error exclusions and explicitly approved unsupported layouts remain recorded in the local generation inputs. Cmsghdr now contains only the fixed C header fields: 16 bytes with 8-byte alignment on Linux and 12 bytes with 4-byte alignment on macOS, with no trailing flexible-array field. The interface enumeration function is retained as IfNameindex and its record is named IfNameindexT. Restored Linux resolver declarations link libresolv. Selected packed, explicitly aligned, bit-field, flexible-array, wide-scalar, and unresolved callback layouts remain filtered. This is an incomplete selection of common native APIs.

The unified raw outputs passed their measured original-header C/Go ABI comparisons on all three targets. The classified packages pass `llgo test` with GC and `nogc` on macOS arm64 and Linux amd64/arm64. The networking tests check native layouts and binary socket I/O; they do not exercise every generated symbol. This cleanup preserves each package's C symbol coverage on all three targets and both GC modes. The glibc 2.24 and macOS 14 SDK header baselines do not establish minimum executable runtime versions.
