package main

type BinaryEncoderDecoder interface {
	Encode(text string) []byte
	Decode(binary []byte) string
}

type HuffmanCodeEncoderDecoder struct {

}