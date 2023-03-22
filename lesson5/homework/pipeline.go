package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := make(Out)
	res := make(chan any)

	out = in
	for _, s := range stages {
		out = s(out)
	}

	go func(ctx context.Context, in In, out chan any) {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-in:
				if ok {
					out <- m
				} else {
					return
				}
			}
		}
	}(ctx, out, res)

	return res
}
