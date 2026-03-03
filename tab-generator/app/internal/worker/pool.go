package worker

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

type Handler func(ctx context.Context) error

type Pool struct {
	sem *semaphore.Weighted
	wg  sync.WaitGroup
}

func NewPool(n int64) *Pool {
	return &Pool{
		sem: semaphore.NewWeighted(n),
	}
}

func (p *Pool) Submit(ctx context.Context, fn Handler) error {
	if err := p.sem.Acquire(ctx, 1); err != nil {
		return err
	}

	p.wg.Go(func() {
		defer p.sem.Release(1)

		_ = fn(ctx)
	})

	return nil
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
