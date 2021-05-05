package kafun

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func decodeBodyResponseFixture(t *testing.T, responseBody []byte) *http.Response {
	t.Helper()

	encoder := japanese.ShiftJIS.NewEncoder()
	sjisStr, sjisLen, _ := transform.Bytes(encoder, responseBody)

	return &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewReader(sjisStr)),
		ContentLength: int64(sjisLen),
	}
}

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

func Test_decodeBody(t *testing.T) {
	type args struct {
		resp *http.Response
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "standard case",
			args: args{
				resp: decodeBodyResponseFixture(
					t,
					[]byte(`[{
						"SKT_CD": "00000000",
						"AMeDAS_CD": "00000",
						"SKT_NNGP": "00000000",
						"SKT_HH": "00",
						"SKT_NM": "テスト観測所",
						"SKT_TYPE": "0",
						"TDFKN_CD": "00",
						"TDFKN_NM": "テスト都道府県",
						"SKCHSN_CD": "000000",
						"SKCHSN_NM": "テスト市町村",
						"KFN_NUM": "0",
						"AMeDAS_WD": "00",
					}]`),
				),
				out: SokuteiData{},
			},
			want: SokuteiData{
				&HourlySokuteiData{
					SokuteikyokuCode:     "00000000",
					AMeDASCode:           "00000",
					SokuteiNengappi:      "00000000",
					SokuteiJikoku:        "00",
					SokuteikyokuName:     "テスト観測所",
					SokuteiType:          "0",
					TodofukenCode:        "00",
					TodofukenName:        "テスト都道府県",
					SokuteiShichosonCode: "000000",
					SokuteiShichosonName: "テスト市町村",
					KafunNum:             0,
					AMeDASWindDirect:     "00",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := decodeBody(tt.args.resp, tt.args.out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("decodeBody() got = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}
