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

type GetPaymentSuccessRefundByOrderSnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentSuccessRefundByOrderSnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentSuccessRefundByOrderSnLogic {
	return &GetPaymentSuccessRefundByOrderSnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据订单sn查询流水记录
func (l *GetPaymentSuccessRefundByOrderSnLogic) GetPaymentSuccessRefundByOrderSn(in *pb.GetPaymentSuccessRefundByOrderSnReq) (*pb.GetPaymentSuccessRefundByOrderSnResp, error) {
	whereBuilder := l.svcCtx.ThirdPaymentModel.SelectBuilder().Where(
		"order_sn = ? and (trade_state = ? or trade_state = ? )",
		in.OrderSn, model.ThirdPaymentPayTradeStateSuccess, model.ThirdPaymentPayTradeStateRefund,
	)
	thirdPayments, err := l.svcCtx.ThirdPaymentModel.FindAll(l.ctx, whereBuilder, "id desc")
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrMsg("get payment record fail"), "get payment record fail FindOneByQuery  err : %v , orderSn:%s", err, in.OrderSn)
	}

	var resp pb.PaymentDetail
	if len(thirdPayments) > 0 {
		thirdPayment := thirdPayments[0]
		if thirdPayment != nil {
			_ = copier.Copy(&resp, thirdPayment)
			resp.CreateTime = thirdPayment.CreateTime.Unix()
			resp.UpdateTime = thirdPayment.UpdateTime.Unix()
			resp.PayTime = thirdPayment.PayTime.Unix()
		}
	}

	return &pb.GetPaymentSuccessRefundByOrderSnResp{
		PaymentDetail: &resp,
	}, nil
}
