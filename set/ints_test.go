// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package set_test

import (
	"sort"
	"testing"

	tc "github.com/howbazaar/checkers"

	"github.com/juju/collections/set"
)

type intSetSuite struct {
	tc.Test
}

func TestIntSet(t *testing.T) {
	tc.RunSuite(t, &intSetSuite{})
}

// Helper methods for the tests.
func (i *intSetSuite) AssertIntValues(s set.Ints, expected ...int) {
	values := s.Values()
	// Expect an empty slice, not a nil slice for values.
	if expected == nil {
		expected = []int{}
	}
	sort.Ints(expected)
	sort.Ints(values)
	i.Assert(values, tc.DeepEquals, expected)
	i.Assert(s.Size(), tc.Equals, len(expected))
	// Check the sorted values too.
	sorted := s.SortedValues()
	i.Assert(sorted, tc.DeepEquals, expected)
}

// Actual tests start here.

func (i *intSetSuite) TestEmpty() {
	s := set.NewInts()
	i.AssertIntValues(s)
}

func (i *intSetSuite) TestInitialValues() {
	values := []int{1, 2, 3}
	s := set.NewInts(values...)
	i.AssertIntValues(s, values...)
}

func (i *intSetSuite) TestSize() {
	// Empty sets are empty.
	s := set.NewInts()
	i.Assert(s.Size(), tc.Equals, 0)

	// Size returns number of unique values.
	s = set.NewInts(1, 1, 2)
	i.Assert(s.Size(), tc.Equals, 2)
}

func (i *intSetSuite) TestIsEmpty() {
	// Empty sets are empty.
	s := set.NewInts()
	i.Assert(s.IsEmpty(), tc.IsTrue)

	// Non-empty sets are not empty.
	s = set.NewInts(1)
	i.Assert(s.IsEmpty(), tc.IsFalse)
	// Newly empty sets work too.
	s.Remove(1)
	i.Assert(s.IsEmpty(), tc.IsTrue)
}

func (i *intSetSuite) TestAdd() {
	s := set.NewInts()
	s.Add(1)
	s.Add(1)
	s.Add(2)
	i.AssertIntValues(s, 1, 2)
}

func (i *intSetSuite) TestRemove() {
	s := set.NewInts(1, 2)
	s.Remove(1)
	i.AssertIntValues(s, 2)
}

func (i *intSetSuite) TestContains() {
	s := set.NewInts(1, 2)
	i.Assert(s.Contains(1), tc.IsTrue)
	i.Assert(s.Contains(2), tc.IsTrue)
	i.Assert(s.Contains(3), tc.IsFalse)
}

func (i *intSetSuite) TestRemoveNonExistent() {
	s := set.NewInts()
	s.Remove(1)
	i.AssertIntValues(s)
}

func (i *intSetSuite) TestUnion() {
	s1 := set.NewInts(1, 2)
	s2 := set.NewInts(1, 3, 4)
	union1 := s1.Union(s2)
	union2 := s2.Union(s1)

	i.AssertIntValues(union1, 1, 2, 3, 4)
	i.AssertIntValues(union2, 1, 2, 3, 4)
}

func (i *intSetSuite) TestIntersection() {
	s1 := set.NewInts(1, 2)
	s2 := set.NewInts(1, 3, 4)
	int1 := s1.Intersection(s2)
	int2 := s2.Intersection(s1)

	i.AssertIntValues(int1, 1)
	i.AssertIntValues(int2, 1)
}

func (i *intSetSuite) TestDifference() {
	s1 := set.NewInts(1, 2)
	s2 := set.NewInts(1, 3, 4)
	diff1 := s1.Difference(s2)
	diff2 := s2.Difference(s1)

	i.AssertIntValues(diff1, 2)
	i.AssertIntValues(diff2, 3, 4)
}

func (i *intSetSuite) TestUninitializedPanics() {
	f := func() {
		var s set.Ints
		s.Add(1)
	}
	i.Assert(f, tc.PanicMatches, "uninitalised set")
}
