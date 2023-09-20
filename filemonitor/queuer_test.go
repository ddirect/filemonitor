package main

import (
	cm "filemonitor/common_test"
	"sync"
	"testing"
	"time"
)

func newQueuer(channelDepth int) (*queuer, <-chan FileStatus) {
	checkerChan := make(chan FileStatus, channelDepth)
	return StartQueuer(&Context{CheckerChan: checkerChan}, channelDepth), checkerChan
}

func TestQueuerFillAndEmpty(t *testing.T) {
	q, outChan := newQueuer(3)
	defer q.StopAndWait()

	var sent []string
	var received []FileStatus

	core := func(count int, sendDelay, recvDelay time.Duration) {
		var wg sync.WaitGroup
		wg.Add(1)
		// sender
		go func() {
			defer wg.Done()
			for i := 0; i < count; i++ {
				name := cm.RandomString()
				sent = append(sent, name)
				q.NewFileChan() <- name
				time.Sleep(sendDelay)
			}
		}()

		// receiver
		for i := 0; i < count; i++ {
			received = append(received, <-outChan)
			time.Sleep(recvDelay)
		}
		wg.Wait()
	}

	// fill
	core(100, 10*time.Millisecond, 20*time.Millisecond)
	// empty
	core(200, 20*time.Millisecond, 10*time.Millisecond)

	// check that the queue is empty
	found := false
	select {
	case <-outChan:
		found = true
	default:
	}

	if found {
		t.Fatal("queue has unexpected items")
	}

	if len(sent) != len(received) {
		t.Fatal("sent and received count mismatch")
	}

	for i := range sent {
		if sent[i] != received[i].FileName {
			t.Fatal("queue has wrong items")
		}
	}
}

func TestQueuerItemsAreNotRepeatedUntilError(t *testing.T) {
	q, outChan := newQueuer(0)
	defer q.StopAndWait()

	expect := func(x string) {
		if x != (<-outChan).FileName {
			t.FailNow()
		}
	}

	q.NewFileChan() <- "a"
	q.NewFileChan() <- "b"
	q.NewFileChan() <- "a"
	q.NewFileChan() <- "c"
	q.ErrorFileChan() <- "a"
	q.NewFileChan() <- "a"
	q.NewFileChan() <- "b"
	q.NewFileChan() <- "d"

	expect("a")
	expect("b")
	expect("c")
	expect("a")
	expect("d")
}
