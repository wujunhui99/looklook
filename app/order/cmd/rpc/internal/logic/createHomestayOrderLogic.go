package logic

import (
	"context"

	"github.com/wujunhui99/looklook/app/order/cmd/rpc/internal/svc"
	"github.com/wujunhui99/looklook/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿下订单
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {

	return &pb.CreateHomestayOrderResp{}, nil
}
