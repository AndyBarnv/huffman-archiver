package archiver

import (
	"bufio"
	"encoding/binary"
	"huffman-archiver/internal/bitio"
	"huffman-archiver/internal/core"
	"io"
)

func Compress(input io.Reader, output io.Writer) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return err
	}

	freq := make(map[byte]uint64)
	for _, b := range data {
		freq[b]++
	}

	tree := core.BuildTree(freq)
	codes := core.BuildCanonicalCodes(tree)

	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	if _, err := bufWriter.Write([]byte(core.FileSignature)); err != nil {
		return err
	}

	buf8 := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf8, uint64(len(data)))
	if _, err := bufWriter.Write(buf8); err != nil {
		return err
	}

	buf2 := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf2, uint16(len(codes)))
	if _, err := bufWriter.Write(buf2); err != nil {
		return err
	}

	codeMap := make(map[byte]core.CodeInfo, len(codes))
	tableBuf := make([]byte, 2*len(codes))
	for i, ci := range codes {
		codeMap[ci.Char] = ci
		tableBuf[i*2] = ci.Char
		tableBuf[i*2+1] = ci.Len
	}
	if _, err := bufWriter.Write(tableBuf); err != nil {
		return err
	}

	writer := bitio.NewBitWriter(bufWriter)
	for _, b := range data {
		ci := codeMap[b]
		for i := uint8(0); i < ci.Len; i++ {
			bit := (ci.Code >> (ci.Len - 1 - i)) & 1
			if err := writer.WriteBit(bit == 1); err != nil {
				return err
			}
		}
	}

	return writer.Flush()
}
