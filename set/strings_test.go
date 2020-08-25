// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package set_test

import (
	"sort"
	"testing"

	tc "github.com/howbazaar/checkers"

	"github.com/juju/collections/set"
)

type stringSetSuite struct {
	tc.Test
}

func TestStringSet(t *testing.T) {
	tc.RunSuite(t, &stringSetSuite{})
}

// Helper methods for the tests.
func (s *stringSetSuite) AssertValues(set set.Strings, expected ...string) {
	values := set.Values()
	// Expect an empty slice, not a nil slice for values.
	if expected == nil {
		expected = []string{}
	}
	sort.Strings(expected)
	sort.Strings(values)
	s.Assert(values, tc.DeepEquals, expected)
	s.Assert(set.Size(), tc.Equals, len(expected))
	// Check the sorted values too.
	sorted := set.SortedValues()
	s.Assert(sorted, tc.DeepEquals, expected)
}

// Actual tests start here.

func (s *stringSetSuite) TestEmpty() {
	s.AssertValues(set.NewStrings())
}

func (s *stringSetSuite) TestInitialValues() {
	values := []string{"foo", "bar", "baz"}
	v := set.NewStrings(values...)
	s.AssertValues(v, values...)
}

func (s *stringSetSuite) TestSize() {
	// Empty sets are empty.
	v := set.NewStrings()
	s.Assert(v.Size(), tc.Equals, 0)

	// Size returns number of unique values.
	v = set.NewStrings("foo", "foo", "bar")
	s.Assert(v.Size(), tc.Equals, 2)
}

func (s *stringSetSuite) TestIsEmpty() {
	// Empty sets are empty.
	v := set.NewStrings()
	s.Assert(v.IsEmpty(), tc.IsTrue)

	// Non-empty sets are not empty.
	v = set.NewStrings("foo")
	s.Assert(v.IsEmpty(), tc.IsFalse)
	// Newly empty sets work too.
	v.Remove("foo")
	s.Assert(v.IsEmpty(), tc.IsTrue)
}

func (s *stringSetSuite) TestAdd() {
	v := set.NewStrings()
	v.Add("foo")
	v.Add("foo")
	v.Add("bar")
	s.AssertValues(v, "foo", "bar")
}

func (s *stringSetSuite) TestRemove() {
	v := set.NewStrings("foo", "bar")
	v.Remove("foo")
	s.AssertValues(v, "bar")
}

func (s *stringSetSuite) TestContains() {
	v := set.NewStrings("foo", "bar")
	s.Assert(v.Contains("foo"), tc.IsTrue)
	s.Assert(v.Contains("bar"), tc.IsTrue)
	s.Assert(v.Contains("baz"), tc.IsFalse)
}

func (s *stringSetSuite) TestRemoveNonExistent() {
	v := set.NewStrings()
	v.Remove("foo")
	s.AssertValues(v)
}

func (s *stringSetSuite) TestUnion() {
	s1 := set.NewStrings("foo", "bar")
	s2 := set.NewStrings("foo", "baz", "bang")
	union1 := s1.Union(s2)
	union2 := s2.Union(s1)

	s.AssertValues(union1, "foo", "bar", "baz", "bang")
	s.AssertValues(union2, "foo", "bar", "baz", "bang")
}

func (s *stringSetSuite) TestIntersection() {
	s1 := set.NewStrings("foo", "bar")
	s2 := set.NewStrings("foo", "baz", "bang")
	int1 := s1.Intersection(s2)
	int2 := s2.Intersection(s1)

	s.AssertValues(int1, "foo")
	s.AssertValues(int2, "foo")
}

func (s *stringSetSuite) TestDifference() {
	s1 := set.NewStrings("foo", "bar")
	s2 := set.NewStrings("foo", "baz", "bang")
	diff1 := s1.Difference(s2)
	diff2 := s2.Difference(s1)

	s.AssertValues(diff1, "bar")
	s.AssertValues(diff2, "baz", "bang")
}

func (s *stringSetSuite) TestUninitializedPanics() {
	f := func() {
		var s set.Strings
		s.Add("foo")
	}
	s.Assert(f, tc.PanicMatches, "uninitalised set")
}
