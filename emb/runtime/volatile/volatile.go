// Package volatile provides definitions for volatile loads and stores. These
// are implemented as compiler builtins.
//
// The load operations load a volatile value. The store operations store to a
// volatile value. The compiler will emit exactly one load or store operation
// when possible and will not reorder volatile operations. However, the compiler
// may move other operations across load/store operations, so make sure that all
// relevant loads/stores are done in a volatile way if this is a problem.
//
// These loads and stores are commonly used to read/write values from memory
// mapped peripheral devices. They do not provide atomicity, use the sync/atomic
// package for that.
//
// For more details: https://llvm.org/docs/LangRef.html#volatile-memory-accesses
// and https://blog.regehr.org/archives/28.
package volatile

import _ "unsafe"

const (
	LLGoPackage = "link"
)

// LoadUint8 loads the volatile value *addr.
//
//go:linkname LoadUint8 llgo.atomicLoad
func LoadUint8(addr *uint8) (val uint8)

// LoadUint16 loads the volatile value *addr.
//
//go:linkname LoadUint16 llgo.atomicLoad
func LoadUint16(addr *uint16) (val uint16)

// LoadUint32 loads the volatile value *addr.
//
//go:linkname LoadUint32 llgo.atomicLoad
func LoadUint32(addr *uint32) (val uint32)

// LoadUint64 loads the volatile value *addr.
//
//go:linkname LoadUint64 llgo.atomicLoad
func LoadUint64(addr *uint64) (val uint64)

// StoreUint8 stores val to the volatile value *addr.
//
//go:linkname StoreUint8 llgo.atomicStore
func StoreUint8(addr *uint8, val uint8)

// StoreUint16 stores val to the volatile value *addr.
//
//go:linkname StoreUint16 llgo.atomicStore
func StoreUint16(addr *uint16, val uint16)

// StoreUint32 stores val to the volatile value *addr.
//
//go:linkname StoreUint32 llgo.atomicStore
func StoreUint32(addr *uint32, val uint32)

// StoreUint64 stores val to the volatile value *addr.
//
//go:linkname StoreUint64 llgo.atomicStore
func StoreUint64(addr *uint64, val uint64)
