package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/wujunhui99/looklook/app/usercenter/cmd/api/internal/svc"
	"github.com/wujunhui99/looklook/app/usercenter/cmd/api/internal/types"
	"github.com/wujunhui99/looklook/app/usercenter/cmd/rpc/usercenter"
	"github.com/wujunhui99/looklook/app/usercenter/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{
		AuthType: model.UserAuthTypeSystem,
		AuthKey:  req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.LoginResp{}
	_ = copier.Copy(resp, loginResp)

	return resp, nil

}
