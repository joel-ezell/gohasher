package statistics

import (
	"encoding/json"
	"log"
	"math"
	"sync"
)

type statistics struct {
	// Stores the number of requested indices
	count int
	// Stores the number of items averaged
	NumAveraged int `json:"total"`
	// Stores the average rounded to the nearest integer
	Average int64 `json:"average"`
	// Stores the actual average as a float for best accuracy
	floatAverage float64
	mu           sync.RWMutex
}

var (
	s    *statistics
	once sync.Once
)

func getInstance() *statistics {
	if s == nil {
		once.Do(func() {
			s = &statistics{
				count:       0,
				NumAveraged: 0,
				Average:     0}
		})
	}
	return s
}

// NextIndex This thread-safe function increments the current count and returns the resulting value
func NextIndex() int {
	stats := getInstance()
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.count++
	return stats.count
}

// UpdateAverage This thread-safe function calculates and stores a new average, incorporating the provided value.
// The updated average is returned.
func UpdateAverage(newDuration int64) int64 {
	stats := getInstance()
	stats.mu.Lock()
	defer stats.mu.Unlock()
	newTotal := (stats.floatAverage*float64(stats.NumAveraged) + float64(newDuration))
	stats.NumAveraged++
	stats.floatAverage = newTotal / float64(stats.NumAveraged)
	stats.Average = int64(math.Round(stats.floatAverage))
	return stats.Average
}

// GetStats returns a JSON encoded string representing the statistics
func GetStats() string {
	stats := getInstance()
	stats.mu.RLock()
	defer stats.mu.RUnlock()
	b, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Received error when marshalling JSON: %s", err)
	}
	return string(b)
}
