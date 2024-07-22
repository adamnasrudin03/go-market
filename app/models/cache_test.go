package models

import "testing"

func TestKeyCacheProductDetail(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				id: 1,
			},
			want: "product_detail_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeyCacheProductDetail(tt.args.id); got != tt.want {
				t.Errorf("KeyCacheProductDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
