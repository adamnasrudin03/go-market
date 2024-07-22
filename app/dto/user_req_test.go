package dto

import "testing"

func TestUserDetailReq_Validate(t *testing.T) {
	tests := []struct {
		name    string
		m       *UserDetailReq
		wantErr bool
	}{
		{
			name:    "invalid params",
			m:       &UserDetailReq{},
			wantErr: true,
		},
		{
			name: "success",
			m: &UserDetailReq{
				Email: "adam@example.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserDetailReq.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
