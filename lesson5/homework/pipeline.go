package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

// ExecutePipeline создает пайплайн из функций Stage, который пропускает
// через себя данные из In, и выводит преобразованные данные в Out.
// Также работа функции может быть завершена через контекст ctx
func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := make(chan any)

	// создаем пайплайн
	pipeline := in
	for _, s := range stages {
		pipeline = s(pipeline)
	}

	// здесь горутина будет принимать данные из пайплайна и
	// записывать их в out
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
	}(ctx, pipeline, out)

	return out
}
