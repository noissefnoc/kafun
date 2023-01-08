# kafun

[![GitHub Actions](https://github.com/noissefnoc/kafun/workflows/CI/badge.svg)](https://github.com/noissefnoc/kafun/actions?workflow=CI)

*DEPRECATED: [環境省花粉観測システム（愛称：はなこさん）事業の廃止に伴う花粉自動計測器を用いた花粉観測の終了について - 2021年12月24日](https://www.env.go.jp/press/110339.html)にあるように、2022年以降該当のAPIサービスが停止されたため、このコマンドは利用できません。*

Golang製の [環境庁花粉観測システムAPI](https://kafun.env.go.jp/apiManual) のクライアントです。

**アルファ版なので様々な変更の可能性があります**

## インストール

### コマンド

現在(2021-04) `go get` ないしはローカル環境でのビルドのみの対応です。

```shell
go get github.com/noissefnoc/kafun/cmd/kafun
```

ないしは

```shell
git clone git@github.com:noissefnoc/kafun.git
cd kafun
go build -o kafun cmd/kafun/main.go
```

### ライブラリ

`go.mod` に `github.com/noissefnoc/kafun` を追加してください。

## 使い方

### コマンド

#### オプション

```
Usage of kafun:
  -endYM string
        終了年月 (format: yyyyMM)
  -sokuteikyokuCode string
        測定局コード
  -startYM string
        開始年月 (format: yyyyMM) (必須)
  -todofukenCode string
        都道府県コード (range: 01 to 47) (必須)
```

#### 具体用例

* 取得期間：2021-02〜2021-03
* 対象都道府県：東京都(13)
* 対象観測所：新宿区役所第二分庁舎 (51320100)
    * https://kafun.env.go.jp/preview/table/51320100/today のURLから判断

```shell
kafun -startYM 202102 -endYM 202103 -todofukenCode 13 -sokuteikyokuCode 51320100
```

### ライブラリ

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/noissefnoc/kafun"
)

// エラー処理省略
func main() {
	client, _ := kafun.NewClient(kafun.DefaultURL)
	response, _ := client.Search(context.Background(), "202102", "202103", "13", "51320100")
	printJSON, _ := json.MarshalIndent(response[len(response) - 1], "", "\t") // 最新の測定データのみを表示対象にする
	fmt.Println(printJSON)
}
```


## 環境庁花粉測定システムAPI公式サイト

* [APIの説明ページ](https://kafun.env.go.jp/apiManual): APIトップページ
* [1-1. 環境省花粉観測システムで提供しているAPIに関する情報](https://kafun.env.go.jp/apiManual/apiPage1/api-1-1): APIデータの提供期間について。2月1日〜6月30日
* [2-2.提供データについて](https://kafun.env.go.jp/apiManual/apiPage2/api-2-2): APIレスポンスの各項目説明
* [3-1-1.パラメータ一覧](https://kafun.env.go.jp/apiManual/apiPage3/api-3-1-1): APIのリクエストパラメータ(`kafun` コマンドのオプション)


## Author

noissefnoc@gmail.com
