package task

import "sync"

type ItemHandler[I any] interface {
	HandleItem(I)
}

type Task[I any] struct {
	inChan chan I
	wg     sync.WaitGroup
}

func Start[I any](goCount int, channelDepth int, handlerFactory func() ItemHandler[I]) *Task[I] {
	t := &Task[I]{inChan: make(chan I)}
	t.wg.Add(goCount)
	for i := 0; i < goCount; i++ {
		go func() {
			defer t.wg.Done()
			handler := handlerFactory()
			for item := range t.inChan {
				handler.HandleItem(item)
			}
		}()
	}
	return t
}

func (t *Task[I]) StopAndWait() {
	close(t.inChan)
	t.wg.Wait()
}

func (t *Task[I]) InChan() chan<- I {
	return t.inChan
}
