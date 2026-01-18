package listen

import (
	"context"

	"github.com/wujunhui99/looklook/app/order/cmd/mq/internal/config"
	"github.com/wujunhui99/looklook/app/order/cmd/mq/internal/svc"
	"github.com/zeromicro/go-zero/core/service"
)

func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	//kq ：pub sub
	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services
}
