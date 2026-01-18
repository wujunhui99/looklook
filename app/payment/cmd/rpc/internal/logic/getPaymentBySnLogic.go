package logic

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/internal/svc"
	"github.com/wujunhui99/looklook/app/payment/cmd/rpc/pb"
	"github.com/wujunhui99/looklook/app/payment/model"
	"github.com/wujunhui99/looklook/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentBySnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentBySnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentBySnLogic {
	return &GetPaymentBySnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据sn查询流水记录
func (l *GetPaymentBySnLogic) GetPaymentBySn(in *pb.GetPaymentBySnReq) (*pb.GetPaymentBySnResp, error) {
	thirdPayment, err := l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GetPaymentBySn  FindOneBySn  db err:%v , in : %+v", err, in)
	}

	var resp pb.PaymentDetail
	if thirdPayment != nil {
		_ = copier.Copy(&resp, thirdPayment)

		resp.CreateTime = thirdPayment.CreateTime.Unix()
		resp.UpdateTime = thirdPayment.UpdateTime.Unix()
		resp.PayTime = thirdPayment.PayTime.Unix()
	}

	return &pb.GetPaymentBySnResp{
		PaymentDetail: &resp,
	}, nil
}
