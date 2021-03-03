package cli

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	NRuns               int
	NIter               int
	Levels              int
	Timeout             time.Duration
	FileName            string
	StabilizationFactor int
}

var (
	filename = flag.String("file", "rc_201.1.txt", "file name of a test case")
	timeout  = flag.String("time", "10000", "execution timeout")
	factor   = flag.String("factor", "10", "tree stabilization factor")
	iter     = flag.String("iter", "10", "iterations in the next level, number of children per node")
	levels   = flag.String("levels", "5", "tree levels")
	runs     = flag.String("runs", "4", "number of trees")
)

func LoadConfig() *Config {
	return &Config{
		NRuns:               getValue(*runs),
		NIter:               getValue(*iter),
		Levels:              getValue(*levels),
		Timeout:             time.Second*time.Duration(getValue(*timeout)) - (time.Millisecond * 5),
		FileName:            *filename,
		StabilizationFactor: getValue(*factor),
	}
}

func getValue(command string) int {
	v, err := strconv.Atoi(command)
	if err != nil {
		fmt.Println("error parsing value of command (using default): " + command)
		return 0
	}
	return v
}
