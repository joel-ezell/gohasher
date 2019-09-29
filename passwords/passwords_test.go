package passwords

import (
	"fmt"
	"testing"
	"time"

	"github.com/joel-ezell/gohasher/statistics"
)

func TestStats(t *testing.T) {
	for i := 0; i < 10; i++ {
		pwd := fmt.Sprintf("Password%d", i)
		index, _ := HashAndStore(pwd)
		stats := statistics.GetStats()
		t.Logf("i = %d, index = %d, stats = %s", i, index, stats)

		hashedPwd, _ := GetHashedPassword(index)
		if hashedPwd == "" {
			t.Logf("hashedPwd is empty, as expected")
		} else {
			t.Logf("hashedPwd is not empty! It's %s", hashedPwd)
		}
	}

	start := time.Now()
	WaitToComplete()
	duration := time.Since(start)
	t.Logf("Waited %d milliseconds to complete", duration.Nanoseconds()/1000000)

	for i := 1; i < 11; i++ {
		hashedPwd, _ := GetHashedPassword(i)
		if hashedPwd == "" {
			t.Logf("hashedPwd should not be empty!")
		} else {
			t.Logf("hashedPwd is now not empty, as expected: %s", hashedPwd)
		}
	}

	t.Logf("Final stats: %s", statistics.GetStats())

}
