package main

import (
	"unsafe"
	_ "unsafe"

	"github.com/goplus/lib/emb/machine"
)

//go:linkname StoreUint32 llgo.atomicStore
func StoreUint32(addr *uint32, val uint32)

//go:linkname sleep sleep
func sleep(n int)

func main() {
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff480A4))), 0x50D83AA1)
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff4808C))), 0)
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff5f048))), 0)
	buttonPin := machine.GPIO34
	buttonPin.Configure(machine.PinConfig{Mode: machine.PinInput})

	for i := 0; i < 10; i++ {
		if buttonPin.Get() {
			println("yes")
		} else {
			println("no")
		}
		sleep(1)
	}
}
