package kafun

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

const testSokuteiDataJSONStringOptional = `[
	{
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
		"KFN_NUM": 0,
		"AMeDAS_WD": "00"
	}
]`

var removePattern = regexp.MustCompile(`\s+|\t|\n`)

func TestCLI_Run(t *testing.T) {
	sjisStr, _ := encodeUTF8ToSJIS(t, sokuteiDataByteOmitOptional)
	type fields struct {
		mockServerHandlerFunc func(w http.ResponseWriter, r *http.Request)
	}
	type args struct {
		args []string
	}
	type want struct {
		returnCode int
		stdout     string
		errout     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "standard case: full flag given",
			fields: fields{
				mockServerHandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write(sjisStr)
				},
			},
			args: args{
				[]string{
					"kafun",
					"-startYM",
					"000000",
					"-endYM",
					"000001",
					"-todofukenCode",
					"01",
					"-sokuteikyokuCode",
					"00000000",
				},
			},
			want: want{
				returnCode: ExitCodeOK,
				stdout:     testSokuteiDataJSONStringOptional,
				errout:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdOut := new(bytes.Buffer)
			errOut := new(bytes.Buffer)
			c := &CLI{
				OutStream: stdOut,
				ErrStream: errOut,
			}
			testServer := httptest.NewServer(http.HandlerFunc(tt.fields.mockServerHandlerFunc))
			defer testServer.Close()
			DefaultEndpoint = testServer.URL

			if got := c.Run(tt.args.args); got != tt.want.returnCode {
				t.Errorf("Run() return code = %v, want %v", got, tt.want.returnCode)
			}
			if stdOut.String() != tt.want.stdout {
				t.Errorf("Run() stdout = %v, want %v", stdOut.String(), tt.want.stdout)
			}
			if errOut.String() != tt.want.errout {
				t.Errorf("Run() errout = %v, want %v", errOut.String(), tt.want.stdout)
			}
		})
	}
}
