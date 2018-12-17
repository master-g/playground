package l337c0d3

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_newListNodeFromString(t *testing.T) {
	l, err := newListNodeFromString("1234")
	if err != nil {
		t.Error(err)
	}
	log.Info(l.String())
}

func Test_addTwoNumbers(t *testing.T) {
	type args struct {
		l1 *ListNode
		l2 *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addTwoNumbers(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addTwoNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
