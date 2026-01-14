package homestayBussiness

import (
	"net/http"

	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/logic/homestayBussiness"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// boss detail
func HomestayBussinessDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayBussinessDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayBussiness.NewHomestayBussinessDetailLogic(r.Context(), svcCtx)
		resp, err := l.HomestayBussinessDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
