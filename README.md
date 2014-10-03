vecstring
=========

vecstring is a Go library for a vector representation of strings.

Conceptually, vecstring represents a vector V[0...num), and each
value V[i] represents a string of variable length.
Strings are concatenated in a byte array, and each offset and its length are
encoded by unary coding, and stored using rsdic package.

Since rsdic supports constant time decode of unary code (Select and RunZeros)
the offset and its length can be decoded in constant time, and vecstring can retrieve
any string in constant time.

For strings of total lengths = TotalLen, VecString stores them in
at most TotalLen + Num * (2 + log_2 TotalLen / Num) / 8 bytes.

Usage
=====

```
import "github.com/hillbig/vecstring"

vs := vecstring.New()  // vs represents a vector of string
vs.PushBack("Hello") // Push String
vs.PushBack("New")
vs.PushBack("World")

fmt.Printf("%d %d\n", vs.Num(), vs.TotalLen()) // 3 13

for i := 0; i < vs.Num(); i++ {
	fmt.Printf("%d:%s\n", i, vs.Get(i))
}
// Hello
// New
// World

bytes, err := vs.MarshalBinary() // Encode to binary representation
newvs := vecstring.New()
err := newvs.UnmarshalBinary(bytes) // Decode from binary presentation
```
