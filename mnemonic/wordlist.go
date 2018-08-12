// https://github.com/bitcoin/bips/blob/master/bip-0039/bip-0039-wordlists.md

package mnemonic

const (
	LANG_EN  = 0
	LANG_CHS = 1
)

var wordValueMap map[int]map[string]int

func init() {
	wordValueMap = make(map[int]map[string]int)
	wordValueMap[LANG_EN] = make(map[string]int)
	wordValueMap[LANG_CHS] = make(map[string]int)
	for i := 0; i < len(wordlist_en); i++ {
		wordValueMap[LANG_EN][wordlist_en[i]] = i
		wordValueMap[LANG_CHS][wordlist_chs[i]] = i
	}
}
