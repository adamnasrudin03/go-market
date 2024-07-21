package middlewares

import (
	"fmt"
	"strings"
	"time"

	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(params JWTClaims) (string, error) {
	now := time.Now().Unix()
	expired := time.Now().AddDate(0, 0, cfg.App.ExpiredToken).Unix()

	claims := &JWTClaims{
		ID:    params.ID,
		Name:  params.Name,
		Phone: params.Phone,
		Email: params.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprintf("%d", params.ID),
			IssuedAt:  now,
			NotBefore: now,
			ExpiresAt: expired,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.App.SecretKey))
	if err != nil {
		logger.Errorf("failed generate token: %v", err)
		err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Gagal melakukan signed token string",
				EN: "Failed to get signed token string",
			},
		))
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(ctx *gin.Context) (tokenString string, err error) {
	headerToken := ctx.Request.Header.Get("Authorization")
	bearer := strings.HasPrefix(headerToken, "Bearer")

	if !bearer {
		err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Bearer token tidak ditemukan",
				EN: "Bearer token not found",
			},
		))
		return "", err
	}

	if len(headerToken) <= 7 {
		err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Token tidak ditemukan",
				EN: "Token not found",
			},
		))
		return "", err
	}

	return headerToken[7:], nil
}

func VerifyToken(ctx *gin.Context) (interface{}, error) {
	tokenString, err := ExtractToken(ctx)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Errorf("failed unexpected signing method : %v", err)
			err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
				response_mapper.MultiLanguages{
					ID: fmt.Sprintf("metode penandatanganan yang tidak terduga: %v", token.Header["alg"]),
					EN: fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
				},
			))
			return nil, err
		}
		return []byte(cfg.App.SecretKey), nil
	})

	if _, ok := err.(*jwt.ValidationError); ok {
		errCustom := extractErrorJwtValidation(err, tokenString)
		if errCustom != nil {
			return nil, errCustom
		}
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Gagal melakukan parse claims",
				EN: "Failed to parse claims",
			},
		))
		return nil, err
	}

	if ok && token.Valid {
		ctx.Set("id", claims.ID)
		ctx.Set("name", claims.Name)
		ctx.Set("phone", claims.Phone)
		ctx.Set("email", claims.Email)
		ctx.Set("expired_at", claims.ExpiresAt)

		return token.Claims.(*JWTClaims), nil
	}

	return nil, err

}

func extractErrorJwtValidation(err error, tokenString string) error {
	if ve, ok := err.(*jwt.ValidationError); ok {
		var errorCustom error
		msgErr := err.Error()
		msgErrID, errTrs := help.Translate(err.Error(), help.Auto, help.LangID)
		if errTrs != nil || msgErrID == "" {
			msgErrID = msgErr
		}

		errorCustom = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: fmt.Sprintf("Gagal melakukan parsing token, Error: %s", msgErrID),
				EN: fmt.Sprintf("Failed to parse token, Error: %s", msgErr),
			},
		))
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errorCustom = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
				response_mapper.MultiLanguages{
					ID: fmt.Sprintf("Format token tidak benar: %s ", tokenString),
					EN: fmt.Sprintf("Invalid token format: %s", tokenString),
				},
			))
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			errorCustom = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
				response_mapper.MultiLanguages{
					ID: "Token telah kadaluarsa",
					EN: "Token has been expired",
				},
			))
		}

		return errorCustom
	}

	logger.Errorf("unknown token error : %v", err)
	err = response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
		response_mapper.MultiLanguages{
			ID: "Token tidak benar",
			EN: "Invalid token",
		},
	))
	return err

}
