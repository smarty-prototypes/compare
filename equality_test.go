package equality_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/michael-whatcott/equality"
)

type Case struct {
	Skip     bool
	Expected interface{}
	Actual   interface{}
	AreEqual bool
	Options  []equality.Option
}

var now = time.Now()

var notUTC, _ = time.LoadLocation("America/Los_Angeles")

var (
	cases = []Case{
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
			Actual:   now.UTC(),
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
)

func TestEqual(t *testing.T) {
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

func TestFormatV(t *testing.T) {
	now := time.Now()

	t.Logf("%v", []int{1, 2, 3})  // [1 2 3]
	t.Logf("%+v", []int{1, 2, 3}) // [1 2 3]
	t.Logf("%#v", []int{1, 2, 3}) // []int{1, 2, 3}
	t.Log()
	t.Logf("%v", []time.Time{now})  // [2020-09-18 13:24:31.940914 -0600 MDT m=+0.000471278]
	t.Logf("%+v", []time.Time{now}) // [2020-09-18 13:24:31.940914 -0600 MDT m=+0.000471278]
	t.Logf("%#v", []time.Time{now}) // []time.Time{time.Time{wall:0xbfd1603bf8153550, ext:471278, loc:(*time.Location)(0x122d5c0)}}
	t.Log()
	t.Logf("%v", []*time.Time{&now})  // [2020-09-18 13:24:31.940914 -0600 MDT m=+0.000471278]
	t.Logf("%+v", []*time.Time{&now}) // [2020-09-18 13:24:31.940914 -0600 MDT m=+0.000471278]
	t.Logf("%#v", []*time.Time{&now}) // []*time.Time{(*time.Time)(0xc00000c200)}
	t.Log()
	t.Logf("%v", errors.New("hi"))  // hi
	t.Logf("%+v", errors.New("hi")) // hi
	t.Logf("%#v", errors.New("hi")) // &errors.errorString{s:"hi"}
	t.Log()
	t.Logf("%v", []error{errors.New("hi")})  // [hi]
	t.Logf("%+v", []error{errors.New("hi")}) // [hi]
	t.Logf("%#v", []error{errors.New("hi")}) // []error{(*errors.errorString)(0xc00008e8a0)}
	t.Log()
	t.Logf("%v", make(chan int))  // 0xc000020540
	t.Logf("%+v", make(chan int)) // 0xc000020540
	t.Logf("%#v", make(chan int)) // (chan int)(0xc000020540)
	t.Log()
	t.Logf("%v", map[int]int{1: 2})  // map[1:2]
	t.Logf("%+v", map[int]int{1: 2}) // map[1:2]
	t.Logf("%#v", map[int]int{1: 2}) // map[int]int{1:2}
	t.Log()
	t.Logf("%v", Thing{Integer: 42})  // {42}
	t.Logf("%+v", Thing{Integer: 42}) // {Integer: 42}
	t.Logf("%#v", Thing{Integer: 42}) // equality.Thing{Integer: 42}
	t.Log()
	t.Logf("%v", []interface{}{Thing{Integer: 42}})  // [{42}]
	t.Logf("%+v", []interface{}{Thing{Integer: 42}}) // [{Integer:42}]
	t.Logf("%#v", []interface{}{Thing{Integer: 42}}) // []interface {}{equality.Thing{Integer:42}}
	t.Log()
	t.Logf("%v", &Thing{Integer: 42})  // &{42}
	t.Logf("%+v", &Thing{Integer: 42}) // &{Integer: 42}
	t.Logf("%#v", &Thing{Integer: 42}) // &equality.Thing{Integer: 42}
	t.Log()
	t.Logf("%v", []*Thing{{Integer: 42}})  // [0xc000024968]
	t.Logf("%+v", []*Thing{{Integer: 42}}) // [0xc000024988]
	t.Logf("%#v", []*Thing{{Integer: 42}}) // []*equality.Thing{(*equality.Thing)(0xc0000249a8)}
	t.Log()
	t.Logf("%v", []interface{}{&Thing{Integer: 42}})  // [0xc000024968]
	t.Logf("%+v", []interface{}{&Thing{Integer: 42}}) // [0xc000024988]
	t.Logf("%#v", []interface{}{&Thing{Integer: 42}}) // []interface {}{(*equality.Thing)(0xc0000249a8)}
	t.Log()
}
