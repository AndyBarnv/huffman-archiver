package bitio

import "io"

type BitWriter struct {
	writer io.Writer
	byte   byte
	count  uint8
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{writer: w}
}

func (bw *BitWriter) WriteBit(bit bool) error {
	bw.count++
	if bit {
		bw.byte |= 1 << (8 - bw.count)
	}
	if bw.count == 8 {
		_, err := bw.writer.Write([]byte{bw.byte})
		if err != nil {
			return err
		}
		bw.byte = 0
		bw.count = 0
	}
	return nil
}

func (bw *BitWriter) Flush() error {
	if bw.count > 0 {
		_, err := bw.writer.Write([]byte{bw.byte})
		if err != nil {
			return err
		}
		bw.byte = 0
		bw.count = 0
	}
	return nil
}

type BitReader struct {
	reader io.Reader
	byte   byte
	count  uint8
	err    error
}

func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{reader: r}
}

func (br *BitReader) ReadBit() (bool, error) {
	if br.err != nil {
		return false, br.err
	}

	if br.count == 0 {
		buf := make([]byte, 1)
		_, br.err = br.reader.Read(buf)
		if br.err != nil {
			return false, br.err
		}
		br.byte = buf[0]
		br.count = 8
	}

	br.count--
	bit := (br.byte & (1 << br.count)) != 0
	return bit, nil
}
