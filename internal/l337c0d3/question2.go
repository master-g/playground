package l337c0d3

import (
	"errors"
	"fmt"
	"strings"
)

// ListNode struct define
type ListNode struct {
	Val  int
	Next *ListNode
}

func newListNodeFromString(str string) (list *ListNode, err error) {
	var n *ListNode
	for i := len(str) - 1; i >= 0; i-- {
		val := str[i] - '0'
		if val < 0 || val > 9 {
			err = errors.New("invalid digit")
			return
		}

		n = &ListNode{
			Val:  int(val),
			Next: nil,
		}
		if list == nil {
			list = n
		}
		n = n.Next
	}
	return
}

func (l *ListNode) String() string {
	if l == nil {
		return "nil"
	}

	sb := strings.Builder{}
	for n := l; n != nil; n = n.Next {
		sb.WriteString(fmt.Sprintf("%v", n.Val))
	}
	return sb.String()
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	n1 := l1
	n2 := l2
	var result *ListNode
	r := result
	overflow := 0
	sum := 0
	for {
		if n1 == nil && n2 == nil {
			// no more data
			if overflow > 0 {
				r = &ListNode{
					Val:  overflow,
					Next: nil,
				}
			}
			break
		}

		if n1 != nil && n2 != nil {
			sum = n1.Val + n2.Val + overflow
			n1 = n1.Next
			n2 = n2.Next
		} else if n1 != nil {
			sum = n1.Val + overflow
			n1 = n1.Next
		} else {
			sum = n2.Val + overflow
			n2 = n2.Next
		}

		if sum >= 10 {
			overflow = 1
		} else {
			overflow = 0
		}
		sum = sum % 10

		r = &ListNode{
			Val:  sum,
			Next: nil,
		}

		r = r.Next
	}
	return result
}
