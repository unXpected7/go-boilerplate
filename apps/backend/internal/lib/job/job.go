package job

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/sriniously/go-boilerplate/apps/backend/internal/config"
)

type JobService struct {
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) *JobService {
	redisAddr := cfg.RedisAddress

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6, // Higher priority queue for important emails
				"default":  3, // Default priority for most emails
				"low":      1, // Lower priority for non-urgent emails
			},
		},
	)

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

func (j *JobService) Start() error {
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}
