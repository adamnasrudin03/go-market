package controller

import (
	"net/http"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authController struct {
	Service   service.AuthService
	Logger    *logrus.Logger
	Validator *validator.Validate
}

func NewAuthController(srv service.AuthService, logger *logrus.Logger, validator *validator.Validate) AuthController {
	return &authController{
		Service:   srv,
		Logger:    logger,
		Validator: validator,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var (
		opName = "AuthController-Register"
		input  dto.RegisterReq
	)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}

	err = c.Validator.Struct(input)
	if err != nil {
		c.Logger.Errorf("%v error validate struct: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.FormatValidationError(err))
		return
	}

	res, err := c.Service.Register(ctx, input)
	if err != nil {
		response_mapper.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(ctx.Writer, http.StatusCreated, res)
}

func (c *authController) Login(ctx *gin.Context) {
	var (
		opName = "AuthController-Login"
		input  dto.LoginReq
	)

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}

	err = c.Validator.Struct(input)
	if err != nil {
		c.Logger.Errorf("%v error validate struct: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.FormatValidationError(err))
		return
	}

	res, err := c.Service.Login(ctx, input)
	if err != nil {
		response_mapper.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(ctx.Writer, http.StatusOK, res)
}
