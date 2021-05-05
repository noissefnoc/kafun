package kafun

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var sokuteiDataByteOmitOptional = []byte(`[{
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
	"AMeDAS_WD": "00"
}]`)

func encodeUTF8ToSJIS(t *testing.T, b []byte) ([]byte, int) {
	t.Helper()

	encoder := japanese.ShiftJIS.NewEncoder()
	sjisStr, sjisLen, _ := transform.Bytes(encoder, b)

	return sjisStr, sjisLen
}

func decodeBodyResponseFixture(t *testing.T, responseBody []byte) *http.Response {
	t.Helper()

	sjisStr, sjisLen := encodeUTF8ToSJIS(t, responseBody)

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

func TestClient_newRequest(t *testing.T) {
	ctx := context.Background()
	endpointURL, _ := url.Parse(DefaultEndpoint)
	fullQueryURL, _ := url.Parse(
		DefaultEndpoint + "/search?End_YM=000001&SKT_CD=00000000&Start_YM=000000&TDFKN_CD=01",
	)
	omitQueryURL, _ := url.Parse(
		DefaultEndpoint + "/search?Start_YM=000000&TDFKN_CD=01",
	)

	type args struct {
		method  string
		apiPath string
		params  map[string]string
		body    io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "standard case: full query given",
			args: args{
				method:  "GET",
				apiPath: "/search",
				params: map[string]string{
					"Start_YM": "000000",
					"End_YM":   "000001",
					"TDFKN_CD": "01",
					"SKT_CD":   "00000000",
				},
			},
			want: (&http.Request{
				Method:     "GET",
				URL:        fullQueryURL,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: map[string][]string{
					"User-Agent": {
						userAgent,
					},
				},
				Host: "kafun.env.go.jp",
			}).WithContext(ctx),
		},
		{
			name: "standard case: omit optional query",
			args: args{
				method:  "GET",
				apiPath: "/search",
				params: map[string]string{
					"Start_YM": "000000",
					"TDFKN_CD": "01",
				},
			},
			want: (&http.Request{
				Method:     "GET",
				URL:        omitQueryURL,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: map[string][]string{
					"User-Agent": {
						userAgent,
					},
				},
				Host: "kafun.env.go.jp",
			}).WithContext(ctx),
		},
		{
			name: "error case: invalid method request",
			args: args{
				method: "無効",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				URL:        endpointURL,
				HTTPClient: http.DefaultClient,
			}
			got, err := c.newRequest(ctx, tt.args.method, tt.args.apiPath, tt.args.params, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Fatalf("newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRequest() got = %v, want %v", got, tt.want)
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

func TestClient_Search(t *testing.T) {
	ctx := context.Background()
	sjisStr, _ := encodeUTF8ToSJIS(t, sokuteiDataByteOmitOptional)

	type fields struct {
		mockServerHandlerFunc func(w http.ResponseWriter, r *http.Request)
	}
	type args struct {
		param *SearchParam
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    SokuteiData
		wantErr bool
	}{
		{
			name: "standard case",
			fields: fields{
				mockServerHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write(sjisStr)
				},
			},
			args: args{
				param: &SearchParam{
					StartYM:          "000000",
					EndYM:            "000001",
					TodofukenCode:    "01",
					SokuteikyokuCode: "00000000",
				},
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
		},
		{
			name: "error case: invalid parameter",
			fields: fields{
				mockServerHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write(sjisStr)
				},
			},
			args: args{
				param: &SearchParam{},
			},
			wantErr: true,
		},
		{
			name: "error case: server returns bad request",
			fields: fields{
				mockServerHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
				},
			},
			args: args{
				param: &SearchParam{
					StartYM:          "000000",
					EndYM:            "000001",
					TodofukenCode:    "01",
					SokuteikyokuCode: "00000000",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(tt.fields.mockServerHandlerFunc))
			defer testServer.Close()
			testServerURL, _ := url.Parse(testServer.URL)
			c := &Client{
				URL:        testServerURL,
				HTTPClient: http.DefaultClient,
			}
			got, err := c.Search(ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
