package passwords

import (
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/joel-ezell/gohasher/statistics"
)

type passwords struct {
	passMap map[int]string
	mu      sync.RWMutex
}

const delaySecs = 5

var instance *passwords
var hashWg sync.WaitGroup

var once sync.Once

func getInstance() *passwords {
	once.Do(func() {
		instance = &passwords{
			passMap: make(map[int]string)}
	})
	return instance
}

func HashAndStore(pwd string) (int, error) {
	start := time.Now()
	hashWg.Add(1)
	index := statistics.NextIndex()
	go hashWorker(index, pwd, start)
	return index, nil
}

func hashWorker(index int, pwd string, start time.Time) {
	time.Sleep(delaySecs * time.Second)
	//TODO: add sha-512 hash
	encodedPwd := base64.StdEncoding.EncodeToString([]byte(pwd))
	p := getInstance()
	p.mu.Lock()
	defer p.mu.Unlock()
	p.passMap[index] = encodedPwd
	duration := time.Since(start)
	statistics.UpdateAverage(duration.Nanoseconds() / 1000)
	hashWg.Done()
}

func GetHashedPassword(index int) (string, error) {
	p := getInstance()
	hashedPwd := p.passMap[index]
	var err error
	if hashedPwd == "" {
		err = errors.New("Requested index not found")
	}

	return hashedPwd, err
}

func WaitToComplete() {
	hashWg.Wait()
}
