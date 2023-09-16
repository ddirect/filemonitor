package main

import (
	"filemonitor/common"
	"filemonitor/common/log"
	"sync"
)

type queuer struct {
	ctx           *Context
	wg            sync.WaitGroup
	newFileChan   chan string
	errorFileChan chan string
	queue         []string
	monitored     map[string]struct{}
}

func StartQueuer(ctx *Context, channelDepth int) *queuer {
	q := &queuer{
		ctx:           ctx,
		newFileChan:   make(chan string, channelDepth),
		errorFileChan: make(chan string, channelDepth),
		monitored:     make(map[string]struct{}),
	}
	q.wg.Add(1)
	go q.run()
	return q
}

func (q *queuer) NewFileChan() chan<- string {
	return q.newFileChan
}

func (q *queuer) ErrorFileChan() chan<- string {
	return q.errorFileChan
}

func (q *queuer) StopAndWait() {
	close(q.newFileChan)
	q.wg.Wait()
}

func (q *queuer) run() {
	defer q.wg.Done()
	for {
		if len(q.queue) > 0 {
			select {
			case q.ctx.CheckerChan <- FileStatus{FileInfo: common.FileInfo{FileName: q.queue[0]}}:
				q.queue = q.queue[1:]
			case name, ok := <-q.newFileChan:
				if !ok {
					return
				}
				q.handleNew(name)
			case name := <-q.errorFileChan:
				q.handleError(name)
			}
		} else {
			select {
			case name, ok := <-q.newFileChan:
				if !ok {
					return
				}
				q.handleNew(name)
			case name := <-q.errorFileChan:
				q.handleError(name)
			}
		}
	}
}

func (q *queuer) handleNew(name string) {
	log.Debug("NEW <- %s", name)
	if _, found := q.monitored[name]; found {
		return
	}
	q.monitored[name] = struct{}{}
	q.queue = append(q.queue, name)
}

func (q *queuer) handleError(name string) {
	log.Debug("ERR <- %s", name)
	delete(q.monitored, name)
}
