package dawg_test

import (
	"fmt"
	"testing"

	"github.com/ciscorn/dawg-go"
)

func Test_Foobar(t *testing.T) {
	builder := dawg.NewBuilder()
	builder.AddWord("top")
	builder.AddWord("tops")
	builder.AddWord("tap")
	builder.AddWord("taps")

	d := builder.Build()
	// dawg.DumpAsDot(os.Stdout)
	// dawg.DumpAsMermaid(os.Stdout)

	// FIXME
	fmt.Println(d.Contains("top"))
	fmt.Println(d.Contains("tops"))
	fmt.Println(d.Contains("tap"))
	fmt.Println(d.Contains("taps"))

	fmt.Println(d.Contains("zaps"))
	fmt.Println(d.Contains("ta"))
	fmt.Println(d.Contains("to"))
	fmt.Println(d.Contains("tapsy"))

	fmt.Println(d.ContainsPrefix("t"))
	fmt.Println(d.ContainsPrefix("to"))
	fmt.Println(d.ContainsPrefix("top"))
	fmt.Println(d.ContainsPrefix("tops"))
	fmt.Println(d.ContainsPrefix("topsi"))
	fmt.Println(d.ContainsPrefix("tapsy"))

	for r := range dawg.FuzzySearch(d, "topz") {
		fmt.Printf("%+v\n", r)
	}

	fmt.Println(dawg.ExtractKeywords(d, "this is tops top hits tap taps feee"))
}
