package passwords

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/joel-ezell/gohasher/statistics"
)

type passwords struct {
	// Stores the number of requested indices
	count   int
	passMap map[int]string
	mu      sync.RWMutex
}

const delaySecs = 5

var instance *passwords
var hashWg sync.WaitGroup

var once sync.Once

func getInstance() *passwords {
	// This ensures that the singleton is instantiated only once, even if multiple initial requests arrive at the same time
	once.Do(func() {
		instance = &passwords{
			passMap: make(map[int]string),
			count:   0}
	})
	return instance
}

// HashAndStore Computes a SHA-512 hash of the specified password, encodes it in Base64, then stores the password in a map
func HashAndStore(pwd string) (int, error) {
	// Use a WorkGroup to track all active workers for graceful shutdown
	start := time.Now()
	hashWg.Add(1)
	index := nextIndex()
	go hashWorker(index, pwd, start)
	return index, nil
}

// This thread-safe function increments the current count and returns the resulting value
func nextIndex() int {
	p := getInstance()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.count++
	return p.count
}

func hashWorker(index int, pwd string, start time.Time) {
	time.Sleep(delaySecs * time.Second)
	sha := sha512.New()
	sha.Write([]byte(pwd))
	encodedPwd := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	p := getInstance()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.passMap[index] = encodedPwd
	duration := time.Since(start)
	statistics.UpdateAverage(duration.Nanoseconds() / 1000)
	hashWg.Done()
}

// GetHashedPassword Returns the hashed password at the specified index
func GetHashedPassword(index int) (string, error) {
	p := getInstance()
	p.mu.RLock()
	defer p.mu.RUnlock()
	hashedPwd := p.passMap[index]
	var err error
	if hashedPwd == "" {
		err = errors.New("Requested index not found")
	}

	return hashedPwd, err
}

// WaitToComplete Blocks until all worker Goroutines have completed
func WaitToComplete() {
	hashWg.Wait()
}
