# hfcode

Compress an utf-8 encoded text file into a `.hfcode` file.

```bash
hfcode encode <file>
```

# File Format

A hfcode file consists of a header containing the frequency of each unicode character present in the original file and a body containing the compressed bytes. The header can be decoded into a huffman tree which can be used to decode the compressed bytes.

The header tells us the number of unicode characters and their frequency.

In the body, we store the number of bits and the actual bits. The number of decodable bits may be shorter than the bytes present so there may be extra padding of 0s at the end.

```
[header]
[4 bytes | N number of unicode characters sorted alphabetically uint32]
[4 bytes | unicode code point int32] [8 bytes | frequency in uint64]
... x N

[body]
[8 bytes | M number of bits uint64]
[M bits + 0 Padding]
```
