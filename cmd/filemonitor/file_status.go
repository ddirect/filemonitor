package main

import (
	"filemonitor/pkg/common"
	"time"
)

type FileStatus struct {
	common.FileInfo
	checkDelay   time.Duration
	isPersistent bool
}

const DEFAULT_CHECK_DELAY = 5 * time.Second
const MIN_CHECK_DELAY = 1 * time.Second
const MAX_CHECK_DELAY = 10 * time.Second

func (s FileStatus) WithMinDelay() FileStatus {
	s.checkDelay = MIN_CHECK_DELAY
	return s
}

func (s FileStatus) Schedule(sink chan<- FileStatus) {
	checkDelay := max(s.checkDelay, MIN_CHECK_DELAY)
	s.checkDelay = min(checkDelay*2, MAX_CHECK_DELAY)
	time.AfterFunc(checkDelay, func() { sink <- s })
}
