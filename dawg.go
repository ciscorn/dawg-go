package dawg

import (
	"encoding/binary"
	"fmt"
	"io"
)

// DAWG
type DAWG struct {
	DFA []map[rune]int32 // 決定性有限オートマトン
}

func (d *DAWG) skipPrefix(prefix string, partial bool) (idx int32, count int) {
	for i, c := range prefix {
		node := d.DFA[idx]
		if idxTo, ok := node[c]; ok && idxTo >= 0 {
			idx = idxTo
		} else {
			if partial {
				return idx, i
			} else {
				return -1, i
			}
		}
	}
	return idx, len(prefix)
}

func (d *DAWG) Contains(word string) bool {
	if idx, _ := d.skipPrefix(word, false); idx >= 0 {
		_, ok := d.DFA[idx][-1]
		return ok
	}
	return false
}

func (d *DAWG) ContainsPrefix(word string) bool {
	idx, _ := d.skipPrefix(word, false)
	return idx != -1
}

func (d *DAWG) StartsWith(prefix string, partial bool, maxSize int) []string {
	idx, count := d.skipPrefix(prefix, partial)
	if idx < 0 {
		return nil
	}
	node := d.DFA[idx]

	type stackItem struct {
		node map[rune]int32
		s    string
	}

	results := []string{}
	stack := []stackItem{{node: node, s: ""}}
	for len(stack) > 0 {
		top := stack[0]
		stack = stack[1:]

		if _, ok := top.node[-1]; ok {
			results = append(results, prefix[:count]+top.s)
			if maxSize > 0 && len(results) > maxSize {
				break
			}
		}
		runes := []rune{}
		for r := range top.node {
			if r != -1 {
				runes = append(runes, r)
			}
		}
		idxs := []int32{}
		for _, r := range runes {
			idxs = append(idxs, top.node[r])
		}

		for i := 0; i < len(runes); i++ {
			idx := idxs[i]
			r := runes[i]
			stack = append(stack, stackItem{d.DFA[idx], top.s + string(r)})
		}
	}
	return results
}

func (d *DAWG) Serialize(w io.Writer) {
	// ノード数を書き込む
	numNodes := int32(len(d.DFA))
	binary.Write(w, binary.LittleEndian, numNodes)

	// 各ノードを記録
	for _, node := range d.DFA {
		// エッジ数を記録
		lenEdges := int32(len(node))
		binary.Write(w, binary.LittleEndian, lenEdges)
		// エッジを記録
		for r, idx := range node {
			binary.Write(w, binary.LittleEndian, r)
			binary.Write(w, binary.LittleEndian, idx)
		}
	}
}

func Deserialize(r io.Reader) (*DAWG, error) {
	// ノードを用意
	var numNodes int32
	if err := binary.Read(r, binary.LittleEndian, &numNodes); err != nil {
		return nil, err
	}
	nodes := make([]map[rune]int32, numNodes)
	for i := 0; i < int(numNodes); i++ {
		nodes[i] = make(map[rune]int32)
	}

	// 各ノードを読み込む
	for idxFrom := 0; idxFrom < int(numNodes); idxFrom++ {
		// エッジ数を読み込む
		var lenEdges int32
		if err := binary.Read(r, binary.LittleEndian, &lenEdges); err != nil {
			return nil, err
		}
		// エッジを読み込む
		buf := make([]struct {
			Rune  rune
			IdxTo int32
		}, lenEdges)
		if err := binary.Read(r, binary.LittleEndian, &buf); err != nil {
			return nil, err
		}
		for _, edge := range buf {
			nodes[idxFrom][edge.Rune] = edge.IdxTo
		}
	}

	return &DAWG{
		DFA: nodes,
	}, nil
}

// DumpAsDot は DAWG を Graphviz の DOT ファイルとして出力します
func (d *DAWG) DumpAsDot(w io.Writer) {
	fmt.Fprint(w, "digraph {\n")
	fmt.Fprint(w, "  graph [rankdir = LR];\n")
	for idx, node := range d.DFA {
		if _, ok := node[-1]; ok {
			fmt.Fprintf(w, "  \"%d\" [peripheries = 2];\n", idx)
		}
	}
	for idxFrom, node := range d.DFA {
		for r, idxTo := range node {
			if idxTo >= 0 {
				fmt.Fprintf(w,
					"  \"%d\" -> \"%d\" [label=\"%c\"];\n",
					idxFrom, idxTo, r,
				)
			}
		}
	}
	fmt.Fprint(w, "}\n")
}

// DumpAsMermaid は DAWG を Mermaid の flowchart として出力します
func (d *DAWG) DumpAsMermaid(w io.Writer) {
	fmt.Fprint(w, "flowchart LR\n")
	for idx, node := range d.DFA {
		if _, ok := node[-1]; ok {
			fmt.Fprintf(w, "  N%d(((%d)))\n", idx, idx)
		} else {
			fmt.Fprintf(w, "  N%d((%d))\n", idx, idx)
		}
	}
	for idxFrom, node := range d.DFA {
		for r, idxTo := range node {
			if idxTo >= 0 {
				fmt.Fprintf(w,
					"  N%d -- %c --> N%d\n",
					idxFrom, r, idxTo,
				)
			}
		}
	}
}
