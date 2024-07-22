package controller

import (
	"net/http"
	"strconv"
	"strings"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/dto"
	"github.com/adamnasrudin03/go-market/app/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ProductController interface {
	GetDetail(ctx *gin.Context)
	GetList(ctx *gin.Context)
}

type productController struct {
	Service   service.ProductService
	Logger    *logrus.Logger
	Validator *validator.Validate
}

func NewProductController(srv service.ProductService, logger *logrus.Logger, validator *validator.Validate) ProductController {
	return &productController{
		Service:   srv,
		Logger:    logger,
		Validator: validator,
	}
}

func (c *productController) GetDetail(ctx *gin.Context) {
	var (
		opName  = "ProductController-GetDetail"
		idParam = strings.TrimSpace(ctx.Param("id"))
		err     error
	)

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Logger.Errorf("%v error parse param: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.ErrInvalid("ID Produk", "Product ID"))
		return
	}

	res, err := c.Service.GetByID(ctx, id)
	if err != nil {
		response_mapper.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(ctx.Writer, http.StatusOK, res)
}

func (c *productController) GetList(ctx *gin.Context) {
	var (
		opName = "ProductController-GetList"
		input  dto.ProductListReq
		err    error
	)

	err = ctx.ShouldBindQuery(&input)
	if err != nil {
		c.Logger.Errorf("%v error bind json: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusBadRequest, response_mapper.ErrGetRequest())
		return
	}

	res, err := c.Service.GetList(ctx, input)
	if err != nil {
		c.Logger.Errorf("%v error: %v ", opName, err)
		response_mapper.RenderJSON(ctx.Writer, http.StatusInternalServerError, err)
		return
	}

	response_mapper.RenderJSON(ctx.Writer, http.StatusOK, res)
}
