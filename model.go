package kafun

// HourlySokuteiData は時間毎の測定データを表します。元データの都合上数値でもクォートして取り扱っています。
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
	KafunNum int `json:"KFN_NUM,string"`

	// 風向き
	AMeDASWindDirect string `json:"AMeDAS_WD"`

	// 風速(m/s)
	AMeDASWindSpeed *int `json:"AMeDAS_WS,string"`

	// 気温(度)
	AMeDASTemperature *float64 `json:"AMeDAS_TP,string"`

	// 降水量(mm)
	AMeDASPrecipitation *int `json:"AMeDAS_PR,string"`

	// レーダー降雨降雪の有無
	AMeDASRadarPrecipitation *int `json:"AMeDAS_RDPR,string"`
}

// SokuteiData はData Search APIのレスポンスを表します。
type SokuteiData []*HourlySokuteiData
