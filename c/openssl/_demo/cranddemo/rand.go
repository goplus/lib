package main

import (
	"fmt"

	"github.com/goplus/lib/c/openssl"
)

func main() {
	b := make([]byte, 10)

	openssl.RANDBytes(b)
	fmt.Printf("%x\n", b)

	openssl.RANDPrivBytes(b)
	fmt.Printf("%x\n", b)
}
