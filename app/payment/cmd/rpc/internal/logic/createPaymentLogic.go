package logic

import (
	"context"

	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/internal/svc"
	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建微信支付预处理订单
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreatePaymentResp{}, nil
}
