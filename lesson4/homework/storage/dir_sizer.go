package storage

import (
	"context"
	"fmt"
	"sync"
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
		maxWorkersCount: 4,
	}
}

// worker обрабатывает указатели на массивы файлов из канала filesChan.
// увеличивает кол-во обработанных файлов count и общий вес файлов count,
// используя атомарные инструкции
func worker(ctx context.Context, filesChan chan *[]File, wg *sync.WaitGroup, m *sync.Mutex,
	size *int64, count *int64, runtimeError *error) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// если контекст закрылся, то завершаем работу
			return
		case files, ok := <-filesChan:
			// если канал закрылся, то завершаем работу
			if !ok {
				return
			}

			// проходим по всем файлам
			for _, file := range *files {
				s, err := file.Stat(ctx)
				// если произошла ошибка при работе с файлом,
				// то записываем ошибку в runtimeError и останавливаем горутину
				if err != nil {
					m.Lock()
					*runtimeError = fmt.Errorf("file does not exist: %v", err)
					m.Unlock()
					return
				}
				// увеличиваем size и count
				m.Lock()
				*size += s
				*count++
				m.Unlock()
			}
		}
	}
}

// dirsAdder это горутина, которая добавляет файлы из директории d в канал файлов.
// При этом она может вызывать себя рекурсивно,
// если dirAddersCount не превосходит maxWorkersCount
func dirsAdder(ctx context.Context, filesChan chan *[]File, wg *sync.WaitGroup, m *sync.Mutex,
	dir Dir, dirAddersCount *int64, maxWorkersCount int, runtimeError *error) {
	// после завершения необходимо уменьшить счетчик работающих adder'ов на 1
	defer func() {
		m.Lock()
		*dirAddersCount--
		m.Unlock()
		wg.Done()
	}()

	// очередь для директорий
	// будет использоваться для обхода директории в ширину
	queue := make([]Dir, 1)
	queue[0] = dir
	for len(queue) != 0 {
		dirs, files, err := queue[0].Ls(ctx)
		// передаем в канал новые файлы
		// здесь может произойти deadlock, если не запущен хотя бы 1 worker,
		// который читает этот канал
		filesChan <- &files
		if err != nil {
			m.Lock()
			*runtimeError = fmt.Errorf("file does not exist: %v", err)
			m.Unlock()
			return
		}

		// удаляем первый элемент из очереди
		// т.к. он уже обработан
		queue = queue[1:]

		if len(dirs) != 0 {
			for _, d := range dirs {
				m.Lock()
				if *dirAddersCount < int64(maxWorkersCount) {
					// Если у нас есть пространство для создания еще одной горутины, то мы создаём ее
					*dirAddersCount++
					wg.Add(1)
					go dirsAdder(ctx, filesChan, wg, m, d, dirAddersCount, maxWorkersCount, runtimeError)
				} else {
					// Если мы не можем создать еще одну горутину, то тогда
					// функция сама обработает все папки из директории
					queue = append(queue, d)
				}
				m.Unlock()
			}
		}
	}
}

// Size возвращает структуру Result, которая характеризует директорию d.
// Также функция может вернуть ошибку, которая произошла во время работы
func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	var (
		size, count  int64          // размер и кол-во файлов соответственно
		workersWG    sync.WaitGroup // wait группа worker'ов
		addersWG     sync.WaitGroup // wait группа adder'ов
		filesChan    chan *[]File   // поток с файлами
		addersCount  int64          // кол-во горутин-adder'ов
		runtimeError error          // сюда будем записывать ошибку, которая могла произойти в рантайме
		m            sync.Mutex     // мьютекс для блокировки записи в runtimeError
	)

	// контекст с отменой используется в горутинах worker и adder
	ctxWithCancel, cancel := context.WithCancel(ctx)

	// поток файлов является буферизированным каналом с размером максимального кол-ва горутин
	filesChan = make(chan *[]File, a.maxWorkersCount)

	// половину горутин выделяем для worker'ов
	for i := 1; i <= a.maxWorkersCount/2; i++ {
		workersWG.Add(1)
		go worker(ctxWithCancel, filesChan, &workersWG, &m, &size, &count, &runtimeError)
	}

	// запускаем adder, который будет добавлять в канал файлов новые файлы
	// так же dirsAdder может рекурсивно вызывать себя в качестве горутины,
	// поэтому мы оставили свободное место размером maxWorkersCount/2
	addersWG.Add(1)
	go dirsAdder(ctxWithCancel, filesChan, &addersWG, &m, d, &addersCount, a.maxWorkersCount/2, &runtimeError)

	addersWG.Wait()
	// когда все adder'ы закончили работу и добавили все файлы из директории в канал,
	// то его можно закрыть, тем самым обеспечив возможность выхода worker'ам
	close(filesChan)

	// если в рантайме произошла какая-то ошибка,
	// то завершаем контекст со всеми работающими горутинами
	m.Lock()
	switch runtimeError {
	case nil:
		m.Unlock()
	default:
		cancel()
		m.Unlock()
	}

	// ожидаем, пока все worker'ы закончат свою работу
	workersWG.Wait()
	cancel()

	return Result{
		size,
		count,
	}, runtimeError
}
