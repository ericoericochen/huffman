package huffman

type Bit bool
type Bits []Bit

const (
	Zero Bit = false
	One Bit = true
)

func (bs Bits) String() string {
	rep := ""
	for _, bit := range bs {
		if bit {
			rep += "1"
		} else {
			rep += "0"	
		}
	}
	return rep
}