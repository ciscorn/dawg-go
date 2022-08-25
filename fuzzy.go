package dawg

type FuzzySearchResult struct {
	Word  string
	Error float32
}

func FuzzySearch(d *DAWG, word string) []FuzzySearchResult {
	foundWords := make(map[string]float32, 0)
	runes := []rune(word)
	lenRunes := len(runes)

	// FIXME: 誤差に関するパラメータ。決め打ちになっているが、変更可能であるべき。
	k := float32(2.0)  // 最大許容誤差
	ie := float32(1.0) // insertion error
	de := float32(1.0) // deletion error
	se := float32(2.0) // substitution error
	te := float32(2.0) // transposition error

	type Stack struct {
		Error  float32
		Node   map[rune]int32
		Pos    int
		Substr string
	}
	stack := []Stack{{
		Node: d.DFA[0],
	}}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// deletion error
		if curr.Error+de <= k {
			for r, idxTo := range curr.Node {
				if r != -1 {
					stack = append(stack, Stack{
						Error:  curr.Error + de,
						Node:   d.DFA[idxTo],
						Pos:    curr.Pos,
						Substr: curr.Substr + string(r),
					})
				}
			}
		}

		// found
		if curr.Pos == lenRunes {
			if _, ok := curr.Node[-1]; ok {
				if merr, ok := foundWords[curr.Substr]; !ok || merr < curr.Error {
					foundWords[curr.Substr] = curr.Error
				}
			}
			continue
		}

		currRune := runes[curr.Pos]

		// hit
		if idxTo, ok := curr.Node[currRune]; ok && idxTo >= 0 {
			stack = append(stack, Stack{
				Error:  curr.Error,
				Node:   d.DFA[idxTo],
				Pos:    curr.Pos + 1,
				Substr: curr.Substr + string(currRune),
			})
		}

		// insertion error
		if curr.Error+ie <= k {
			stack = append(stack, Stack{
				Error:  curr.Error + ie,
				Node:   curr.Node,
				Pos:    curr.Pos + 1,
				Substr: curr.Substr,
			})
		}

		// substitution error
		if curr.Error+se <= k {
			for r, idxTo := range curr.Node {
				if r != -1 && r != currRune {
					stack = append(stack, Stack{
						Error:  curr.Error + se,
						Node:   d.DFA[idxTo],
						Pos:    curr.Pos + 1,
						Substr: curr.Substr + string(r),
					})
				}
			}
		}

		// transposition error
		if (curr.Error+te <= k) && (curr.Pos < lenRunes-1) {
			runeNext := runes[curr.Pos+1]
			if idxNext, ok := curr.Node[runeNext]; ok && idxNext >= 0 {
				if idxTo, ok := d.DFA[idxNext][currRune]; ok && idxTo >= 0 {
					stack = append(stack, Stack{
						Error:  curr.Error + te,
						Node:   d.DFA[idxTo],
						Pos:    curr.Pos + 1,
						Substr: curr.Substr + string(currRune),
					})
				}
			}
		}
	}

	// 結果を配列に変換
	results := make([]FuzzySearchResult, 0)
	for word, merr := range foundWords {
		results = append(results, FuzzySearchResult{
			Word:  word,
			Error: merr,
		})
	}
	return results
}
