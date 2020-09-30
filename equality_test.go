package equality_test

import (
	"fmt"
	"testing"
	"time"

	"bitbucket.org/michael-whatcott/equality"
)

var now = time.Now()

var notUTC, _ = time.LoadLocation("America/Los_Angeles")

func TestGeneralEquality(t *testing.T) {
	cases := []struct {
		Skip     bool
		Expected interface{}
		Actual   interface{}
		AreEqual bool
	}{
		{
			Skip:     false,
			Expected: 0,
			Actual:   0,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: 0,
			Actual:   1,
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: 0.0,
			Actual:   0.0,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: Thing{},
			Actual:   Thing{},
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: Thing{},
			Actual:   Thing{Integer: 1},
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: &Thing{Integer: 2},
			Actual:   &Thing{Integer: 2},
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: []int{1, 2, 3},
			Actual:   []int{1, 2, 3},
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: [3]int{1, 2, 3},
			Actual:   [3]int{1, 2, 3},
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: map[int]int{1: 2},
			Actual:   map[int]int{1: 2},
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: true,
			Actual:   true,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: make(chan int),
			Actual:   make(chan int),
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: "hi",
			Actual:   "hi",
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: now.In(notUTC),
			Actual:   now.In(time.UTC),
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: now.In(notUTC),
			Actual:   now.UTC().Add(time.Nanosecond),
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: int32(0),
			Actual:   0,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: int32(4),
			Actual:   4.0,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: uint32(4),
			Actual:   4.0,
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: complex128(4),
			Actual:   4.0,
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: int32(0),
			Actual:   1,
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: (func())(nil),
			Actual:   (func())(nil),
			AreEqual: true,
		},
		{
			Skip:     false,
			Expected: func() {},
			Actual:   (func())(nil),
			AreEqual: false,
		},
		{
			Skip:     false,
			Expected: func() {},
			Actual:   func() {},
			AreEqual: false,
		},
	}
	for x, test := range cases {
		title := fmt.Sprintf(
			"%d.Equal(%+v,%+v)==%t",
			x,
			test.Expected,
			test.Actual,
			test.AreEqual,
		)
		t.Run(title, func(t *testing.T) {
			if test.Skip {
				t.Skip()
			}
			if test.AreEqual {
				_ = equality.T(t).Assert(test.Expected, test.Actual)
			} else {
				report := equality.Report(test.Expected, test.Actual)
				if report == "" {
					t.Errorf("unequal values %v and %v erroneously deemed equal", test.Expected, test.Actual)
				} else {
					t.Log("(report printed below for visual inspection)", report)
				}
			}
		})
	}
}

type Thing struct {
	Integer int
}
