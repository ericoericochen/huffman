package huffman

type Encoder interface {
	Encode(text string) []byte
}

type Decoder interface {
	Decode(binary []byte) string
}