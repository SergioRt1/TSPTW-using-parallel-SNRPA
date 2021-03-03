package tsptw

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"alda/cli"
	"alda/entities"
	"alda/nrpa"
	"alda/utils"
)

func LoadInstance(config *cli.Config) error {
	start := time.Now()
	tsptwInstance, err := loadData(config)
	if err != nil {
		return err
	}

	nrpaInstance := nrpa.NewNRPA(tsptwInstance, config.Levels, config.NIter, config.StabilizationFactor)
	policy := nrpaInstance.PreAllocate()

	done := make(chan *nrpa.Rollout)
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	go run(ctx, config, nrpaInstance, policy, done)
	var bestRollout *nrpa.Rollout
	select {
	case <-ctx.Done():
		bestRollout = nrpaInstance.FindCurrentBest()
		fmt.Println("Finish with timeout in:", time.Since(start))
	case bestRollout = <-done:
		fmt.Println("Finish successfully in:", time.Since(start))
	}
	fmt.Printf("%v violations: %v, Score: %f,  makespan: %f\n", bestRollout.Tour, bestRollout.Violations, -bestRollout.Score, bestRollout.Makespan)

	return nil
}

func run(ctx context.Context, config *cli.Config, nrpaInstance *nrpa.NRPA, policy [][]float64, done chan *nrpa.Rollout) {
	bestRollout := &nrpa.Rollout{Score: -math.MaxFloat64}
	for i := 0; i < config.NRuns; i++ {
		rollout := nrpaInstance.StableNRPA(ctx, config.Levels, nrpaInstance.DataPerLevel[config.Levels], policy)
		if rollout.Score > bestRollout.Score {
			bestRollout = rollout
		}
	}
	done <- bestRollout
}

func loadData(config *cli.Config) (*entities.TSPTW, error) {
	b, err := utils.GetFileAsBytes("cases/" + config.FileName)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)
	tsptwInstance, err := readData(reader)
	return tsptwInstance, err
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
