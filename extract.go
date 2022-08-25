package dawg

type ExtractKeywordsResult struct {
	Found    bool
	Fragment string
}

func ExtractKeywords(d *DAWG, document string) (result []ExtractKeywordsResult) {
	type BufItem struct {
		Rune    rune
		Node    map[rune]int32
		Advance int
	}

	buf := []*BufItem{}
	nonhitRunes := []rune{}

	for _, r := range document + string(rune(-1)) {
		if len(buf) == 0 {
			if _, ok := d.DFA[0][r]; !ok {
				nonhitRunes = append(nonhitRunes, r)
				continue
			}
		}
		buf = append(buf, &BufItem{
			Rune:    r,
			Node:    d.DFA[0],
			Advance: 0,
		})

		for i, bi := range buf {
			if bi.Node == nil {
				continue
			}
			if nn, ok := bi.Node[r]; !ok || nn == -1 {
				bi.Node = nil
			} else {
				bi.Node = d.DFA[nn]
				if _, ok := bi.Node[-1]; ok {
					bi.Advance = len(buf) - i
				}
			}
		}

		var i = 0
		for i < len(buf) && buf[i].Node == nil {
			bi := buf[i]
			lenk := bi.Advance
			if lenk > 0 {
				runes := make([]rune, lenk)
				for i, r := range buf[i : i+lenk] {
					runes[i] = r.Rune
				}
				if len(nonhitRunes) > 0 {
					result = append(result, ExtractKeywordsResult{
						Fragment: string(nonhitRunes),
					})
				}
				nonhitRunes = nil
				result = append(result, ExtractKeywordsResult{
					Found:    true,
					Fragment: string(runes),
				})
				i += lenk
			} else {
				nonhitRunes = append(nonhitRunes, bi.Rune)
				i += 1
			}
		}
		if i > 0 {
			buf = buf[i:]
		}
	}

	if len(nonhitRunes) > 1 {
		result = append(result, ExtractKeywordsResult{
			Fragment: string(nonhitRunes[:len(nonhitRunes)-1]),
		})
	}

	return result
}
