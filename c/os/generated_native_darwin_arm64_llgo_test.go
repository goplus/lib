//go:build llgo && darwin && arm64 && !ios

package os

import (
	"testing"
	"unsafe"
)

func TestGeneratedMachMessageHeaderLayout(t *testing.T) {
	var header MachMsgHeaderT
	if unsafe.Sizeof(header) != 24 || unsafe.Alignof(header) != 4 {
		t.Fatalf("mach_msg_header_t = %d/%d, want 24/4", unsafe.Sizeof(header), unsafe.Alignof(header))
	}
}
func TestGeneratedUUIDOutputBuffer(t *testing.T) {
	var value UuidT
	for i := range value {
		value[i] = 0xa5
	}
	UuidClear(&value[0])
	for i, b := range value {
		if b != 0 {
			t.Fatalf("uuid_clear left byte %d = %d", i, b)
		}
	}
	if result := UuidIsNull(&value[0]); result == 0 {
		t.Fatal("cleared uuid is not nil")
	}
}
