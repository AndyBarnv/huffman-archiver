// Package archiver реализует высокоуровневую логику сжатия и разжатия
// файлов по алгоритму статического канонического Хаффмана.
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
//
// Формат выходного потока:
// [Сигнатура "HUF" (3 байта)] [Размер оригинала (8 байт)] [Кол-во символов (2 байта)]
// [Таблица длин кодов (2 байта на символ)] [Сжатые биты]
func Compress(input io.ReadSeeker, output io.Writer) error {
	// Буфер 32 КБ для чтения из файла кусками
	buf := make([]byte, 32*1024)
	freq := make(map[byte]uint64)
	dataSize := uint64(0)

	// Подсчет частот символов и общего размера
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

	// Строим дерево частот и генерируем канонические коды
	tree := core.BuildTree(freq)
	codes := core.BuildCanonicalCodes(tree)

	// Оборачиваем выходной поток в буфер для минимизации системных вызовов при записи
	bufWriter := bufio.NewWriter(output)
	defer bufWriter.Flush()

	// Запись заголовка
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

	// Возвращаем указатель чтения в начало файла для второго прохода
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Подготовка таблицы кодов и её запись
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

	// Побитовое кодирование данных
	writer := bitio.NewBitWriter(bufWriter)
	for {
		n, err := input.Read(buf)
		for i := 0; i < n; i++ {
			ci := codeMap[buf[i]]
			// Побитово извлекаем код текущего символа (от старшего бита к младшему)
			for j := uint8(0); j < ci.Len; j++ {
				// Сдвигаем код вправо так, чтобы текущий бит (j) оказался в самой правой позиции (младший бит),
				// затем применяем маску & 1, чтобы отбросить все остальные биты и получить только 0 или 1.
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
