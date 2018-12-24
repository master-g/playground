package l337c0d3

import (
	"fmt"
	"testing"
)

func Test_lengthOfLongestSubstring(t *testing.T) {
	l := lengthOfLongestSubstring("")
	fmt.Println(l)
	l = lengthOfLongestSubstring(" ")
	fmt.Println(l)
	l = lengthOfLongestSubstring("bbbbb")
	fmt.Println(l)
	l = lengthOfLongestSubstring("pwwkew")
	fmt.Println(l)

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lengthOfLongestSubstring(tt.args.s); got != tt.want {
				t.Errorf("lengthOfLongestSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}
