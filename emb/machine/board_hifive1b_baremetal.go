//go:build fe310 && hifive1b

package machine

import "github.com/goplus/lib/emb/device/sifive"

// SPI on the HiFive1.
var (
	SPI1 = &SPI{
		Bus: sifive.QSPI1,
	}
)
