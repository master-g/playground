package l337c0d3

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	n1 := l1
	n2 := l2
	var result *ListNode
	r := result
	for {
		if n1 == nil || n2 == nil {
			break
		}

		sum := n1.Val + n2.Val
		overflow := 0
		if sum >= 10 {
			overflow = 1
		}
		r = &ListNode{
			Val:  sum % 10,
			Next: nil,
		}
		n1 = n1.Next
		n2 = n2.Next
		r.Next = &ListNode{
			Val:  overflow,
			Next: nil,
		}
		r = r.Next
	}
	return result
}
