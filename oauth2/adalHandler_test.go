package oauth2

import (
	"net/url"
	"testing"
)

func TestAdalHandler_UpdateFromRedirect(t *testing.T) {
	type args struct {
		redirect *url.URL
	}
	tests := []struct {
		name    string
		h       AdalHandler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.UpdateFromRedirect(tt.args.redirect); (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.UpdateFromRedirect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdalHandler_GetAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		h       AdalHandler
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.h.GetAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AdalHandler.GetAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdalHandler_UpdateFromCode(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		h       AdalHandler
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.UpdateFromCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("AdalHandler.UpdateFromCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
