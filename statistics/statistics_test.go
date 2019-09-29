package statistics

import (
	"sync"
	"testing"
)

func TestStats(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup, t *testing.T) {
			index := NextIndex()
			average := UpdateAverage(int64(i))
			t.Logf("i = %d, index = %d, average = %d", i, index, average)
			wg.Done()
		}(i, &wg, t)
	}
	wg.Wait()
	s := GetStats()
	t.Logf("Total stats: %s", s)
}
