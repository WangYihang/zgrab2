package main

import (
	"flag"
	"fmt"

	"github.com/WangYihang/zgrab2"
	"github.com/WangYihang/zgrab2/bin"
	_ "github.com/WangYihang/zgrab2/modules"
)

// main wraps the "true" main, bin.ZGrab2Main(), after importing all scan
// modules in ZGrab2.
func main() {
	var ver = flag.Bool("version", false, "Show version")
	flag.Parse()
	if *ver {
		fmt.Println(zgrab2.Commit)
	} else {
		bin.ZGrab2Main()
	}
}
