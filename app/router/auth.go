package router

import (
	"github.com/adamnasrudin03/go-market/app/controller"
	"github.com/gin-gonic/gin"
)

func (r routes) authRouter(rg *gin.RouterGroup, handler controller.AuthController) {
	auth := rg.Group("/auth")
	{
		auth.POST("/sign-up", handler.Register)
		auth.POST("/sign-in", handler.Login)
	}

}
