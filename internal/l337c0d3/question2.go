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
	var cur *ListNode
	for i := len(str) - 1; i >= 0; i-- {
		val := str[i] - '0'
		if val < 0 || val > 9 {
			err = errors.New("invalid digit")
			return
		}

		n := &ListNode{
			Val:  int(val),
			Next: nil,
		}
		if list == nil {
			list = n
			cur = list
		} else {
			cur.Next = n
			cur = n
		}
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
	n1, n2 := l1, l2
	var result, cur *ListNode
	var sum, carry int
	result = &ListNode{
		Val: 0,
	}
	cur = result
	for {
		if n1 == nil && n2 == nil {
			break
		}

		if n1 != nil && n2 != nil {
			sum = n1.Val + n2.Val + carry
			n1 = n1.Next
			n2 = n2.Next
		} else if n1 != nil {
			sum = n1.Val + carry
			n1 = n1.Next
		} else {
			sum = n2.Val + carry
			n2 = n2.Next
		}

		carry = sum / 10
		n := &ListNode{
			Val: sum % 10,
		}

		cur.Next = n
		cur = cur.Next
	}

	if carry > 0 {
		cur.Next = &ListNode{
			Val: carry,
		}
	}
	return result.Next
}
