package main

import (
	"fastgh3/fastbrowser/gh"
	"fmt"
	"time"
)

func main() {
	start := time.Now().UnixMilli()
	for _ = range 10000 {
		_, _ = gh.ReadChart("example")
	}
	total := time.Now().UnixMilli() - start
	fmt.Printf("%d ms\n", total)
}
