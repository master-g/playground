package l337c0d3

func lengthOfLongestSubstring(s string) int {
	marks := make(map[byte]int)
	maxLen := 0
	b := []byte(s)
	for i := 0; i < len(b)-1; {
		curLen := 0
		marks[b[i]] = i
		for j := i + 1; j < len(b); j++ {
			if pos, ok := marks[b[j]]; ok {
				if curLen > maxLen {
					// update maxLen if needed
					maxLen = curLen
				}
				// reset map
				marks = make(map[byte]int)
				// move i forward
				i = pos + 1
				break
			} else {
				curLen++
			}
		}
		if curLen > maxLen {
			maxLen = curLen
		}
	}
	return maxLen
}
