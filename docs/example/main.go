package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/ciscorn/dawg-go"
)

func main() {
	// DAWG を構築
	builder := dawg.NewBuilder()
	builder.AddWord("仙台市")
	builder.AddWord("仙台市青葉区")
	builder.AddWord("仙台市役所")
	builder.AddWord("青葉区")
	builder.AddWord("横浜市")
	builder.AddWord("横浜市役所")
	builder.AddWord("横浜市立大学")
	builder.AddWord("横浜国立大学")
	builder.AddWord("横浜市青葉区")
	wg := builder.Build()

	// DAWGが指定の語を含むかどうか
	if wg.Contains("横浜市") {
		fmt.Println("横浜市")
	}

	// DAWGを辞書として、文字列から語を抽出する
	for _, r := range dawg.ExtractKeywords(wg, "青葉区といえば仙台市青葉区と横浜市青葉区があります") {
		fmt.Println(r)
	}

	// DAWGに含まれる語をあいまい検索する
	for _, hit := range dawg.FuzzySearch(wg, "横浜私立大学") {
		fmt.Println(hit)
	}

	// DAWG を Graphviz の dot ファイルとして出力
	f, err := os.OpenFile("docs/example/example.dot", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	wg.DumpAsDot(w)
	w.Flush()

	// DAWG を Mermaid の flowchart として出力
	buf := bytes.NewBuffer(nil)
	wg.DumpAsMermaid(buf)
	fmt.Println(buf.String())
}
