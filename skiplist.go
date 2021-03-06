// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package skiplist

import (
	"bytes"
	"math/rand"
)

//Node is the main type contained inside the StringSk
type Node struct {
	Next  []*Node
	Value []byte
}

// returns the height of the StringSk
func (n *Node) height() int {
	h := len(n.Next)

	if h == 0 {
		return 0
	}

	return h - 1
}

func newNode(height int) *Node {
	return &Node{Next: make([]*Node, height)}
}

// SkipList is the SkipList data structure
type SkipList struct {
	sentinel  *Node
	nodeCount uint
}

//New creates new SkipList
func New() *SkipList {
	return &SkipList{&Node{Next: make([]*Node, 1)}, 0}
}

// Size returns the nodeCount of Nodes in the StringSk
// except for the sentinel Node
func (s *SkipList) Size() uint {
	return s.nodeCount
}

// Height returns the current height of the StringSk
func (s *SkipList) Height() int {
	return len(s.sentinel.Next) - 1
}

func (s *SkipList) findPrev(value []byte) *Node {
	n := s.sentinel
	h := n.height()
	for ; h >= 0; h-- {
		for n.Next[h] != nil && bytes.Compare(n.Next[h].Value, value) < 0 {
			n = n.Next[h]
		}
	}

	return n
}

// Find tries to look for a value and returns the tuple value, boolean
// true if the value was found false if it wasnt' found
func (s *SkipList) Find(value []byte) bool {
	n := s.findPrev(value)
	if n.Next[0] != nil && bytes.Equal(n.Next[0].Value, value) {
		return true
	}

	return false
}

// RangeFind does a range query from start element till end element returns
// success or failure in form of a boolean and a the list of found values
// fails optmistiacally meaning if the start value is not found the query
// does not start at all, if the end value is not found the query run till
// a bigger value is found or there are no more elements on the list
func (s SkipList) RangeFind(start []byte, end []byte) (ok bool, found [][]byte) {
	n := s.findPrev(start)
	if n.Next[0] != nil && bytes.Equal(n.Next[0].Value, start) {
		for ; n.Next[0] != nil; n = n.Next[0] {
			found = append(found, n.Next[0].Value)
			if bytes.Equal(n.Next[0].Value, end) {
				return true, found
			}
			if bytes.Compare(n.Next[0].Value, end) > 0 {
				return false, found
			}
		}
	}

	return
}

func (s *SkipList) pickHeight() int {
	z := rand.Intn(39751)
	var (
		k int
		m = 1
	)

	for (z & m) != 0 {
		k++
		m <<= 1
	}

	return int(k) + 1
}

// Insert a new value and returns true or false based on success or failure
func (s *SkipList) Insert(value []byte) bool {
	n := s.sentinel
	h := s.sentinel.height()
	stack := make([]*Node, h+1)

	for ; h >= 0; h-- {
		for n.Next[h] != nil && bytes.Compare(n.Next[h].Value, value) < 0 {
			n = n.Next[h]
		}
		if n.Next[h] != nil && bytes.Equal(n.Next[h].Value, value) {
			return false
		}
		stack[h] = n
	}

	new := newNode(s.pickHeight())
	new.Value = value
	for s.sentinel.height() < new.height() {
		stack = append(stack, make([]*Node, 1)...)
		s.sentinel.Next = append(s.sentinel.Next, make([]*Node, 1)...)
		// basically increamenting stack and StringSk height
		stack[s.sentinel.height()] = s.sentinel
	}

	for i := 0; i < len(new.Next); i++ {
		new.Next[i] = stack[i].Next[i]
		stack[i].Next[i] = new
	}

	s.nodeCount++

	return true
}

// Remove a new value and returns true or false based on success or failure
func (s *SkipList) Remove(value []byte) (removed bool) {
	n := s.sentinel
	h := s.sentinel.height()

	for ; h >= 0; h-- {
		for n.Next[h] != nil && bytes.Compare(n.Next[h].Value, value) < 0 {
			n = n.Next[h]
		}
		if n.Next[h] != nil && bytes.Equal(n.Next[h].Value, value) {
			n.Next[h] = n.Next[h].Next[h]
			if n == s.sentinel && n.Next[h] == nil {
				s.sentinel.Next = s.sentinel.Next[:s.Height()]
			}

			removed = true
		}
	}

	if removed {
		s.nodeCount--
	}

	return removed
}
