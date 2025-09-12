package main

import (
	"unsafe"
)

//go:linkname StoreUint32 llgo.atomicStore
func StoreUint32(addr *uint32, val uint32)

//go:linkname sleep sleep
func sleep(tm int)

func main() {
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff480A4))), 0x50D83AA1)
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff4808C))), 0)
	StoreUint32((*uint32)(unsafe.Pointer(uintptr(0x3ff5f048))), 0)
	buttonPin := machine.GPIO34
	buttonPin.Configure(machine.PinConfig{Mode: machine.PinInput})
	println(buttonPin.Get())
	sleep(2)
	println(buttonPin.Get())
	println(buttonState)
	// time.Sleep(200 * time.Millisecond)
}
