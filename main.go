package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"alda/cli"
	"alda/tsptw"
)

var cpuProf = flag.String("cpuprof", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuProf != "" {
		f, err := os.Create(*cpuProf)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	config := cli.LoadConfig()
	err := tsptw.LoadInstance(config)
	if err != nil {
		fmt.Print(err)
	}
}
