package tsptw

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"alda/cli"
	"alda/entities"
	"alda/nrpa"
	"alda/utils"
)

func LoadInstance(config *cli.Config) error {
	b, err := utils.GetFileAsBytes("cases/" + config.FileName)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)
	instance, err := readData(reader)
	if err != nil {
		return err
	}
	bestRollout := &nrpa.Rollout{
		Score: -math.MaxFloat64,
	}
	nrpaInstance := nrpa.NewNRPA(instance, config.Levels, config.NIter)
	policy := nrpaInstance.PreAllocate()

	for i := 0; i < config.NRuns; i++ {
		rollout := nrpaInstance.StableNRPA(config.Levels, nrpaInstance.DataPerLevel[config.Levels], policy, config.StabilizationFactor)
		if rollout.Score > bestRollout.Score {
			bestRollout = rollout
		}
	}
	fmt.Printf("%v violations: %v, Score: %v,  makespan: %v\n", bestRollout.Tour, bestRollout.Violations, bestRollout.Score, bestRollout.Makespan)

	return nil
}

func readData(reader io.Reader) (*entities.TSPTW, error) {
	var n int
	var err error

	if _, err = fmt.Fscanf(reader, "%v\n", &n); err != nil {
		return nil, utils.InvalidInput
	}
	instance := entities.TSPTW{
		N:           n,
		WindowStart: make([]float64, 0, n),
		WindowEnd:   make([]float64, 0, n),
		Distances:   make([][]float64, n),
	}
	for i := 0; i < n; i++ {
		vec := make([]float64, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscanf(reader, "%v", &vec[j]); err != nil {
				return nil, utils.InvalidInput
			}
		}
		_, _ = fmt.Fscanf(reader, "\n")
		instance.Distances[i] = vec
	}
	var s, e float64
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscanf(reader, "%v %v\n", &s, &e); err != nil {
			return nil, utils.InvalidInput
		}
		instance.WindowStart = append(instance.WindowStart, s)
		instance.WindowEnd = append(instance.WindowEnd, e)
	}
	return &instance, nil
}
