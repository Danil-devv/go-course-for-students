package storage

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	// by default - 4
	maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: 10,
	}
}

func worker(ctx context.Context, ch chan Dir, wg *sync.WaitGroup, size *int64, count *int64, errgg *error) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-ch:
			if !ok {
				return
			}

			_, files, err := d.Ls(ctx)
			if err != nil {
				*errgg = fmt.Errorf("file does not exist: %v", err)
				return
			}
			for _, file := range files {
				s, err := file.Stat(ctx)
				if err != nil {
					*errgg = fmt.Errorf("file does not exist: %v", err)
					return
				}
				atomic.AddInt64(size, s)
				atomic.AddInt64(count, 1)
			}
		}
	}
}

func adder(ctx context.Context, ch chan Dir, wg *sync.WaitGroup, d Dir, addersCount *int64, maxWorkersCount int, errgg *error) {
	defer func() {
		atomic.AddInt64(addersCount, -1)
		wg.Done()
	}()
	//runtime.Gosched()

	queue := make([]Dir, 1)
	queue[0] = d
	ch <- d
	for {
		if len(queue) == 0 {
			return
		}

		dirs, _, err := queue[0].Ls(ctx)
		if err != nil {
			*errgg = err
			return
		}

		queue = queue[1:]

		if len(dirs) != 0 {
			for _, d := range dirs {
				if *addersCount < int64(maxWorkersCount)/2 {
					atomic.AddInt64(addersCount, 1)
					wg.Add(1)
					go adder(ctx, ch, wg, d, addersCount, maxWorkersCount, errgg)
				} else {
					queue = append(queue, d)
					ch <- d
				}
			}
		}
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	var (
		size, count int64
		wg          sync.WaitGroup
		wg2         sync.WaitGroup
		ch          chan Dir
		addersCount int64
		errgg       error
	)
	errgg = nil

	ctxWithCancel, cancel := context.WithCancel(ctx)

	ch = make(chan Dir, a.maxWorkersCount)

	for i := 1; i <= a.maxWorkersCount/2; i++ {
		wg.Add(1)
		go worker(ctxWithCancel, ch, &wg, &size, &count, &errgg)
	}

	wg2.Add(1)
	go adder(ctxWithCancel, ch, &wg2, d, &addersCount, a.maxWorkersCount, &errgg)

	wg2.Wait()

	close(ch)

	wg.Wait()

	cancel()

	return Result{
		size,
		count,
	}, errgg
}
