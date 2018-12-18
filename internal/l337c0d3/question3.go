package l337c0d3

func lengthOfLongestSubstring(s string) int {
	set := make(map[byte]bool)
	max := 0
	b := []byte(s)
	for i := 0; i < len(b)-1; {
		curLen := 0
		set[b[i]] = true
		for j := i + 1; j < len(b); j++ {
			if set[b[j]] {
				if curLen > max {
					// update max if needed
					max = curLen
				}
				// reset map
				set = make(map[byte]bool)
				// TODO
				break
			}
		}
	}
	return max
}
