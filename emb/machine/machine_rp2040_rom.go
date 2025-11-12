//go:build llgo && rp2040

package machine

import (
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/lib/emb/runtime/interrupt"
)

const (
	LLGoFiles = "_wrap/machine_rp2040_rom.c"
)

//go:linkname reset_usb_boot C.reset_usb_boot
func reset_usb_boot(c.Uint32T, c.Uint32T)

//go:linkname flash_do_cmd C.flash_do_cmd
func flash_do_cmd(*c.Uint8T, *c.Uint8T, c.SizeT)

//go:linkname flash_range_write C.flash_range_write
func flash_range_write(c.Uint32T, *c.Uint8T, c.SizeT)

//go:linkname flash_erase_blocks C.flash_erase_blocks
func flash_erase_blocks(c.Uint32T, c.SizeT)

func enterBootloader() {
	reset_usb_boot(0, 0)
}

func doFlashCommand(tx []byte, rx []byte) error {
	if len(tx) != len(rx) {
		return errFlashInvalidWriteLength
	}

	// C.flash_do_cmd(
	// 	(*C.uint8_t)(unsafe.Pointer(&tx[0])),
	// 	(*C.uint8_t)(unsafe.Pointer(&rx[0])),
	// 	C.ulong(len(tx)))
	flash_do_cmd(
		(*c.Uint8T)(unsafe.Pointer(&tx[0])),
		(*c.Uint8T)(unsafe.Pointer(&rx[0])),
		c.SizeT(len(tx)),
	)

	return nil
}

// Flash related code
// const memoryStart = C.XIP_BASE // memory start for purpose of erase
const XIP_BASE = 0x10000000
const memoryStart = XIP_BASE // memory start for purpose of erase

func (f flashBlockDevice) writeAt(p []byte, off int64) (n int, err error) {
	if writeAddress(off)+uintptr(XIP_BASE) > FlashDataEnd() {
		return 0, errFlashCannotWritePastEOF
	}

	state := interrupt.Disable()
	defer interrupt.Restore(state)

	// rp2040 writes to offset, not actual address
	// e.g. real address 0x10003000 is written to at
	// 0x00003000
	address := writeAddress(off)
	padded := flashPad(p, int(f.WriteBlockSize()))

	// C.flash_range_write(C.uint32_t(address),
	// 	(*C.uint8_t)(unsafe.Pointer(&padded[0])),
	// 	C.ulong(len(padded)))
	flash_range_write(
		c.Uint32T(address),
		(*c.Uint8T)(unsafe.Pointer(&padded[0])),
		c.SizeT(len(padded)),
	)

	return len(padded), nil
}

func (f flashBlockDevice) eraseBlocks(start, length int64) error {
	address := writeAddress(start * f.EraseBlockSize())
	if address+uintptr(XIP_BASE) > FlashDataEnd() {
		return errFlashCannotErasePastEOF
	}

	state := interrupt.Disable()
	defer interrupt.Restore(state)

	// C.flash_erase_blocks(C.uint32_t(address), C.ulong(length*f.EraseBlockSize()))
	flash_erase_blocks(c.Uint32T(address), c.SizeT(length*f.EraseBlockSize()))
	return nil
}
