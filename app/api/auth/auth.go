package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zngue/go_helper/pkg/api"
	"github.com/zngue/go_open_platform/app/wechat"
)

func AuthLink(ctx *gin.Context) {
	req := &wechat.AuthLinkRequest{}
	if err := ctx.ShouldBind(req); err != nil {
		api.DataWithErr(ctx, err, nil)
		return
	}
	if len(req.CallbackUrl) == 0 {
		api.Error(ctx, api.Msg("回调地址不能为空"))
		return
	}
	platform, err := wechat.NewOpenPlatform(true)
	if err != nil {
		api.DataWithErr(ctx, err, nil)
		return
	}
	auth, errs := platform.AuthLink(req)
	api.DataWithErr(ctx, errs, auth)
	return
}
func AuthLinkByCode(ctx *gin.Context) {
	code := ctx.DefaultQuery("code", "")
	if len(code) == 0 {
		api.DataWithErr(ctx, errors.New("参数错误"), nil)
		return
	}
	platform, err := wechat.NewOpenPlatform(false)

	if err != nil {
		api.DataWithErr(ctx, err, nil)
		return
	}
	byCode, errs := platform.GetLinkByCode(code)
	api.DataWithErr(ctx, errs, byCode)
	return

}
