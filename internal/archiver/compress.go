package archiver

import (
	"bufio"
	"encoding/binary"
	"huffman-archiver/internal/bitio"
	"huffman-archiver/internal/core"
	"io"
)

// Compress выполняет сжатие данных из input в output.
// Функция требует io.ReadSeeker, так как использует два прохода по данным:
// первый для подсчета частот символов, второй для записи сжатого потока.
// Это гарантирует потребление O(1) оперативной памяти независимо от размера файла.
func Compress(input io.ReadSeeker, output io.Writer) error {
	frequencies, originalSize, err := countFrequencies(input)
	if err != nil {
		return err
	}

	tree := core.BuildTree(frequencies)
	codes := core.BuildCanonicalCodes(tree)

	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	if err := writeHeader(bufWriter, originalSize, codes); err != nil {
		return err
	}

	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return err
	}

	codeMap := buildCodeMap(codes)
	writer := bitio.NewBitWriter(bufWriter)

	if err := writeCompressedData(input, writer, codeMap); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return bufWriter.Flush()
}

// countFrequencies делает первый проход по данным, подсчитывая количество каждого байта.
func countFrequencies(input io.Reader) (map[byte]uint64, uint64, error) {
	buf := make([]byte, 32*1024)
	frequencies := make(map[byte]uint64)
	dataSize := uint64(0)

	for {
		n, err := input.Read(buf)
		for i := 0; i < n; i++ {
			frequencies[buf[i]]++
			dataSize++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, 0, err
		}
	}

	return frequencies, dataSize, nil
}

// writeHeader формирует и записывает структурированный заголовок архива.
func writeHeader(writer *bufio.Writer, dataSize uint64, codes []core.CodeInfo) error {
	if _, err := writer.Write([]byte(core.FileSignature)); err != nil {
		return err
	}

	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, dataSize)
	if _, err := writer.Write(sizeBuf); err != nil {
		return err
	}

	codesNumBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(codesNumBuf, uint16(len(codes)))
	if _, err := writer.Write(codesNumBuf); err != nil {
		return err
	}

	tableBuf := make([]byte, 2*len(codes))
	for i, ci := range codes {
		tableBuf[i*2] = ci.Char
		tableBuf[i*2+1] = ci.Len
	}

	_, err := writer.Write(tableBuf)
	return err
}

// buildCodeMap конвертирует слайс канонических кодов в хеш-таблицу для O(1) поиска по символу.
func buildCodeMap(codes []core.CodeInfo) map[byte]core.CodeInfo {
	codeMap := make(map[byte]core.CodeInfo, len(codes))
	for _, ci := range codes {
		codeMap[ci.Char] = ci
	}
	return codeMap
}

// writeCompressedData делает второй проход по данным, побитово записывая коды в выходной поток.
func writeCompressedData(input io.Reader, bitWriter *bitio.BitWriter, codeMap map[byte]core.CodeInfo) error {
	buf := make([]byte, 32*1024)

	for {
		n, err := input.Read(buf)
		for i := 0; i < n; i++ {
			ci := codeMap[buf[i]]
			for j := uint8(0); j < ci.Len; j++ {
				bit := (ci.Code >> (ci.Len - 1 - j)) & 1
				if err := bitWriter.WriteBit(bit == 1); err != nil {
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

	return nil
}
