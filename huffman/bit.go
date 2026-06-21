package huffman

type Bit bool
type Bits []Bit

const (
	Zero Bit = false
	One Bit = true
)

const BITS_IN_BYTE = 8

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

func GetBits(text string) int {
	return len(text) * BITS_IN_BYTE
}

func ByteToBits(b byte) Bits {
	var bits Bits = make([]Bit, 8)
	for i := 0; i < 8; i++ {
		bit := (b >> i) & 1
		bits[7 - i] = bit == 1
	}
	return bits
}

func StringifyByte(b byte) string {
	bits := ByteToBits(b)
	return bits.String()
}