package main

import (
	"crypto/sha256"
	"filemonitor/pkg/log"
	"hash"
	"io"
	"os"
)

func NewHasher(ctx *Context) *hasher {
	return &hasher{
		ctx:        ctx,
		hashEngine: sha256.New(),
		buffer:     make([]byte, 0x20000),
	}
}

type hasher struct {
	ctx        *Context
	hashEngine hash.Hash
	buffer     []byte
}

func (h *hasher) HandleItem(status FileStatus) {
	log.Debug("HAS <- %v", status)
	file, err := os.Open(status.FileName)
	if err != nil {
		h.ctx.HandleFileError(err, "hasher", status)
		return
	}
	defer file.Close()
	log.Info("hash %s", status.FileName)
	h.hashEngine.Reset()
	for {
		done, err := file.Read(h.buffer)
		if err != nil {
			if err == io.EOF {
				status.Hash = h.hashEngine.Sum(nil)
				h.ctx.SenderChan <- status
			} else {
				h.ctx.HandleFileError(err, "hasher", status)
			}
			return
		}
		h.hashEngine.Write(h.buffer[:done])
	}
}
