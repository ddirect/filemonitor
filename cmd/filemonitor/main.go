package main

import (
	"filemonitor/pkg/common"
	"filemonitor/pkg/log"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
)

const CHANNEL_DEPTH = 10

func main() {
	var basePath, logLevel, apiKey, ledgerName, collectionName string

	flag.StringVar(&basePath, "base", ".", "root direcotry")
	flag.StringVar(&logLevel, "log_level", "", "log level")
	flag.StringVar(&apiKey, "api_key", "", "API key")
	flag.StringVar(&ledgerName, "ledger", "default", "immudb ledger")
	flag.StringVar(&collectionName, "collection", common.FILE_MONITOR_DEFAULT_COLLECTION, "immudb collection")
	flag.Parse()

	log.SetLevel(log.LevelFromString(logLevel))

	base, err := filepath.Abs(basePath)
	if err != nil {
		log.Error("main: %s", err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	db, err := SetupDB(ledgerName, collectionName, apiKey)
	if err != nil {
		log.Error("setup DB: %s", err)
		return
	}

	ctx := &Context{Db: db}
	queuer := StartQueuer(ctx, CHANNEL_DEPTH)
	cpuCount := runtime.NumCPU()
	checkerTask := Start(cpuCount, func() ItemHandler { return NewChecker(ctx) })
	hasherTask := Start(2, func() ItemHandler { return NewHasher(ctx) })
	senderTask := Start(1, func() ItemHandler { return NewSender(ctx) })

	ctx.NewFileChan = queuer.NewFileChan()
	ctx.ErrorFileChan = queuer.ErrorFileChan()
	ctx.CheckerChan = checkerTask.InChan()
	ctx.HasherChan = hasherTask.InChan()
	ctx.SenderChan = senderTask.InChan()

	log.Info("starting on %s", base)
	ctx.CheckerChan <- FileStatus{FileInfo: common.FileInfo{FileName: base}, isPersistent: true}

	sig := <-sigCh
	log.Info("%s received - cleaning up...", sig)

	queuer.StopAndWait()
	checkerTask.StopAndWait()
	hasherTask.StopAndWait()
	senderTask.StopAndWait()
}
