package message

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg/api"
	"github.com/zngue/go_open_platform/app/wechat"
	"io/ioutil"
)


func Message(ctx *gin.Context) {
	all, err := ioutil.ReadAll(ctx.Request.Body)
	platform := wechat.NewOpenPlatform()
	if err != nil {
		return
	}
	if len(all)==0{
		return
	}
	platform.DecryptMsg(all)
	ctx.JSON(200,"success")
}
func GetVerifyTicket(ctx *gin.Context)  {
	api.Success(ctx,api.Data(viper.GetString("wechatOpenPlatform.VerifyTicket")))
}
