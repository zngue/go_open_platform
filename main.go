package main


import (
	"github.com/gin-gonic/gin"
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

		}),
		common_run.IsRegisterCenter(true),
	)
}


