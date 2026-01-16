package svc

import (
	"github.com/wujunhui99/looklook/app/order/cmd/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
