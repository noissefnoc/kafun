package kafun

import (
	"encoding/json"
	"strconv"
)

// HourlySokuteiData は時間毎の測定データを表します。
// 数値型のもので、ゼロではなく、空文字列のものはnullとみなしてJSON出力時には項目を出力しない仕様にしています。
// 詳細は https://kafun.env.go.jp/apiManual/apiPage2/api-2-3 確認してください
type HourlySokuteiData struct {
	// 測定局コード
	SokuteikyokuCode string `json:"SKT_CD"`

	// アメダスコード
	AMeDASCode string `json:"AMeDAS_CD"`

	// 測定年月日(yyyyMMdd)
	SokuteiNengappi string `json:"SKT_NNGP"`

	// 測定時刻(HH)
	SokuteiJikoku string `json:"SKT_HH"`

	// 測定局名
	SokuteikyokuName string `json:"SKT_NM"`

	// 測定局のタイプ
	SokuteiType string `json:"SKT_TYPE"`

	// 都道府県コード(JIS)
	TodofukenCode string `json:"TDFKN_CD"`

	// 都道府県名
	TodofukenName string `json:"TDFKN_NM"`

	// 測定局の市区町村コード
	SokuteiShichosonCode string `json:"SKCHSN_CD"`

	// 測定局の市区町村名
	SokuteiShichosonName string `json:"SKCHSN_NM"`

	// 花粉数(個/立方メートル)
	KafunNum int `json:"KFN_NUM"`

	// 風向き
	AMeDASWindDirect string `json:"AMeDAS_WD"`

	// 風速(m/s)
	AMeDASWindSpeed *int `json:"AMeDAS_WS,omitempty"`

	// 気温(度)
	AMeDASTemperature *float64 `json:"AMeDAS_TP,omitempty"`

	// 降水量(mm)
	AMeDASPrecipitation *int `json:"AMeDAS_PR,omitempty"`

	// レーダー降雨降雪の有無
	AMeDASRadarPrecipitation *int `json:"AMeDAS_RDPR,omitempty"`
}

// SokuteiData はData Search APIのレスポンスを表します。
type SokuteiData []*HourlySokuteiData

// UnmarshalJSON は HourlySokuteiData が Valid な JSON ではないために作成したカスタムUnmarshaler
//
// - value が 数値型の場合でもクォートされる。stringタグを使うと出力のさいにクォートついてしまうので対応
// - 数値で空文字列が入ってくる。ゼロはあるので、こちらはnullとみなして key-valueを生成しない
func (hsd *HourlySokuteiData) UnmarshalJSON(data []byte) error {
	var err error
	var v map[string]interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	hsd.SokuteikyokuCode = v["SKT_CD"].(string)
	hsd.AMeDASCode = v["AMeDAS_CD"].(string)
	hsd.SokuteiNengappi = v["SKT_NNGP"].(string)
	hsd.SokuteiJikoku = v["SKT_HH"].(string)
	hsd.SokuteikyokuName = v["SKT_NM"].(string)
	hsd.SokuteiType = v["SKT_TYPE"].(string)
	hsd.TodofukenCode = v["TDFKN_CD"].(string)
	hsd.TodofukenName = v["TDFKN_NM"].(string)
	hsd.SokuteiShichosonCode = v["SKCHSN_CD"].(string)
	hsd.SokuteiShichosonName = v["SKCHSN_NM"].(string)

	kafunNumStr := v["KFN_NUM"].(string)
	kafunNumInt, err := strconv.Atoi(kafunNumStr)
	if err != nil {
		return err
	}
	hsd.KafunNum = kafunNumInt

	hsd.AMeDASWindDirect = v["AMeDAS_WD"].(string)

	intPointerElem, err := validateIntPointerElement(v, "AMeDAS_WS")
	if err != nil {
		return err
	}
	hsd.AMeDASWindSpeed = intPointerElem

	float64PointerElem, err := validateFloat64PointerElement(v, "AMeDAS_TP")
	if err != nil {
		return err
	}
	hsd.AMeDASTemperature = float64PointerElem

	intPointerElem, err = validateIntPointerElement(v, "AMeDAS_PR")
	if err != nil {
		return err
	}
	hsd.AMeDASPrecipitation = intPointerElem

	intPointerElem, err = validateIntPointerElement(v, "AMeDAS_RDPR")
	if err != nil {
		return err
	}
	hsd.AMeDASRadarPrecipitation = intPointerElem

	return nil
}

func validateIntPointerElement(v map[string]interface{}, key string) (*int, error) {
	elem, ok := v[key]
	if ok {
		elemStr := elem.(string)
		if len(elemStr) != 0 {
			elemInt, err := strconv.Atoi(elemStr)
			if err != nil {
				return nil, err
			}

			return &elemInt, nil
		}
	}

	return nil, nil
}

func validateFloat64PointerElement(v map[string]interface{}, key string) (*float64, error) {
	elem, ok := v[key]
	if ok {
		elemStr := elem.(string)
		if len(elemStr) != 0 {
			elemFloat64, err := strconv.ParseFloat(elemStr, 64)
			if err != nil {
				return nil, err
			}

			return &elemFloat64, nil
		}
	}

	return nil, nil
}
