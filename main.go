package main

import (
	"fmt"
	"os"

	"alda/cli"
	"alda/tsptw"
)

func main() {
	config := cli.LoadConfig(os.Args)
	err := tsptw.LoadInstance(config)
	if err != nil {
		fmt.Print(err)
	}
}
