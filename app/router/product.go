package router

import (
	"github.com/adamnasrudin03/go-market/app/controller"
	"github.com/adamnasrudin03/go-market/app/middlewares"
	"github.com/gin-gonic/gin"
)

func (r routes) productRouter(rg *gin.RouterGroup, handler controller.ProductController) {
	product := rg.Group("/products")
	{
		product.Use(middlewares.Authentication())
		product.GET("", middlewares.AuthorizationMustBe(), handler.GetList)
		product.GET("/:id", middlewares.AuthorizationMustBe(), handler.GetDetail)
	}

}
