package l337c0d3

func lengthOfLongestSubstring(s string) int {
	marks := make(map[byte]int)
	maxLen := 0
	b := []byte(s)
	for i := 0; i < len(b)-1; {
		curLen := 1
		marks[b[i]] = i
		duplicated := false
		for j := i + 1; j < len(b); j++ {
			pos := 0
			if pos, duplicated = marks[b[j]]; duplicated {
				// reset map
				marks = make(map[byte]int)
				// move i forward
				i = pos + 1
				break
			} else {
				marks[b[j]] = j
				curLen++
			}
		}
		if curLen > maxLen {
			maxLen = curLen
		}
		if !duplicated {
			// this round haven't found any duplication
			break
		}
	}
	return maxLen
}
