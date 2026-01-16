package homestayOrder

import (
	"net/http"

	"github.com/wujunhui99/looklook/app/order/cmd/api/internal/logic/homestayOrder"
	"github.com/wujunhui99/looklook/app/order/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/order/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建民宿订单
func CreateHomestayOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateHomestayOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayOrder.NewCreateHomestayOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateHomestayOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
