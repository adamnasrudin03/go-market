package models

const (
	TRUE  = "true"
	FALSE = "false"

	// Default ...
	DefaultPage                 = 1
	DefaultLimit                = 10
	DefaultLimitIsTotalDataTrue = 20

	// OrderBy ...
	OrderByASC  = "ASC"
	OrderByDESC = "DESC"
)

var (
	IsValidOrderBy = map[string]bool{
		OrderByASC:  true,
		OrderByDESC: true,
	}
)

type BasedFilter struct {
	Limit             int    `json:"limit" form:"limit"`
	Offset            int    `json:"offset" form:"offset"`
	Page              int    `json:"page" form:"page"`
	OrderBy           string `json:"order_by" form:"order_by"`
	SortBy            string `json:"sort_by" form:"sort_by"`
	IsNoLimit         bool   `json:"is_no_limit" form:"is_no_limit"`
	IsNotDefaultQuery bool   `json:"is_not_default_query" form:"is_not_default_query"`
	CustomColumns     string `json:"custom_columns" form:"custom_columns"`
}

func (c *BasedFilter) DefaultQuery() BasedFilter {
	if c.Limit <= 0 {
		c.Limit = 10
	}

	if c.Page <= 0 {
		c.Page = 1
	}

	if c.Page > 0 {
		c.Offset = (c.Page - 1) * c.Limit
	}

	return *c
}
