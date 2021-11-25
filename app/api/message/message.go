package message

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg/api"
	"github.com/zngue/go_open_platform/app/service"
	"github.com/zngue/go_open_platform/app/wechat"
	"io/ioutil"
)

func Message(ctx *gin.Context) {
	all, err := ioutil.ReadAll(ctx.Request.Body)
	platform, err := wechat.NewOpenPlatform(false)
	if err != nil {
		return
	}
	if len(all) == 0 {
		return
	}
	platform.DecryptMsg(all)
	ctx.JSON(200, "success")
}
func GetVerifyTicket(ctx *gin.Context) {
	api.Success(ctx, api.Data(viper.GetString("wechatOpenPlatform.VerifyTicket")))
}

func Token(ctx *gin.Context) {

	appid := ctx.Param("appid")
	if len(appid)==0 {
		api.Error(ctx,api.Msg("account id is null "))
		return
	}
	var req service.OfficialAccountRequest
	req.Appid=appid
	account, err3 := service.NewOfficialAccount().Detail(&req)
	if err3 != nil {
		api.DataWithErr(ctx,err3,nil)
		return
	}
	platform, err := wechat.NewOpenPlatform(true)
	if err != nil {
		return
	}
	platform.Open(ctx,account.Appid)
	ctx.JSON(200,"success")

}
func Parse(ctx *gin.Context) {
	platform, err := wechat.NewOpenPlatform(true)
	fmt.Println(err)
	msg:=`<xml>
    <ToUserName><![CDATA[gh_6e09d08362fb]]></ToUserName>
    <Encrypt><![CDATA[09eAPYcop1RmHWDsX5+pSGm72dxbN+cflY96jmFgBihqDrTpzaVZs8vEA++jzmFvuGcKNmtV0rvs3NbgATqAZsBF/TI5YC/5Vc1P/tXZRkct0XyhuUUJ70hWq0QU4svLTKuVW/IuL01ZBhzvmp3RwYcE+xCN7LavjWm7LvHwV+uaF08uHmvwZyrJGdbjxxMGBOt8lcZ1jj98VP6oStJ2nWagcQiUV5JYpm1xA8N1VVHnYJefD2n4kQQd0jJWQv5MrMOO4G5jvvwyrGtnJ5shMqK7VMurZ4fVDowpbcH8BqaSLyAME2jXoVXe17+JCkptjUkEThTdjQvY2lmBCRn6T99XNqGj2nWWdodWjQAKOz+7SrEG3tLcaUsuQ4VLcwXWH65ORs1wsDrv3ysOFsean7+VuexkHdRyKkTs538HBZU=]]></Encrypt>
</xml>`

	platform.DecryptMsg([]byte(msg))



}
