package homestay

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/travel/cmd/api/internal/types"
	"github.com/wujunhui99/looklook/app/travel/cmd/rpc/travel"
	"github.com/wujunhui99/looklook/pkg/tool"
	"github.com/wujunhui99/looklook/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// homestay room detail
func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req *types.HomestayDetailReq) (resp *types.HomestayDetailResp, err error) {
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get homestay detail fail"), " get homestay detail db err , id : %d , err : %v ", req.Id, err)
	}

	var typeHomestay types.Homestay
	if homestayResp.Homestay != nil {

		_ = copier.Copy(&typeHomestay, homestayResp.Homestay)

		typeHomestay.FoodPrice = tool.Fen2Yuan(homestayResp.Homestay.FoodPrice)
		typeHomestay.HomestayPrice = tool.Fen2Yuan(homestayResp.Homestay.HomestayPrice)
		typeHomestay.MarketHomestayPrice = tool.Fen2Yuan(homestayResp.Homestay.MarketHomestayPrice)

	}

	return &types.HomestayDetailResp{
		Homestay: typeHomestay,
	}, nil
}
