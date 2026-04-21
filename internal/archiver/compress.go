package archiver

import (
	"bufio"
	"encoding/binary"
	"huffman-archiver/internal/bitio"
	"huffman-archiver/internal/core"
	"io"
)

func Compress(input io.ReadSeeker, output io.Writer) error {
	buf := make([]byte, 32*1024)
	freq := make(map[byte]uint64)
	dataSize := uint64(0)

	for {
		n, err := input.Read(buf)
		for i := 0; i < n; i++ {
			freq[buf[i]]++
			dataSize++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	tree := core.BuildTree(freq)
	codes := core.BuildCanonicalCodes(tree)

	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	if _, err := bufWriter.Write([]byte(core.FileSignature)); err != nil {
		return err
	}

	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, dataSize)
	if _, err := bufWriter.Write(sizeBuf); err != nil {
		return err
	}

	codesNumBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(codesNumBuf, uint16(len(codes)))
	if _, err := bufWriter.Write(codesNumBuf); err != nil {
		return err
	}

	if _, err := input.Seek(0, io.SeekStart); err != nil {
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
	for {
		n, err := input.Read(buf)
		for i := 0; i < n; i++ {
			ci := codeMap[buf[i]]
			for j := uint8(0); j < ci.Len; j++ {
				bit := (ci.Code >> (ci.Len - 1 - j)) & 1
				if err := writer.WriteBit(bit == 1); err != nil {
					return err
				}
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return bufWriter.Flush()
}
