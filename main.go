package main

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/samber/do"
	"github.com/valli0x/apod/di"
	"github.com/valli0x/apod/server"
	"github.com/valli0x/apod/worker"
	"github.com/valli0x/apod/worker/fairshare"
)

func main() {
	i := do.New()
	do.Provide(i, di.NewConfig)
	do.Provide(i, di.NewLogger)
	do.Provide(i, di.NewStor)
	do.Provide(i, di.NewServer)
	do.Provide(i, di.NewWorker)
	do.Provide(i, di.NewJobManager)

	logger := do.MustInvoke[*zerolog.Logger](i)
	worker := do.MustInvoke[*worker.APODjob](i)
	server := do.MustInvoke[*server.Server](i)
	jobmanager := do.MustInvoke[*fairshare.JobManager](i)
	queueID := "one_queue"

	logger.Info().Msg("start worker")

	jobmanager.Start()
	defer jobmanager.Stop()
	
	go func ()  {
		jobmanager.AddJob(worker, queueID)
		for range time.NewTicker(time.Duration(24) * time.Hour).C {
			jobmanager.AddJob(worker, queueID)
		}
	}()

	logger.Info().Msg("server start")

	if err := server.Start(); err != nil {
		panic(err)
	}
}
