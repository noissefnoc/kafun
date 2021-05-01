package kafun

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	defaultURL, _ := url.Parse(DefaultEndpoint)

	type args struct {
		endpoint string
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "normal case: with default url",
			args: args{
				endpoint: DefaultEndpoint,
			},
			want: &Client{
				URL:        defaultURL,
				HTTPClient: http.DefaultClient,
			},
			wantErr: false,
		},
		{
			name: "normal case: with empty url",
			args: args{
				endpoint: "",
			},
			want: &Client{
				URL:        defaultURL,
				HTTPClient: http.DefaultClient,
			},
			wantErr: false,
		},
		{
			name: "error case: with invalid url",
			args: args{
				endpoint: "hoge",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
