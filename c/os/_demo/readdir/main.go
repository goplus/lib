package main

import (
	"github.com/goplus/lib/c"
	"github.com/goplus/lib/c/os"
)

func main() {
	fd := os.Open(c.Str("."), os.O_RDONLY)
	if fd == 0 {
		c.Perror(c.Str("open error"))
		os.Exit(1)
	}
	dir := os.Fdopendir(fd)
	if dir == nil {
		c.Perror(c.Str("fdopendir error"))
		os.Exit(1)
	}
	var check int
	for {
		entry := os.Readdir(dir)
		if entry == nil {
			break
		}
		c.Printf(c.Str("%s %d %d\n"), &entry.Name[0], entry.Namlen, entry.Type)
		switch entry.Namlen {
		case 1:
			if entry.Name[0] == '.' {
				check++
			}
		case 2:
			if entry.Name[0] == '.' && entry.Name[1] == '.' {
				check++
			}
		case 7:
			if c.Strncmp(&entry.Name[0], c.Str("main.go"), 7) == 0 {
				check++
			}
		}
	}
	os.Closedir(dir)
	if check != 3 {
		c.Printf(c.Str("readdir error\n"))
		os.Exit(1)
	}
}
