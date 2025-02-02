package grpcPool

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
	"time"
)

var (
	ErrClosed        = errors.New("grpc pool: client pool is closed")
	ErrTimeout       = errors.New("grpc pool: client pool timed out")
	ErrAlreadyClosed = errors.New("grpc pool: the connection was already closed")
	ErrFullPool      = errors.New("grpc pool: closing a ClientConn into a full pool")
)

type Factory func() (*grpc.ClientConn, error)

type FactoryWithContext func(context.Context) (*grpc.ClientConn, error)

type Pool struct {
	clients         chan ClientConn
	factory         FactoryWithContext
	idleTimeout     time.Duration
	maxLifeDuration time.Duration
	mu              sync.RWMutex
}

type ClientConn struct {
	*grpc.ClientConn
	pool          *Pool
	timeUsed      time.Time
	timeInitiated time.Time
	unhealthy     bool
}

func New(factory Factory, init, capacity int, idleTimeout time.Duration,
	maxLifeDuration ...time.Duration) (*Pool, error) {
	return NewWithContext(context.Background(), func(ctx context.Context) (*grpc.ClientConn, error) { return factory() },
		init, capacity, idleTimeout, maxLifeDuration...)
}

func NewWithContext(ctx context.Context, factory FactoryWithContext, init, capacity int, idleTimeout time.Duration,
	maxLifeDuration ...time.Duration) (*Pool, error) {

	if capacity <= 0 {
		capacity = 1
	}
	if init < 0 {
		init = 0
	}
	if init > capacity {
		init = capacity
	}
	p := &Pool{
		clients:     make(chan ClientConn, capacity),
		factory:     factory,
		idleTimeout: idleTimeout,
	}
	if len(maxLifeDuration) > 0 {
		p.maxLifeDuration = maxLifeDuration[0]
	}
	for i := 0; i < init; i++ {
		c, err := factory(ctx)
		if err != nil {
			return nil, err
		}

		p.clients <- ClientConn{
			ClientConn:    c,
			pool:          p,
			timeUsed:      time.Now(),
			timeInitiated: time.Now(),
		}
	}

	for i := 0; i < capacity-init; i++ {
		p.clients <- ClientConn{
			pool: p,
		}
	}
	return p, nil
}

func (p *Pool) getClients() chan ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.clients
}

func (p *Pool) Close() {
	p.mu.Lock()
	clients := p.clients
	p.clients = nil
	p.mu.Unlock()

	if clients == nil {
		return
	}

	close(clients)
	for client := range clients {
		if client.ClientConn == nil {
			continue
		}
		client.ClientConn.Close()
	}
}

// IsClosed returns true if the client pool is closed.
func (p *Pool) IsClosed() bool {
	return p == nil || p.getClients() == nil
}

func (p *Pool) Get(ctx context.Context) (*ClientConn, error) {
	clients := p.getClients()
	if clients == nil {
		return nil, ErrClosed
	}

	wrapper := ClientConn{
		pool: p,
	}
	select {
	case wrapper = <-clients:
	case <-ctx.Done():
		return nil, ErrTimeout // it would better returns ctx.Err()
	}

	idleTimeout := p.idleTimeout
	if wrapper.ClientConn != nil && idleTimeout > 0 &&
		wrapper.timeUsed.Add(idleTimeout).Before(time.Now()) {

		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	}

	var err error
	if wrapper.ClientConn == nil {
		wrapper.ClientConn, err = p.factory(ctx)
		if err != nil {
			clients <- ClientConn{
				pool: p,
			}
		}
		wrapper.timeInitiated = time.Now()
	}

	return &wrapper, err
}

func (c *ClientConn) Unhealthy() {
	c.unhealthy = true
}

func (c *ClientConn) Close() error {
	if c == nil {
		return nil
	}
	if c.ClientConn == nil {
		return ErrAlreadyClosed
	}
	if c.pool.IsClosed() {
		return ErrClosed
	}

	maxDuration := c.pool.maxLifeDuration
	if maxDuration > 0 && c.timeInitiated.Add(maxDuration).Before(time.Now()) {
		c.Unhealthy()
	}

	wrapper := ClientConn{
		pool:       c.pool,
		ClientConn: c.ClientConn,
		timeUsed:   time.Now(),
	}
	if c.unhealthy {
		wrapper.ClientConn.Close()
		wrapper.ClientConn = nil
	} else {
		wrapper.timeInitiated = c.timeInitiated
	}
	select {
	case c.pool.clients <- wrapper:
	default:
		return ErrFullPool
	}

	c.ClientConn = nil
	return nil
}

func (p *Pool) Capacity() int {
	if p.IsClosed() {
		return 0
	}
	return cap(p.clients)
}

func (p *Pool) Available() int {
	if p.IsClosed() {
		return 0
	}
	return len(p.clients)
}