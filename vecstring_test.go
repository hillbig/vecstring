package vecstring

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

const (
	testIter = 10
	numMax   = 20000
)

func TestVecString(t *testing.T) {
	Convey("When an input is empty", t, func() {
		vs := New()
		So(vs.Num(), ShouldEqual, 0)
		So(vs.TotalLen(), ShouldEqual, 0)
	})
	Convey("When three strings are set", t, func() {
		vs := NewForWX()
		vs.PushBack("abc")
		vs.PushBack("de")
		vs.PushBack("efgh")
		So(vs.GetByte(5), ShouldEqual, byte('e'))
		So(vs.ExactMatch(1, "de"), ShouldBeTrue)
		So(vs.ExactMatch(1, "dee"), ShouldBeFalse)
		l, ok := vs.PrefixMatch(2, "efgh")
		So(ok, ShouldBeTrue)
		So(l, ShouldEqual, 4)
		l, ok = vs.PrefixMatch(2, "efghi")
		So(ok, ShouldBeTrue)
		So(l, ShouldEqual, 4)
		l, ok = vs.PrefixMatch(2, "efggi")
		So(ok, ShouldBeFalse)
		So(l, ShouldEqual, 3)
	})
	Convey("When large VecString is build", t, func() {
		for iter := 0; iter < testIter; iter++ {
			num := rand.Int()%numMax + 1
			vs := New()
			orig := make([]string, num)
			offsets := make([]uint64, num)
			offset := uint64(0)
			for i := 0; i < num; i++ {
				offsets = append(offsets, offset)
				slen := uint32(rand.Int31() % 10)
				buf := make([]byte, slen)
				for j := uint32(0); j < slen; j++ {
					buf[j] = 'a'
				}
				orig[i] = string(buf)
				offset += uint64(len(orig[i]))
				vs.PushBack(orig[i])
			}
			for i := 0; i < 10; i++ {
				ind := uint64(rand.Int() % num)
				So(vs.Get(ind), ShouldEqual, orig[ind])
				if len(orig[ind]) == 0 {
					continue
				}
				pos := rand.Int() % len(orig[ind])
				So(vs.Find(ind, orig[ind][pos]), ShouldEqual, true)
				So(vs.Find(ind, 'b'), ShouldEqual, false)
			}
			out, err := vs.MarshalBinary()
			So(err, ShouldBeNil)
			newvs := New()
			err = newvs.UnmarshalBinary(out)
			So(err, ShouldBeNil)
			So(newvs.Num(), ShouldEqual, num)
			for i := 0; i < 10; i++ {
				ind := uint64(rand.Int() % num)
				So(newvs.Get(ind), ShouldEqual, orig[ind])
				if len(orig[ind]) == 0 {
					continue
				}
				pos := rand.Int() % len(orig[ind])
				So(newvs.Find(ind, orig[ind][pos]), ShouldEqual, true)
				So(newvs.Find(ind, 'b'), ShouldEqual, false)
			}
		}
	})

}
