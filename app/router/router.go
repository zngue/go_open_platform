package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zngue/go_open_platform/app/api/auth"
	"github.com/zngue/go_open_platform/app/api/message"
	"net/http"
)

func Router(router *gin.RouterGroup) {
	messageRouter := router.Group("message")
	{
		messageRouter.Handle(http.MethodGet, "info", message.Message)
		messageRouter.Handle(http.MethodPost, "info", message.Message)
		messageRouter.GET("ticket", message.GetVerifyTicket)

	}
	authRouter := router.Group("auth")

	{
		authRouter.GET("codeLink", auth.AuthLink)
		authRouter.GET("link", auth.AuthLinkByCode)

	}

}
