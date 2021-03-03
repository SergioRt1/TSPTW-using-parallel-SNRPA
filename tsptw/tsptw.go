package tsptw

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"sync"
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

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	bestRollout := runConcurrent(ctx, config, tsptwInstance)
	fmt.Println("Finish", time.Since(start))
	fmt.Printf("%v violations: %v, Score: %f,  makespan: %f\n", bestRollout.Tour, bestRollout.Violations, -bestRollout.Score, bestRollout.Makespan)

	return nil
}

// Runs nRuns trees of NRPA
func runConcurrent(ctx context.Context, config *cli.Config, t *entities.TSPTW) *nrpa.Rollout {
	chOut := make(chan *nrpa.Rollout, config.NRuns)
	var wg sync.WaitGroup
	best := &nrpa.Rollout{Score: -math.MaxFloat64}

	wg.Add(config.NRuns)
	go func() {
		wg.Wait()
		close(chOut)
	}()
	for i := 0; i < config.NRuns; i++ {
		nrpaInstance := nrpa.NewNRPA(t, config.Levels, config.NIter, config.StabilizationFactor)
		go nrpaInstance.RunConcurrent(ctx, config, t, chOut, &wg)
	}
	for rollout := range chOut {
		if rollout.Score > best.Score {
			best = rollout
		}
	}
	return best
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
