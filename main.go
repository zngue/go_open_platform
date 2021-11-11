package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/util"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg/common_run"
	"github.com/zngue/go_helper/pkg/response"
	"github.com/zngue/go_open_platform/app/router"
)

func main() {
	run()
}
func run() {
	common_run.CommonGinRun(
		common_run.FnRouter(func(engine *gin.Engine) {
			engine.NoRoute(func(c *gin.Context) {
				response.HttpFailWithCodeAndMessage(404, "路由不存在", c)
			})
			platform := engine.Group("platform")
			router.Router(platform)
			engine.GET("test", func(context *gin.Context) {
				json, err := util.PostJSON("https://api.weixin.qq.com/cgi-bin/component/api_start_push_ticket", map[string]interface{}{
					"component_appid":  viper.GetString("wechatOpenPlatform.Appid"),
					"component_secret": viper.GetString("wechatOpenPlatform.AppSecret"),
				})
				fmt.Println(json, err)
			})

		}),
		common_run.IsRegisterCenter(true),
	)
}
