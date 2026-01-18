package kq

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/wujunhui99/looklook/app/order/cmd/mq/internal/svc"
	"github.com/wujunhui99/looklook/app/order/cmd/rpc/order"
	"github.com/wujunhui99/looklook/app/order/model"
	paymentModel "github.com/wujunhui99/looklook/app/payment/model"
	"github.com/wujunhui99/looklook/pkg/kqueue"
	"github.com/wujunhui99/looklook/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
*
Listening to the payment flow status change notification message queue
*/
type PaymentUpdateStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentUpdateStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentUpdateStatusMq {
	return &PaymentUpdateStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentUpdateStatusMq) Consume(ctx context.Context, _, val string) error {

	var message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(ctx).Error("PaymentUpdateStatusMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}

	if err := l.execService(ctx, message); err != nil {
		logx.WithContext(ctx).Error("PaymentUpdateStatusMq->execService  err : %v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *PaymentUpdateStatusMq) execService(ctx context.Context, message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage) error {

	orderTradeState := l.getOrderTradeStateByPaymentTradeState(message.PayStatus)
	if orderTradeState != -99 {
		//update homestay order state
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         message.OrderSn,
			TradeState: orderTradeState,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("update homestay order state fail"), "update homestay order state fail err : %v ,message:%+v", err, message)
		}
	}

	return nil
}

// Get order status based on payment status.
func (l *PaymentUpdateStatusMq) getOrderTradeStateByPaymentTradeState(thirdPaymentPayStatus int64) int64 {

	switch thirdPaymentPayStatus {
	case paymentModel.ThirdPaymentPayTradeStateSuccess:
		return model.HomestayOrderTradeStateWaitUse
	case paymentModel.ThirdPaymentPayTradeStateRefund:
		return model.HomestayOrderTradeStateRefund
	default:
		return -99
	}

}
