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
		Timeout:             time.Hour,
		FileName:            "rc_201.1.txt",
		StabilizationFactor: 10,
	}
	for i := 1; i < len(args); i++ {
		command := strings.Split(args[i], "=")
		if len(command) == 2 {
			switch command[0] {
			case "case":
				config.FileName = command[1]
			case "time":
				config.Timeout = time.Second*time.Duration(getValue(command[1])) - (time.Millisecond * 5)
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
