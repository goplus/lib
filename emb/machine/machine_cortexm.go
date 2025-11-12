//go:build cortexm

package machine

import "github.com/goplus/lib/emb/device/arm"

// CPUReset performs a hard system reset.
func CPUReset() {
	arm.SystemReset()
}
