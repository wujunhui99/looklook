package svc

import (
	"github.com/hibiken/asynq"
	"github.com/wujunhui99/looklook/app/mqueue/cmd/scheduler/internal/config"
)

type ServiceContext struct {
	Config config.Config

	Scheduler *asynq.Scheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Scheduler: newScheduler(c),
	}
}
