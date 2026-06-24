# huffman

Implementation of a file compressor using vanilla huffman coding for utf-8 text files.

Exposes a small cli `hfcode` for encoding files into a custom file format `.hfcode` and decoding them back to its original contents.

## Setup

```bash
./setup.sh
```

## Usage

Compress a file. The `output` file must end in `.hfcode`.

```bash
hfcode encode <file> <output>
```

Decompress a `.hfcode` file.

```bash
hfcode decode <file> <output>
```

Get stats on compression. Try running this command to see how much we compressed `./examples/tiny-shakespeare.txt`!

```bash
hfcode stats <file>

# example output
Original bits: 3840 bits
Compressed bits: 2543 bits
Compression ratio: 1.51
Compression percentage: 33.78%
Bits per symbol: 5.32 bits/symbol
```

## File Format

A hfcode file consists of a header containing the frequency of each unicode character present in the original file and a body containing the compressed bytes. The header can be decoded into a huffman tree which can be used to decode the compressed bytes.

Because files store bytes, the last byte may contain extra padding that do not store encoded bits.

```
[header]
[4 bytes | N number of unique unicode characters uint32]
[4 bytes | unicode codepoint uint32] [8 bytes | frequency in uint64]
... x N

[body]
[compressed bytes]
```
