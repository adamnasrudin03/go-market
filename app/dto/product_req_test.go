package dto

import (
	"testing"

	"github.com/adamnasrudin03/go-market/app/models"
)

func TestProductListReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *ProductListReq
		wantErr bool
	}{
		{
			name: "invalid order by",
			m: &ProductListReq{
				BasedFilter: models.BasedFilter{
					OrderBy: "invalid",
					SortBy:  "",
				},
			},
			wantErr: true,
		},
		{
			name: "sort by required if order by provided",
			m: &ProductListReq{
				BasedFilter: models.BasedFilter{
					OrderBy: models.OrderByASC,
					SortBy:  "",
				},
			},
			wantErr: true,
		},
		{
			name:    "success",
			m:       &ProductListReq{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ProductListReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
