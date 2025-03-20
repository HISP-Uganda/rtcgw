package main

import (
	"log"
	"rtcgw/config"
	"rtcgw/tasks"

	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.RTCGwConf.Server.RedisAddress},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: config.RTCGwConf.Server.MaxConcurrent,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSendResults, tasks.HandleResultsTask)
	mux.HandleFunc(tasks.TypeCreateClient, tasks.HandleClientTask)
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
