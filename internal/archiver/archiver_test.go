package archiver

import (
	"bytes"
	"math/rand"
	"testing"
)

func generateRandomData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data)
	return data
}

func TestCompressAndDecompress_Integration(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty file", []byte{}},
		{"One byte", []byte{'A'}},
		{"One unique char (1000 times)", bytes.Repeat([]byte{'B'}, 1000)},
		{"Two chars", bytes.Repeat([]byte{'A'}, 500)},
		{"Small text", []byte("hello world")},
		{"All 256 bytes", generateRandomData(256)},
		{"Random 10KB", generateRandomData(10 * 1024)},
		{"Random 1MB", generateRandomData(1 * 1024 * 1024)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var compressedBuf bytes.Buffer
			var decompressedBuf bytes.Buffer

			originalReader := bytes.NewReader(tt.data)

			err := Compress(originalReader, &compressedBuf)
			if err != nil {
				t.Fatalf("Compress failed: %v", err)
			}

			if len(tt.data) == 0 {
				if compressedBuf.Len() != 13 {
					t.Errorf("Empty file compressed size expected 13, got %d", compressedBuf.Len())
				}
				return
			}

			compressedReader := bytes.NewReader(compressedBuf.Bytes())
			err = Decompress(compressedReader, &decompressedBuf)
			if err != nil {
				t.Fatalf("Decompress failed: %v", err)
			}

			if !bytes.Equal(tt.data, decompressedBuf.Bytes()) {
				t.Errorf("Mismatch! Original len=%d, Decompressed len=%d", len(tt.data), decompressedBuf.Len())
			}
		})
	}
}
