package httpPool

import (
	"net/http"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type PClient struct {
	client       *http.Client
	maxPoolSize  int
	cSemaphore   chan int
	reqPerSecond int
	rateLimiter  <-chan time.Time
}

func NewPClient(stdClient *http.Client, maxPoolSize int, reqPerSec int) *PClient {
	var semaphore chan int = nil
	if maxPoolSize > 0 {
		semaphore = make(chan int, maxPoolSize)
	}

	var emitter <-chan time.Time = nil
	if reqPerSec > 0 {
		emitter = time.NewTicker(time.Second / time.Duration(reqPerSec)).C
	}

	return &PClient{
		client:       stdClient,
		maxPoolSize:  maxPoolSize,
		cSemaphore:   semaphore,
		reqPerSecond: reqPerSec,
		rateLimiter:  emitter,
	}
}

func (c *PClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoPool(req)
}

func (c *PClient) DoPool(req *http.Request) (*http.Response, error) {
	if c.maxPoolSize > 0 {
		c.cSemaphore <- 1
		defer func() {
			<-c.cSemaphore
		}()
	}

	if c.reqPerSecond > 0 {
		<-c.rateLimiter
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
