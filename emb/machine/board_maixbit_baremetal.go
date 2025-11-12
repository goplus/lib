//go:build k210 && maixbit

package machine

import "github.com/goplus/lib/emb/device/kendryte"

// SPI on the MAix Bit.
var (
	SPI0 = &SPI{
		Bus: kendryte.SPI0,
	}
	SPI1 = &SPI{
		Bus: kendryte.SPI1,
	}
)
