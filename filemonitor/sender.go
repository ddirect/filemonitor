package main

import (
	"filemonitor/common/log"
)

func NewSender(ctx *Context) *sender {
	return &sender{ctx: ctx}
}

type sender struct {
	ctx *Context
}

func (s *sender) HandleItem(status FileStatus) {
	log.Debug("SND <- %v", status)
	log.Info("send %s", status.FileName)
	if status.Id == "" {
		if id, err := s.ctx.Db.CreateDocument(status); err != nil {
			log.Error("CreateDocument: %s", err)
		} else {
			status.Id = id
		}
	} else {
		if err := s.ctx.Db.UpdateDocumentById(status.Id, status.WithoutId()); err != nil {
			log.Error("UpdateDocumentBy: %s", err)
		}
	}
	status.WithMinDelay().Schedule(s.ctx.CheckerChan)
}
