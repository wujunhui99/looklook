package thirdPayment

import (
	"net/http"

	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/logic/thirdPayment"
	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/payment/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// third payment：wechat pay callback
func ThirdPaymentWxPayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentWxPayCallback(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
