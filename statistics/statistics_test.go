package statistics

import (
	"encoding/json"
	"math"
	"sync"
	"testing"
)

const numItems = 10

type stats struct {
	// Stores the number of items averaged
	NumAveraged int `json:"total"`
	// Stores the average rounded to the nearest integer
	Average int64 `json:"average"`
}

func TestStats(t *testing.T) {
	var wg sync.WaitGroup
	var total = 0

	for i := 0; i < numItems; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup, t *testing.T) {
			index := NextIndex()
			average := UpdateAverage(int64(i))
			total += i
			t.Logf("i = %d, index = %d, average = %d", i, index, average)
			wg.Done()
		}(i, &wg, t)
	}
	wg.Wait()
	s := GetStats()
	t.Logf("Total stats: %s", s)

	var statsJSON stats
	err := json.Unmarshal([]byte(s), &statsJSON)

	if err != nil {
		t.Errorf("Error unmarshalling JSON: %s", err)
		return
	}

	if statsJSON.NumAveraged != numItems {
		t.Errorf("Number of averaged items: %d, expected %d", statsJSON.NumAveraged, numItems)
	}

	var a = int64(math.Round(float64(total) / float64(numItems)))
	if statsJSON.Average != a {
		t.Errorf("Average: %d, expected %d", statsJSON.Average, a)
	}
}
