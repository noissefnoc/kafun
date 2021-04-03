package kafun

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	defaultURL, _ := url.Parse(DefaultURL)

	type args struct {
		urlStr string
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
				urlStr: DefaultURL,
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
				urlStr: "",
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
				urlStr: "hoge",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.urlStr)
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
