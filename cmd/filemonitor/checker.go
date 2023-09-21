package main

import (
	"filemonitor/pkg/log"
	"io"
	"os"
	"path/filepath"
)

func NewChecker(ctx *Context) *checker {
	return &checker{ctx: ctx}
}

type checker struct {
	ctx *Context
}

func (c *checker) HandleItem(status FileStatus) {
	log.Debug("CHK <- %v", status)
	info, err := os.Lstat(status.FileName)
	if err != nil {
		c.ctx.HandleFileError(err, "checker", status)
		return
	}
	if !info.Mode().IsRegular() && !info.Mode().IsDir() {
		c.ctx.HandleFileError(nil, "", status)
		return
	}
	newModTime := info.ModTime().UnixNano()
	if newModTime == status.ModTime {
		status.Schedule(c.ctx.CheckerChan)
	} else {
		status.ModTime = newModTime
		if info.IsDir() {
			if err := c.scanDir(status.FileName); err != nil {
				c.ctx.HandleFileError(err, "checker", status)
			} else {
				status.WithMinDelay().Schedule(c.ctx.CheckerChan)
			}
		} else {
			c.ctx.HasherChan <- status
		}
	}
}

func (c *checker) scanDir(dirName string) error {
	dir, err := os.Open(dirName)
	log.Info("rdir %s", dirName)
	if err != nil {
		return err
	}
	defer dir.Close()
	for {
		names, err := dir.Readdirnames(100)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		for _, name := range names {
			c.ctx.NewFileChan <- filepath.Join(dirName, name)
		}
	}
}
