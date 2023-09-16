package main

import (
	"filemonitor/common"
	"filemonitor/common/log"
	"filemonitor/filemonitor/task"
)

type Task = task.Task[FileStatus]
type ItemHandler = task.ItemHandler[FileStatus]

func Start(goCount int, factory func() ItemHandler) *Task {
	return task.Start[FileStatus](goCount, CHANNEL_DEPTH, factory)
}

type Context struct {
	NewFileChan   chan<- string
	ErrorFileChan chan<- string
	CheckerChan   chan<- FileStatus
	HasherChan    chan<- FileStatus
	SenderChan    chan<- FileStatus
	Db            Vault
}

func (c *Context) HandleFileError(err error, caller string, status FileStatus) {
	if err != nil {
		log.Error("%s: %s", caller, err)
	}
	if status.isPersistent {
		FileStatus{
			FileInfo:     common.FileInfo{FileName: status.FileName},
			isPersistent: true,
			checkDelay:   DEFAULT_CHECK_DELAY,
		}.Schedule(c.CheckerChan)
	} else {
		c.ErrorFileChan <- status.FileName
	}
}
