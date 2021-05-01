package kafun

import (
	"reflect"
	"testing"
)

func intPointerHelper(t *testing.T, i int) *int {
	t.Helper()
	return &i
}

func float64PointerHelper(t *testing.T, f float64) *float64 {
	t.Helper()
	return &f
}

func TestHourlySokuteiData_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *HourlySokuteiData
		wantErr bool
	}{
		{
			name: "standard case: fill all fields",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05",
					"AMeDAS_WS": "1",
					"AMeDAS_TP": "15.4",
					"AMeDAS_PR": "0",
					"AMeDAS_RDPR": "0"
        		}`),
			},
			want: &HourlySokuteiData{
				SokuteikyokuCode:         "00000000",
				AMeDASCode:               "00000",
				SokuteiNengappi:          "00000101",
				SokuteiJikoku:            "01",
				SokuteikyokuName:         "テスト測定所",
				SokuteiType:              "1",
				TodofukenCode:            "00",
				TodofukenName:            "テスト県",
				SokuteiShichosonCode:     "00000",
				SokuteiShichosonName:     "テスト市",
				KafunNum:                 4,
				AMeDASWindDirect:         "05",
				AMeDASWindSpeed:          intPointerHelper(t, 1),
				AMeDASTemperature:        float64PointerHelper(t, 15.4),
				AMeDASPrecipitation:      intPointerHelper(t, 0),
				AMeDASRadarPrecipitation: intPointerHelper(t, 0),
			},
		},
		{
			name: "standard case: omit all optional fields",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05"
        		}`),
			},
			want: &HourlySokuteiData{
				SokuteikyokuCode:     "00000000",
				AMeDASCode:           "00000",
				SokuteiNengappi:      "00000101",
				SokuteiJikoku:        "01",
				SokuteikyokuName:     "テスト測定所",
				SokuteiType:          "1",
				TodofukenCode:        "00",
				TodofukenName:        "テスト県",
				SokuteiShichosonCode: "00000",
				SokuteiShichosonName: "テスト市",
				KafunNum:             4,
				AMeDASWindDirect:     "05",
			},
		},
		{
			name: "error case: KFN_NUM format is not int",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "invalid",
					"AMeDAS_WD": "05"
        		}`),
			},
			wantErr: true,
		},
		{
			name: "error case: AMeDAS_WS format is not int",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05",
					"AMeDAS_WS": "invalid"
        		}`),
			},
			wantErr: true,
		},
		{
			name: "error case: AMeDAS_WS format is not float64",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05",
					"AMeDAS_WS": "1",
					"AMeDAS_TP": "invalid"
        		}`),
			},
			wantErr: true,
		},
		{
			name: "error case: AMeDAS_PR format is not int",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05",
					"AMeDAS_WS": "1",
					"AMeDAS_TP": "15.4",
					"AMeDAS_PR": "invalid"
        		}`),
			},
			wantErr: true,
		},
		{
			name: "error case: AMeDAS_RDPR format is not int",
			args: args{
				[]byte(`{
					"SKT_CD": "00000000",
					"AMeDAS_CD": "00000",
					"SKT_NNGP": "00000101",
					"SKT_HH": "01",
					"SKT_NM": "テスト測定所",
					"SKT_TYPE": "1",
					"TDFKN_CD": "00",
					"TDFKN_NM": "テスト県",
					"SKCHSN_CD": "00000",
					"SKCHSN_NM": "テスト市",
					"KFN_NUM": "4",
					"AMeDAS_WD": "05",
					"AMeDAS_WS": "1",
					"AMeDAS_TP": "15.4",
					"AMeDAS_PR": "0",
					"AMeDAS_RDPR": "invalid"
        		}`),
			},
			wantErr: true,
		},
		{
			name: "error case: response format is not JSON",
			args: args{
				[]byte(`invalid`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &HourlySokuteiData{}
			err := got.UnmarshalJSON(tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_validateIntPointerElement(t *testing.T) {
	type args struct {
		v   map[string]interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *int
		wantErr bool
	}{
		{
			name: "standard case: key exists and int value",
			args: args{
				v: map[string]interface{}{
					"key1": "1",
					"key2": "2",
				},
				key: "key1",
			},
			want: intPointerHelper(t, 1),
		},
		{
			name: "standard case: key exists and empty string value",
			args: args{
				v: map[string]interface{}{
					"key1": "",
					"key2": "2",
				},
				key: "key1",
			},
			want: nil,
		},
		{
			name: "standard case: key doesn't exist",
			args: args{
				v: map[string]interface{}{
					"key2": "2",
				},
				key: "key1",
			},
			want: nil,
		},
		{
			name: "error case: key exists and invalid string",
			args: args{
				v: map[string]interface{}{
					"key1": "invalid",
					"key2": "2",
				},
				key: "key1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateIntPointerElement(tt.args.v, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIntPointerElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateIntPointerElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateFloat64PointerElement(t *testing.T) {
	type args struct {
		v   map[string]interface{}
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *float64
		wantErr bool
	}{
		{
			name: "standard case: key exists and int value",
			args: args{
				v: map[string]interface{}{
					"key1": "1.0",
					"key2": "2",
				},
				key: "key1",
			},
			want: float64PointerHelper(t, 1),
		},
		{
			name: "standard case: key exists and empty string value",
			args: args{
				v: map[string]interface{}{
					"key1": "",
					"key2": "2",
				},
				key: "key1",
			},
			want: nil,
		},
		{
			name: "standard case: key doesn't exist",
			args: args{
				v: map[string]interface{}{
					"key2": "2",
				},
				key: "key1",
			},
			want: nil,
		},
		{
			name: "error case: key exists and invalid string",
			args: args{
				v: map[string]interface{}{
					"key1": "invalid",
					"key2": "2",
				},
				key: "key1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateFloat64PointerElement(tt.args.v, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFloat64PointerElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateFloat64PointerElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}
