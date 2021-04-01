// Copyright noissefnoc@gmail.com

/*
kafun は環境庁花粉観測システムAPIのクライアントです。

API Provider

環境庁花粉観測システムAPI https://kafun.env.go.jp/apiManual

Synopsys

APIクライアントライブラリとして使う場合、go.mod に追加した後に main.go として以下のようなコードが書ける。

	package main

	import (
		"context"
		"fmt"

		"github.com/noissefnoc/kafun"
	)

	// NOTE: エラーハンドリング省略しているので、実利用のさいは適宜対応してください。
	func main() {
		// APIクライアントを
		client, _ := kafun.NewClient(kafun.DefaultURL)

		// 2021-02から2021-03の東京都の新宿区役所の測定局のデータ取得を指定
		response, _ := client.Search(context.Background, "202102", "202103", "13", "51320100")

		// 最新の時間の測定データを取得
		latestHourlySokuteiData := response[len(response) - 1]

		fmt.Printf(
			"測定日時: %s%s, 気温: %f(度), 花粉量: %d(個/平方メートル)",
			latestHourlySokuteiData.SokuteiNengappi,
			latestHourlySokuteiData.SokuteiJikoku,
			*(latestHourlySokuteiData.AMeDASTempature),
			*(latestHourlySokuteiData.KafunNum),
		)
	}

コマンドラインで使う場合は go get

	go get github.com/noissefnoc/kafun/cmd/kafun

ないしはローカルでビルドした後

	kafun -startYM 202102 -endYM 202103 -todofukenCode 13 -sokuteikyokuCode 51320100

などで実行する。

*/
package kafun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

// DefaultURL は環境庁花粉観測システムAPIのベースURLを表す。
const DefaultURL = "https://kafun.env.go.jp/hanako/api"

// バリデーションのための正規表現。
var yMPattern = regexp.MustCompile(`\d{6}`)
var sokuteikyokuCodePattern = regexp.MustCompile(`\d{8}`)

// APIリクエスト時のUser Agent文字列。
var userAgent = fmt.Sprintf("KafunGoClient/%s (%s)", Version, runtime.Version())

// Client は環境庁花粉観測システムAPIのクライアントを表す。
type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
}

// NewClient は新しいAPIクライアントを作成する。
func NewClient(urlStr string) (*Client, error) {
	if len(urlStr) == 0 {
		urlStr = DefaultURL
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse url: %s: %v", urlStr, err)
	}

	return &Client{
		URL:        parsedURL,
		HTTPClient: http.DefaultClient,
	}, nil
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	spath string,
	params map[string]string,
	body io.Reader,
) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	// 値が空でない場合だけパラメータクエリを作る
	if len(params) != 0 {
		q := u.Query()

		for k, v := range params {
			if len(v) != 0 {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), body)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()

	// APIレスポンスがShift-JISなので、Goで扱えるようにUTF-8変換する。
	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	byteArray, err := ioutil.ReadAll(reader)
	if err != nil {
		xerrors.Errorf("cannot read response: %v", err)
	}

	// APIレスポンスがValidなJSONではないので文字列置換で対応。
	// レスポンスのイメージは以下で、JSONのkey-valueがクォートされていない
	// かつ、valueがない場合何も存在しない (今回は `null` に置換する)
	// [{ KEY11: VALUE11, KEY12: VALUE12 }, { KEY21: VALUE21, KEY22: }]
	str := string(byteArray)
	str = strings.ReplaceAll(str, `{`, `{"`)
	str = strings.ReplaceAll(str, `}`, `"}`)
	str = strings.ReplaceAll(str, `:`, `":"`)
	str = strings.ReplaceAll(str, `,`, `","`)
	str = strings.ReplaceAll(str, `}","{`, `},{`)
	str = strings.ReplaceAll(str, `""`, `null`)

	return json.Unmarshal([]byte(str), out)
}

// Search は 環境庁花粉観測システムAPIの data_search API をコールするメソッド
func (c *Client) Search(
	ctx context.Context,
	startYM string,
	endYM string,
	todofukenCode string,
	sokuteikyokuCode string,
) (SokuteiData, error) {
	params := make(map[string]string)

	if !yMPattern.MatchString(startYM) {
		return nil, xerrors.New(fmt.Sprintf("invalid pattern Start_YM: %s. please specify yyyyMM", startYM))
	}
	params["Start_YM"] = startYM

	if len(endYM) > 0 && !yMPattern.MatchString(endYM) {
		return nil, xerrors.New(fmt.Sprintf("invalid pattern End_YM: %s. please specify yyyyMM", endYM))
	}
	params["End_YM"] = endYM

	todofukenInt, err := strconv.Atoi(todofukenCode)

	if err != nil || 1 > todofukenInt || todofukenInt > 47 {
		return nil, xerrors.Errorf("invalid pattern TDFKN_CD: %s. please specify 01 to 47: %v", endYM, err)
	}
	params["TDFKN_CD"] = todofukenCode

	if len(sokuteikyokuCode) > 0 && !sokuteikyokuCodePattern.MatchString(sokuteikyokuCode) {
		return nil, xerrors.New(
			fmt.Sprintf("invalid pattern SKT_CD: %s. please specify 8 digits", sokuteikyokuCode),
		)
	}
	params["SKT_CD"] = sokuteikyokuCode

	spath := fmt.Sprintf("/data_search")
	req, err := c.newRequest(ctx, http.MethodGet, spath, params, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, xerrors.New(
			fmt.Sprintf(
				"failed to http request with url=%s, status_code=%d",
				req.URL.String(),
				res.StatusCode,
			),
		)
	}

	var response SokuteiData
	if err := decodeBody(res, &response); err != nil {
		return nil, err
	}

	return response, nil
}
