package bitio

import (
	"bytes"
	"testing"
)

func TestBitIO_WriteAndRead(t *testing.T) {
	var buf bytes.Buffer
	writer := NewBitWriter(&buf)

	// Записываем 8 бит
	bits := []bool{true, false, true, true, false, false, false, true}
	for i, b := range bits {
		if err := writer.WriteBit(b); err != nil {
			t.Fatalf("WriteBit %d failed: %v", i, err)
		}
	}
	writer.Flush()

	// Читаем биты и сравниваем с ожидаемыми
	reader := NewBitReader(&buf)
	for i, expected := range bits {
		got, err := reader.ReadBit()
		if err != nil {
			t.Fatalf("ReadBit %d failed: %v", i, err)
		}
		if got != expected {
			t.Errorf("Bit %d: expected %v, got %v", i, expected, got)
		}
	}
}

func TestBitIO_FlushPadding(t *testing.T) {
	var buf bytes.Buffer
	writer := NewBitWriter(&buf)

	// Пишем только 3 бита: 1, 0, 1
	writer.WriteBit(true)
	writer.WriteBit(false)
	writer.WriteBit(true)
	writer.Flush() // Должен дописать 5 нулей: 10100000

	reader := NewBitReader(&buf)

	// Читаем первые 3 бита
	b1, err := reader.ReadBit()
	if err != nil {
		t.Fatalf("ReadBit failed: %v", err)
	}
	b2, err := reader.ReadBit()
	if err != nil {
		t.Fatalf("ReadBit failed: %v", err)
	}
	b3, err := reader.ReadBit()
	if err != nil {
		t.Fatalf("ReadBit failed: %v", err)
	}
	if !(b1 && !b2 && b3) {
		t.Errorf("Original bits corrupted")
	}

	// Читаем и проверяем оставшиеся 5 бит (должны быть дописанные нули)
	for i := 0; i < 5; i++ {
		padBit, err := reader.ReadBit()
		if err != nil {
			t.Fatalf("ReadBit failed: %v", err)
		}
		if padBit != false {
			t.Errorf("Flush did not pad with zeros! Bit %d is true", i)
		}
	}
}
