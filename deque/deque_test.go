// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package deque_test

import (
	"container/list"
	"testing"

	tc "github.com/howbazaar/checkers"

	"github.com/juju/collections/deque"
)

type suite struct {
	tc.Test
	deque *deque.Deque
}

func TestDeque(t *testing.T) {
	tc.RunSuite(t, &suite{})
}

const testLen = 1000

func (s *suite) SetUpTest() {
	s.deque = deque.New()
}

func (s *suite) TestInit() {
	s.checkEmpty()
}

func (s *suite) TestStackBack() {
	// Push many values on to the back.
	for i := 0; i < testLen; i++ {
		s.Assert(s.deque.Len(), tc.Equals, i)
		s.deque.PushBack(i)
	}

	// Pop them all off from the back.
	for i := testLen; i > 0; i-- {
		s.Assert(s.deque.Len(), tc.Equals, i)
		v, ok := s.deque.PopBack()
		s.Assert(ok, tc.IsTrue)
		s.Assert(v.(int), tc.Equals, i-1)
	}

	s.checkEmpty()
}

func (s *suite) TestStackFront() {
	// Push many values on to the front.
	for i := 0; i < testLen; i++ {
		s.Assert(s.deque.Len(), tc.Equals, i)
		s.deque.PushFront(i)
	}

	// Pop them all off from the front.
	for i := testLen; i > 0; i-- {
		s.Assert(s.deque.Len(), tc.Equals, i)
		v, ok := s.deque.PopFront()
		s.Assert(ok, tc.IsTrue)
		s.Assert(v.(int), tc.Equals, i-1)
	}

	s.checkEmpty()
}

func (s *suite) TestQueueFromFront() {
	// Push many values on to the back.
	for i := 0; i < testLen; i++ {
		s.deque.PushBack(i)
	}

	// Pop them all off the front.
	for i := 0; i < testLen; i++ {
		v, ok := s.deque.PopFront()
		s.Assert(ok, tc.IsTrue)
		s.Assert(v.(int), tc.Equals, i)
	}

	s.checkEmpty()
}

func (s *suite) TestQueueFromBack() {
	// Push many values on to the front.
	for i := 0; i < testLen; i++ {
		s.deque.PushFront(i)
	}

	// Pop them all off the back.
	for i := 0; i < testLen; i++ {
		v, ok := s.deque.PopBack()
		s.Assert(ok, tc.IsTrue)
		s.Assert(v.(int), tc.Equals, i)
	}

	s.checkEmpty()
}

func (s *suite) TestFrontEmpty() {
	v, ok := s.deque.Front()
	s.Assert(ok, tc.IsFalse)
	s.Assert(v, tc.IsNil)
}

func (s *suite) TestFrontValue() {
	s.deque.PushFront(42)
	v, ok := s.deque.Front()
	s.Assert(ok, tc.IsTrue)
	s.Assert(v.(int), tc.Equals, 42)
	// Item is still there.
	s.Assert(s.deque.Len(), tc.Equals, 1)
}

func (s *suite) TestFrontBack() {
	// Populate from the front and back.
	for i := 0; i < testLen; i++ {
		s.Assert(s.deque.Len(), tc.Equals, i*2)
		s.deque.PushFront(i)
		s.deque.PushBack(i)
	}

	//  Remove half the items from the front and back.
	for i := testLen; i > testLen/2; i-- {
		s.Assert(s.deque.Len(), tc.Equals, i*2)

		vb, ok := s.deque.PopBack()
		s.Assert(ok, tc.IsTrue)
		s.Assert(vb.(int), tc.Equals, i-1)

		vf, ok := s.deque.PopFront()
		s.Assert(ok, tc.IsTrue)
		s.Assert(vf.(int), tc.Equals, i-1)
	}

	// Expand out again.
	for i := testLen / 2; i < testLen; i++ {
		s.Assert(s.deque.Len(), tc.Equals, i*2)
		s.deque.PushFront(i)
		s.deque.PushBack(i)
	}

	// Consume all.
	for i := testLen; i > 0; i-- {
		s.Assert(s.deque.Len(), tc.Equals, i*2)

		vb, ok := s.deque.PopBack()
		s.Assert(ok, tc.IsTrue)
		s.Assert(vb.(int), tc.Equals, i-1)

		vf, ok := s.deque.PopFront()
		s.Assert(ok, tc.IsTrue)
		s.Assert(vf.(int), tc.Equals, i-1)
	}

	s.checkEmpty()
}

func (s *suite) TestMaxLenFront() {
	const max = 5
	d := deque.NewWithMaxLen(max)

	// Exceed the maximum length by 2
	for i := 0; i < max+2; i++ {
		d.PushFront(i)
	}

	// Observe the the first 2 items on the back were dropped.
	v, ok := d.PopBack()
	s.Assert(ok, tc.IsTrue)
	s.Assert(v.(int), tc.Equals, 2)
}

func (s *suite) TestMaxLenBack() {
	const max = 5
	d := deque.NewWithMaxLen(max)

	// Exceed the maximum length by 3
	for i := 0; i < max+3; i++ {
		d.PushBack(i)
	}

	// Observe the the first 3 items on the front were dropped.
	v, ok := d.PopFront()
	s.Assert(ok, tc.IsTrue)
	s.Assert(v.(int), tc.Equals, 3)
}

func (s *suite) TestBlockAllocation() {
	// This test confirms that the Deque allocates and deallocates
	// blocks as expected.

	for i := 0; i < testLen; i++ {
		s.deque.PushFront(i)
		s.deque.PushBack(i)
	}
	// 2000 items at a blockLen of 64:
	// 31 full blocks + 1 partial front + 1 partial back = 33
	s.Assert(deque.GetDequeBlocks(s.deque), tc.Equals, 33)

	for i := 0; i < testLen; i++ {
		s.deque.PopFront()
		s.deque.PopBack()
	}
	// At empty there should be just 1 block.
	s.Assert(deque.GetDequeBlocks(s.deque), tc.Equals, 1)
}

func (s *suite) checkEmpty() {
	s.Assert(s.deque.Len(), tc.Equals, 0)

	_, ok := s.deque.PopFront()
	s.Assert(ok, tc.IsFalse)

	_, ok = s.deque.PopBack()
	s.Assert(ok, tc.IsFalse)
}

func BenchmarkPushBackList(t *testing.B) {
	l := list.New()
	for i := 0; i < t.N; i++ {
		l.PushBack(i)
	}
}

func BenchmarkPushBackDeque(t *testing.B) {
	d := deque.New()
	for i := 0; i < t.N; i++ {
		d.PushBack(i)
	}
}

func BenchmarkPushFrontList(t *testing.B) {
	l := list.New()
	for i := 0; i < t.N; i++ {
		l.PushFront(i)
	}
}

func BenchmarkPushFrontDeque(t *testing.B) {
	d := deque.New()
	for i := 0; i < t.N; i++ {
		d.PushFront(i)
	}
}

func BenchmarkPushPopFrontList(t *testing.B) {
	l := list.New()
	for i := 0; i < t.N; i++ {
		l.PushFront(i)
	}
	for i := 0; i < t.N; i++ {
		elem := l.Front()
		_ = elem.Value
		l.Remove(elem)
	}
}

func BenchmarkPushPopFrontDeque(t *testing.B) {
	d := deque.New()
	for i := 0; i < t.N; i++ {
		d.PushFront(i)
	}
	for i := 0; i < t.N; i++ {
		_, _ = d.PopFront()
	}
}

func BenchmarkPushPopBackList(t *testing.B) {
	l := list.New()
	for i := 0; i < t.N; i++ {
		l.PushBack(i)
	}
	for i := 0; i < t.N; i++ {
		elem := l.Back()
		_ = elem.Value
		l.Remove(elem)
	}
}

func BenchmarkPushPopBackDeque(t *testing.B) {
	d := deque.New()
	for i := 0; i < t.N; i++ {
		d.PushBack(i)
	}
	for i := 0; i < t.N; i++ {
		_, _ = d.PopBack()
	}
}

func iterToSlice(iter deque.Iterator) []string {
	var result []string
	var value string
	for iter.Next(&value) {
		result = append(result, value)
	}
	return result
}

func (s *suite) TestIterEmpty() {
	s.Assert(iterToSlice(s.deque.Iterator()), tc.HasLen, 0)
}

func (s *suite) TestIter() {
	s.deque.PushFront("second")
	s.deque.PushBack("third")
	s.deque.PushFront("first")
	s.Assert(iterToSlice(s.deque.Iterator()), tc.DeepEquals, []string{"first", "second", "third"})
}

func (s *suite) TestIterOverBlocksBack() {
	for i := 0; i < testLen; i++ {
		s.deque.PushBack(i)
	}

	iter := s.deque.Iterator()
	expect := 0
	var obtained int
	for iter.Next(&obtained) {
		s.Assert(obtained, tc.Equals, expect)
		expect++
	}
	s.Assert(expect, tc.Equals, testLen)
}

func (s *suite) TestIterOverBlocksFront() {
	for i := 0; i < testLen; i++ {
		s.deque.PushFront(i)
	}

	iter := s.deque.Iterator()
	expect := testLen - 1
	var obtained int
	for iter.Next(&obtained) {
		s.Assert(obtained, tc.Equals, expect)
		expect--
	}
	s.Assert(expect, tc.Equals, -1)
}

func (s *suite) TestWrongTypePanics() {
	// Use an int in the deque, and try to get strings out using the iterToSlice method.
	s.deque.PushFront(14)

	// TODO: impl panic matches
	// s.Assert(func() {
	// 	iterToSlice(s.deque.Iterator())
	// }, gc.PanicMatches, "reflect.Set: value of type int is not assignable to type string")

	// s.Assert(func() {
	// 	iter := s.deque.Iterator()
	// 	var i int
	// 	iter.Next(i)
	// }, gc.PanicMatches, "value is not a pointer")
}
