package archiver

import (
	"bufio"
	"encoding/binary"
	"errors"
	"huffman-archiver/internal/bitio"
	"huffman-archiver/internal/core"
	"io"
)

// ErrInvalidArchive возвращается при попытке разжать файл,
// не соответствующий формату архива (неверная сигнатура или поврежденные данные).
var ErrInvalidArchive = errors.New("invalid archive format")

// Decompress выполняет восстановление оригинальных данных из сжатого потока.
// В отличие от Compress, принимает обычный io.Reader, так как делает только один проход
// и не нуждается в возврате указателя чтения. Обеспечивает O(1) потребление памяти.
func Decompress(input io.Reader, output io.Writer) error {
	// Буферизуем вывод. При побитовом чтении мы будем распаковывать данные
	// по одному байту, и bufio предотвратит обращение к диску на каждый распакованный символ.
	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	// Чтение заголовка
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

	// Выделяем память под все записи таблицы и читаем её одним куском
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

	// Восстанавливаем дерево для декодирования.
	tree := core.BuildDecodingTree(codes)

	// Побитовое декодирование
	reader := bitio.NewBitReader(input)
	var written uint64

	for written < size {
		node := tree

		// Спускаемся по дереву до листового узла
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

		// Записываем распакованный символ из листа в буфер вывода
		if err := bufWriter.WriteByte(node.Char); err != nil {
			return err
		}
		written++
	}

	return nil
}
