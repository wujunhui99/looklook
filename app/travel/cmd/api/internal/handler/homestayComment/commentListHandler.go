package homestayComment

import (
	"net/http"

	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/logic/homestayComment"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// homestay comment list
func CommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayComment.NewCommentListLogic(r.Context(), svcCtx)
		resp, err := l.CommentList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
