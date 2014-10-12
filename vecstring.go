// package vecstring provides a space-efficient representation of a vector of strings of variable length.
package vecstring

import (
	"github.com/hillbig/rsdic"
	"github.com/ugorji/go/codec"
)

// VecString represents a vector V[0...Num), and each value V[i] represents a string
// Internally, VecString are stored in space-efficient way.
// Strings are concatenated in a byte array, and each offset, length information is
// encoded in unary coding, and represented by rsdic to achieve further space reduction. // For strings of total lengths = TotalLen, VecString stores them in
// at most TotalLen + Num * (2 + log_2 TotalLen / Num) / 8 bytes.
type VecString interface {
	// Get returns V[ind]
	Get(ind uint64) string

	// Find returns true if c is found in V[ind], and returns false otherwise.
	Find(ind uint64, c byte) bool

	// Num returns the number of strings
	Num() uint64

	// TotalLen returns the total length of strings
	TotalLen() uint64

	// ExactMatch returns true if V[ind] == str, and returns false otherwise
	ExactMatch(ind uint64, str string) bool

	// PrefixMatch returns true if V[ind] == Prefix of str, and returns false otherwise
	PrefixMatch(ind uint64, str string) (uint64, bool)

	// OffsetAndLen returns the offset and the length of V[ind]
	OffsetAndLen(ind uint64) (uint64, uint64)

	// MarshalBinary encodes VecString into a binary form and returns the result.
	MarshalBinary() ([]byte, error)

	// UnmarshalBinary decodes the FixVec form a binary from generated MarshalBinary
	UnmarshalBinary([]byte) error

	// PushBack set V[Num] = v, and Num++
	PushBack(v string)
}

// VecStringForWX represents a vector of strings, and provides extended interface
// for library WX.
type VecStringForWX interface {
	VecString

	// GetByte returns the (offset+1)-th byte in the concatenated strings.
	GetByte(offset uint64) byte

	// FindZeroRank finds c in V[ind], and returns (offset, true) if found,
	// and returns (0, false) otherwise.
	FindZeroRank(ind uint64, c byte) (uint64, bool)

	// IthCharInd returns the index corresponding to i-th child.
	// lens.Rank(lens.Select(i-1, false), true) - 1
	IthCharInd(i uint64) uint64
}

// New returns VecString with 0 strings.
func New() VecString {
	rsd := rsdic.New()
	rsd.PushBack(true)
	return &vecStringImpl{
		lens:  rsd,
		bytes: make([]byte, 0),
	}
}

// New returns VecStringForWX with 0 strings.
func NewForWX() VecStringForWX {
	rsd := rsdic.New()
	rsd.PushBack(true)
	return &vecStringImpl{
		lens:  rsd,
		bytes: make([]byte, 0),
	}
}

type vecStringImpl struct {
	lens  rsdic.RSDic
	bytes []byte
}

func (vv vecStringImpl) OffsetAndLen(ind uint64) (uint64, uint64) {
	onePos := vv.lens.Select(ind, true)
	nextOnePos := vv.lens.Select(ind+1, true)
	offset := onePos - ind // zeroNum
	l := nextOnePos - onePos - 1
	return offset, l
}

func (vv vecStringImpl) Get(ind uint64) string {
	beg, l := vv.OffsetAndLen(ind)
	return string(vv.bytes[beg : beg+l])
}

func (vv vecStringImpl) GetByte(offset uint64) byte {
	return vv.bytes[offset]
}

func (vv vecStringImpl) Num() uint64 {
	return vv.lens.OneNum() - 1
}

func (vv *vecStringImpl) TotalLen() uint64 {
	return vv.lens.ZeroNum()
}

func (vv vecStringImpl) IthCharInd(i uint64) uint64 {
	return vv.lens.Rank(vv.lens.Select(i-1, false), true) - 1
}

func (vv vecStringImpl) Find(ind uint64, c byte) bool {
	beg, l := vv.OffsetAndLen(ind)
	for i := uint64(0); i < l; i++ {
		if c == vv.bytes[beg+i] {
			return true
		}
	}
	return false
}

func (vv vecStringImpl) FindZeroRank(ind uint64, c byte) (uint64, bool) {
	beg, l := vv.OffsetAndLen(ind)
	for i := uint64(0); i < l; i++ {
		if c == vv.bytes[beg+i] {
			return beg + i, true
		}
	}
	return 0, false
}

func (vv vecStringImpl) ExactMatch(ind uint64, str string) bool {
	onePos := vv.lens.Select(ind, true)
	beg := onePos - ind
	l := vv.lens.RunZeros(onePos + 1)
	return string(vv.bytes[beg:beg+l]) == str
}

func (vv vecStringImpl) PrefixMatch(ind uint64, str string) (uint64, bool) {
	beg, l := vv.OffsetAndLen(ind)
	for i := uint64(0); i < l; i++ {
		if int(i) >= len(str) || vv.bytes[beg+i] != str[i] {
			return i, false
		}
	}
	return l, true
}

func (vb *vecStringImpl) PushBack(str string) {
	for i := 0; i < len(str); i++ {
		vb.lens.PushBack(false)
	}
	vb.lens.PushBack(true)
	vb.bytes = append(vb.bytes, []byte(str)...)
}

func (vs vecStringImpl) MarshalBinary() (out []byte, err error) {
	var bh codec.MsgpackHandle
	enc := codec.NewEncoderBytes(&out, &bh)
	err = enc.Encode(vs.lens)
	if err != nil {
		return
	}
	err = enc.Encode(vs.bytes)
	if err != nil {
		return
	}

	return
}

func (vs *vecStringImpl) UnmarshalBinary(in []byte) (err error) {
	var bh codec.MsgpackHandle
	dec := codec.NewDecoderBytes(in, &bh)
	err = dec.Decode(&vs.lens)
	if err != nil {
		return
	}
	err = dec.Decode(&vs.bytes)
	if err != nil {
		return
	}
	return nil
}
