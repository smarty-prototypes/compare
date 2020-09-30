package equality_test

import (
	"errors"
	"testing"
	"time"
)

func TestFormatV(t *testing.T) {
	t.Skip("just for demonstration")

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
