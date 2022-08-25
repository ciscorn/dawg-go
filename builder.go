package dawg

import (
	"sort"
	"strconv"
	"strings"
)

// Builder
type Builder struct {
	nfa       []map[rune]*[]int // 非決定性有限オートマトン
	suffixMap map[string]int
}

func NewBuilder() *Builder {
	nodes := make([]map[rune]*[]int, 0)
	nodes = append(nodes, make(map[rune]*[]int))
	return &Builder{
		nfa:       nodes,
		suffixMap: make(map[string]int),
	}
}

func (b *Builder) AddWord(word string) {
	if len(word) == 0 {
		return
	}
	var ok bool
	var r rune
	var nodeIdx = 0
	var suffixIdx = 0
	var i int = -1
	runes := []rune(word)
	for i, r = range runes {
		suffix := string(runes[i+1:])
		nodeCurr := b.nfa[nodeIdx]
		if suffixIdx, ok = b.suffixMap[suffix]; ok {
			if edges, ok := nodeCurr[r]; ok {
				*edges = append(*edges, suffixIdx)
			} else {
				nodeCurr[r] = &[]int{suffixIdx}
			}
			break
		} else {
			suffixIdx = len(b.nfa)
			b.nfa = append(b.nfa, make(map[rune]*[]int))
			b.suffixMap[suffix] = suffixIdx
			if edges, ok := nodeCurr[r]; ok {
				*edges = append(*edges, suffixIdx)
			} else {
				nodeCurr[r] = &[]int{suffixIdx}
			}
			nodeIdx = suffixIdx
		}
	}
	if i == len(runes)-1 {
		b.nfa[suffixIdx][-1] = &[]int{-1}
	}
}

func (b *Builder) Build() *DAWG {
	// NFA をもとに DFA (DAWG) を構築する
	nodes := make([]map[rune]int32, 0)
	nodes = append(nodes, make(map[rune]int32))
	dawg := &DAWG{
		DFA: nodes,
	}
	closureMap := make(map[string]int32)
	closureMap["0"] = 0

	stack := []map[int]struct{}{
		{0: struct{}{}},
	}
	for len(stack) > 0 {
		fromClosure := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		fromIdx := closureMap[convertIntSetToStringKey(fromClosure)]

		res := make(map[rune]map[int]struct{})
		for nfaIdx := range fromClosure {
			for r, edges := range b.nfa[nfaIdx] {
				if cl, ok := res[r]; ok {
					for _, e := range *edges {
						cl[e] = struct{}{}
					}
				} else {
					res[r] = make(map[int]struct{})
					for _, e := range *edges {
						res[r][e] = struct{}{}
					}
				}
			}
		}
		for r, cl := range res {
			toClosure := convertIntSetToStringKey(cl)
			if r == -1 {
				dawg.DFA[fromIdx][-1] = -1
				continue
			}
			if toIdx, ok := closureMap[toClosure]; ok {
				dawg.DFA[fromIdx][r] = toIdx
			} else {
				toIdx = int32(len(dawg.DFA))
				dawg.DFA = append(dawg.DFA, make(map[rune]int32))
				closureMap[toClosure] = toIdx
				stack = append(stack, cl)
				dawg.DFA[fromIdx][r] = toIdx
			}
		}
	}
	return dawg
}

// int の集合を "5,11,130" のような文字列に変換する
func convertIntSetToStringKey(s map[int]struct{}) string {
	ss := []int{}
	for v := range s {
		ss = append(ss, v)
	}
	sort.Ints(ss)
	sss := make([]string, 0, len(ss))
	for _, v := range ss {
		sss = append(sss, strconv.Itoa(v))
	}
	return strings.Join(sss, ",")
}
