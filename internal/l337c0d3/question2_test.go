package l337c0d3

import (
	"fmt"
	"testing"
)

func Test_newListNodeFromString(t *testing.T) {
	l, err := newListNodeFromString("1234")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(l.String())
}

func Test_addTwoNumbers(t *testing.T) {
	var l1 *ListNode
	var l2 *ListNode

	r := addTwoNumbers(l1, l2)
	if r != nil {
		t.Error("add two nil list must emit error")
	}

	var err error
	l1, err = newListNodeFromString("1234")
	if err != nil {
		t.Error(err)
	}

	l2 = nil
	r = addTwoNumbers(l1, l2)
	fmt.Println(r.String())

	l1, err = newListNodeFromString("876543211")
	if err != nil {
		t.Error(err)
	}
	l2, err = newListNodeFromString("123456789")
	if err != nil {
		t.Error(err)
	}
	r = addTwoNumbers(l1, l2)
	fmt.Println(r.String())

	l1, err = newListNodeFromString("342")
	if err != nil {
		t.Error(err)
	}
	l2, err = newListNodeFromString("865")
	if err != nil {
		t.Error(err)
	}
	r = addTwoNumbers(l1, l2)
	fmt.Println(r.String())
}
