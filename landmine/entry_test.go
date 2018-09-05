package landmine

import "testing"

func TestGetACat(t *testing.T) {
	if GetACat() == nil {
		t.Errorf("Forgot to return a real cat!")
	}
}
