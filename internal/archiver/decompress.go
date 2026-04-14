package archiver

import (
	"bufio"
	"encoding/binary"
	"errors"
	"huffman-archiver/internal/bitio"
	"huffman-archiver/internal/core"
	"io"
)

var ErrInvalidArchive = errors.New("invalid archive format")

func Decompress(input io.Reader, output io.Writer) error {
	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	signature := make([]byte, 3)
	if _, err := io.ReadFull(input, signature); err != nil {
		return err
	}

	if string(signature) != core.FileSignature {
		return ErrInvalidArchive
	}

	sizeBuf := make([]byte, 8)
	if _, err := io.ReadFull(input, sizeBuf); err != nil {
		return err
	}
	size := binary.LittleEndian.Uint64(sizeBuf)

	codesNumBuf := make([]byte, 2)
	if _, err := io.ReadFull(input, codesNumBuf); err != nil {
		return err
	}
	codesNum := binary.LittleEndian.Uint16(codesNumBuf)

	codes := make([]core.CodeInfo, codesNum)

	tableBuf := make([]byte, 2*codesNum)
	if _, err := io.ReadFull(input, tableBuf); err != nil {
		return err
	}

	for i := 0; i < int(codesNum); i++ {
		codes[i].Char = tableBuf[i*2]
		codes[i].Len = tableBuf[i*2+1]
	}

	if size == 0 {
		return nil
	}

	tree := core.BuildDecodingTree(codes)

	reader := bitio.NewBitReader(input)
	var written uint64

	for written < size {
		node := tree

		for node.Left != nil || node.Right != nil {
			bit, err := reader.ReadBit()
			if err != nil {
				return err
			}
			if bit {
				node = node.Right
			} else {
				node = node.Left
			}
		}

		if err := bufWriter.WriteByte(node.Char); err != nil {
			return err
		}
		written++
	}

	return nil
}
