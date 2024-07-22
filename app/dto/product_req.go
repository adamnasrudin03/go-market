package dto

import (
	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-market/app/models"
)

type ProductDetailReq struct {
	ID           uint64 `json:"id" form:"id"`
	WarehouseID  uint64 `json:"warehouse_id" form:"warehouse_id"`
	CustomColumn string `json:"custom_column"`
}

type ProductListReq struct {
	Search string `json:"search" form:"search"`
	models.BasedFilter
}

func (m *ProductListReq) Validate() error {
	if m.Page <= 0 {
		m.Page = 1
	}

	if m.Limit <= 0 {
		m.Limit = 10
	}

	m.Search = help.ToLower(m.Search)

	m.OrderBy = help.ToUpper(m.OrderBy)
	if !models.IsValidOrderBy[m.OrderBy] && m.OrderBy != "" {
		return response_mapper.ErrInvalidFormat("order_by", "order_by")
	}

	m.SortBy = help.ToLower(m.SortBy)
	if m.OrderBy != "" && m.SortBy == "" {
		return response_mapper.ErrIsRequired("sort_by", "sort_by")
	}

	return nil
}
