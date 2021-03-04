package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"alda/cli"
	"alda/tsptw"
)

var (
	cpuProf = flag.String("cpuprof", "", "write cpu profile to file")
	all     = flag.Bool("all", false, "read all case files")
)

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
	fileNames := getFileNames(config)

	for _, config.FileName = range fileNames {
		fmt.Println("Case:", config.FileName)
		err := tsptw.LoadInstance(config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getFileNames(config *cli.Config) []string {
	if *all {
		if files, err := ioutil.ReadDir("./cases"); err == nil {
			names := make([]string, 0, len(files))
			for i := range files {
				names = append(names, files[i].Name())
			}
			return names
		}
	}
	return []string{config.FileName}
}
