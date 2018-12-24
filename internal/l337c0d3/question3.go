package l337c0d3

func lengthOfLongestSubstring(s string) int {
	// char and its position
	marks := make(map[byte]int)
	maxLen := 0
	n := len(s)

	for i, j := 0, 0; j < n; j++ {
		if pos, duplicated := marks[s[j]]; duplicated && pos > i {
			i = pos
		}
		if j-i+1 > maxLen {
			maxLen = j - i + 1
		}
		marks[s[j]] = j + 1
	}
	return maxLen
}
