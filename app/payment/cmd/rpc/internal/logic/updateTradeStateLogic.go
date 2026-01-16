package logic

import (
	"context"

	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/internal/svc"
	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateTradeStateResp{}, nil
}
