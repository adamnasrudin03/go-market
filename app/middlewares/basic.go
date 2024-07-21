package middlewares

import (
	"net/http"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/gin-gonic/gin"
)

func SetAuthBasic() gin.HandlerFunc {
	cfg := configs.GetInstance()
	return func(ctx *gin.Context) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := ctx.Request.BasicAuth()
		isValid := hasAuth && user == cfg.App.BasicUsername && password == cfg.App.BasicPassword
		if !isValid {
			err := response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
				response_mapper.MultiLanguages{
					ID: "Token tidak valid",
					EN: "Invalid token",
				},
			))

			response_mapper.RenderJSON(ctx.Writer, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
