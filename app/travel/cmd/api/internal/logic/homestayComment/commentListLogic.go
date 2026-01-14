package homestayComment

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/types"
	"github.com/wujunhui99/looklook/app/usercenter/cmd/rpc/usercenter"
	"github.com/wujunhui99/looklook/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay comment list
func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	// 构建查询
	builder := l.svcCtx.HomestayCommentModel.SelectBuilder()

	// 使用游标分页获取评论列表（按ID降序）
	comments, err := l.svcCtx.HomestayCommentModel.FindPageListByIdDESC(l.ctx, builder, req.LastId, req.PageSize)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "CommentList FindPageListByIdDESC db fail, lastId: %d, pageSize: %d, err: %v", req.LastId, req.PageSize, err)
	}

	var list []types.HomestayComment
	for _, comment := range comments {
		// 计算平均星级
		star := l.calculateAverageStar(comment.Star)

		// 获取用户信息
		var nickname, avatar string
		userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
			Id: comment.UserId,
		})
		if err != nil {
			l.Logger.Errorf("get user info fail, userId: %d, err: %v", comment.UserId, err)
		} else if userResp.User != nil && userResp.User.Id > 0 {
			nickname = userResp.User.Nickname
			avatar = userResp.User.Avatar
		}

		list = append(list, types.HomestayComment{
			Id:         comment.Id,
			HomestayId: comment.HomestayId,
			Content:    comment.Content,
			Star:       star,
			UserId:     comment.UserId,
			Nickname:   nickname,
			Avatar:     avatar,
		})
	}

	return &types.CommentListResp{
		List: list,
	}, nil
}

// calculateAverageStar 计算多维度星级的平均值
func (l *CommentListLogic) calculateAverageStar(starStr string) float64 {
	if starStr == "" {
		return 0
	}

	// star 字段格式可能是多个星级用逗号分隔，如 "4.5,5.0,4.0"
	parts := strings.Split(starStr, ",")
	if len(parts) == 0 {
		return 0
	}

	var total float64
	var count int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		star, err := strconv.ParseFloat(part, 64)
		if err != nil {
			continue
		}
		total += star
		count++
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}
