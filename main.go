package main

import (
	"fmt"
	"os"
	"time"

	"alda/cli"
	"alda/tsptw"
)

func main() {
	start := time.Now()
	config := cli.LoadConfig(os.Args)
	err := tsptw.LoadInstance(config)
	if err != nil {
		fmt.Print(err)
	}
	duration := time.Since(start)
	fmt.Println("Time: ", duration)
}
