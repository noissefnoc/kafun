package kafun

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
)

// 終了コードの状態。
const (
	ExitCodeOK              = iota // 正常終了
	ExitCodeParseFlagError         // コマンドラインフラグのパースエラー終了
	ExitCodeInitializeError        // 初期化のエラー終了
	ExitCodeAPIRequestError        // APIリクエストエラー終了
)

// コマンドラインフラグ
var (
	startYM          string // 取得開始の年月を指定するフラグ
	endYM            string // 取得終了の年月を指定するフラグ
	todofukenCode    string // 都道府県コードを指定するフラグ
	sokuteikyokuCode string // 測定局コードを指定するフラグ
)

// CLI はコマンドを作成するさいの入出力を表す。
type CLI struct {
	OutStream io.Writer // 出力のストリーム
	ErrStream io.Writer // エラー出力のストリーム
}

// Run はコマンドを実行する関数
func (c *CLI) Run(args []string) int {
	// コマンドライン引数をパース
	flags := flag.NewFlagSet("kafun", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)
	flags.StringVar(
		&startYM,
		"startYM",
		"",
		"開始年月 (format: yyyyMM) (必須)",
	)
	flags.StringVar(
		&endYM,
		"endYM",
		"",
		"終了年月 (format: yyyyMM)",
	)
	flags.StringVar(
		&todofukenCode,
		"todofukenCode",
		"",
		"都道府県コード (range: 01 to 47) (必須)",
	)
	flags.StringVar(
		&sokuteikyokuCode,
		"sokuteikyokuCode",
		"",
		"測定局コード。複数指定の場合はカンマ区切りで指定",
	)

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	// コマンドライン引数が指定されていない場合は用例を表示
	if len(args) == 1 {
		flags.Usage()
		return ExitCodeOK
	}

	// data_search API の実行
	client, err := NewClient(DefaultEndpoint)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialize API client with url=%s: %v\n", DefaultEndpoint, err)
		return ExitCodeInitializeError
	}

	param := &SearchParam{
		StartYM:          startYM,
		EndYM:            endYM,
		TodofukenCode:    todofukenCode,
		SokuteikyokuCode: sokuteikyokuCode,
	}

	response, err := client.Search(context.Background(), param)
	if err != nil {
		fmt.Fprintf(
			c.ErrStream,
			"failed to request to API with args startYM=%s, endYM=%s, todofukenCode=%s, sokuteikyokuCode=%s: %v\n",
			startYM,
			endYM,
			todofukenCode,
			sokuteikyokuCode,
			err,
		)
		return ExitCodeAPIRequestError
	}

	printJSON, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialize API client: %v\n", err)
	}

	fmt.Fprintf(c.OutStream, "%s", printJSON)

	return ExitCodeOK
}
