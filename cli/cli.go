package cli

import (
	"fmt"
	"strconv"
	"strings"
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

func LoadConfig(args []string) *Config {
	config := &Config{
		NRuns:               4,
		NIter:               10,
		Levels:              5,
		Timeout:             2,
		FileName:            "rc_201.1.txt",
		StabilizationFactor: 10,
	}
	if len(args) > 1 {
		config.FileName = args[1]
	} else {
		fmt.Println("No test case has been set, using default rc_201.1.txt")
	}
	for i := 2; i < len(args); i++ {
		command := strings.Split(args[i], "=")
		if len(command) == 2 {
			switch command[0] {
			case "time":
				config.Timeout = time.Second * time.Duration(getValue(command[1]))
			case "iter":
				config.NIter = getValue(command[1])
			case "runs":
				config.NRuns = getValue(command[1])
			case "levels":
				config.Levels = getValue(command[1])
			case "sfactor":
				config.StabilizationFactor = getValue(command[1])
			}
		} else {
			fmt.Println("an invalid command has been received (using default): " + args[i])
		}
	}
	return config
}

func getValue(command string) int {
	v, err := strconv.Atoi(command)
	if err != nil {
		fmt.Println("error parsing value of command (using default): " + command)
		return 0
	}
	return v
}
