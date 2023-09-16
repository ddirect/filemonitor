package main

import (
	"filemonitor/common"
	"filemonitor/common/log"
	"filemonitor/immudb"
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
	if _, err := s.ctx.Db.UpdateDocumentBy(common.FILE_INFO_FILE_FIELD_NAME, status.FileName, status); err != nil {
		if immudb.HttpStatusCode(err) != 404 {
			log.Error("UpdateDocumentBy: %s", err)
		} else if _, err := s.ctx.Db.CreateDocument(status); err != nil {
			log.Error("CreateDocument: %s", err)
		}
	}
	status.WithMinDelay().Schedule(s.ctx.CheckerChan)
}
