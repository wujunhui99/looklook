package homestay

import (
	"net/http"

	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/logic/homestay"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// homestay room detail
func HomestayDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestay.NewHomestayDetailLogic(r.Context(), svcCtx)
		resp, err := l.HomestayDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
