package svc

import "github.com/wujunhui99/looklook/app/payment/cmd/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
