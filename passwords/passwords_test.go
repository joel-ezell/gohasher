package passwords

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/joel-ezell/gohasher/statistics"
)

func TestPasswords(t *testing.T) {
	for i := 0; i < 10; i++ {
		pwd := fmt.Sprintf("Password%d", i)
		index, _ := HashAndStore(pwd)
		stats := statistics.GetStats()
		t.Logf("i = %d, index = %d, stats = %s", i, index, stats)

		hashedPwd, _ := GetHashedPassword(index)
		if hashedPwd != "" {
			t.Errorf("hashedPwd is not empty! It's %s", hashedPwd)
		}
	}

	start := time.Now()
	WaitToComplete()
	duration := time.Since(start)
	t.Logf("Waited %d milliseconds to complete", duration.Nanoseconds()/1000000)

	for i := 1; i < 11; i++ {
		pwd := fmt.Sprintf("Password%d", i-1)
		sha := sha256.New()
		sha.Write([]byte(pwd))
		h := base64.StdEncoding.EncodeToString(sha.Sum(nil))

		hashedPwd, _ := GetHashedPassword(i)
		if hashedPwd != h {
			t.Logf("Retrieved hash %s doesn't match expected value %s", hashedPwd, h)
		}
	}

	t.Logf("Final stats: %s", statistics.GetStats())

}
