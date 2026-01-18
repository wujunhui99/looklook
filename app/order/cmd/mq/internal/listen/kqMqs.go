package listen

import (
	"context"

	"github.com/wujunhui99/looklook/app/order/cmd/mq/internal/config"
	kqMq "github.com/wujunhui99/looklook/app/order/cmd/mq/internal/mqs/kq"
	"github.com/wujunhui99/looklook/app/order/cmd/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		kq.MustNewQueue(c.PaymentUpdateStatusConf, kqMq.NewPaymentUpdateStatusMq(ctx, svcContext)),
		//.....
	}

}
