package cli

import (
	"flag"
	"time"
)

type Config struct {
	NRuns               int
	NIter               int
	Levels              int
	Timeout             time.Duration
	FileName            string
	StabilizationFactor int
	NActors             int
	PActors             int
}

var (
	filename = flag.String("file", "rc_201.1.txt", "file name of a test case")
	timeout  = flag.Int("time", 10000, "execution timeout")
	factor   = flag.Int("factor", 10, "tree stabilization factor, leaves in the last node")
	iter     = flag.Int("iter", 10, "iterations in the next level, number of children per node")
	levels   = flag.Int("levels", 7, "tree levels")
	runs     = flag.Int("runs", 4, "number of trees")
	nActors  = flag.Int("nactors", -1, "number of actors that computes the NRPA tree (default runs)")
	pActors  = flag.Int("pactors", -1, "number of actors that computes the leaves (default factor)")
)

func LoadConfig() *Config {
	nAct := *runs
	pAct := *factor
	if *nActors > 0 {
		nAct = *nActors
	}
	if *pActors > 0 {
		pAct = *pActors
	}
	return &Config{
		NRuns:               *runs,
		NIter:               *iter,
		Levels:              *levels,
		Timeout:             time.Second*time.Duration(*timeout) - (time.Millisecond * 5),
		FileName:            *filename,
		StabilizationFactor: *factor,
		NActors:             nAct,
		PActors:             pAct,
	}
}
