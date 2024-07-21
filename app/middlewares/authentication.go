package middlewares

import (
	"errors"
	"net/http"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/models"
	"github.com/adamnasrudin03/go-market/configs"
	"github.com/adamnasrudin03/go-market/pkg/database"
	"github.com/adamnasrudin03/go-market/pkg/driver"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	cfg    = configs.GetInstance()
	logger = driver.Logger(cfg)
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := VerifyToken(c)
		if err != nil {
			response_mapper.RenderJSON(c.Writer, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		c.Set("userData", claims)
		c.Next()
	}
}

func AuthorizationMustBe() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationMustBeValidation(c)
		c.Next()
	}
}

func authorizationMustBeValidation(c *gin.Context) {
	var (
		db        = database.GetDB()
		userID    = c.MustGet("id").(uint64)
		userEmail = c.MustGet("email").(string)
		user      = models.User{}
	)

	err := db.Select("id").Where("id = ? AND email = ?", userID, userEmail).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || user.ID == 0 {
		err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Masuk kembali dengan user terdaftar",
				EN: "Log in again with registered user",
			},
		))
		response_mapper.RenderJSON(c.Writer, http.StatusUnauthorized, err)
		c.Abort()
		return
	}

	if err != nil {
		logger.Errorf("Failed to check user log in: %v ", err)
		response_mapper.RenderJSON(c.Writer, http.StatusUnauthorized, response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Gagal mengecek user log in",
				EN: "Failed to check user log in",
			},
		)))
		c.Abort()
		return
	}

}
