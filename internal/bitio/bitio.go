// Package bitio предоставляет низкоуровневые утилиты для побитового
// чтения и записи данных в потоки, реализующие интерфейсы io.Reader и io.Writer.
// Используется для кодирования и декодирования потоков битов в алгоритме Хаффмана.
package bitio

import "io"

// BitWriter выполняет побитовую запись в базовый поток.
// Накапливает биты во внутреннем буфере (слева направо, MSB first)
// и сбрасывает их по мере заполнения байта.
type BitWriter struct {
	writer io.Writer
	byte   byte
	count  uint8
}

// NewBitWriter создает новый BitWriter для указанного потока вывода.
func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{writer: w}
}

// WriteBit записывает один бит в буфер.
func (bw *BitWriter) WriteBit(bit bool) error {
	bw.count++
	if bit {
		// Сдвигаем единицу на нужную позицию слева и применяем побитовое ИЛИ (|).
		bw.byte |= 1 << (8 - bw.count)
	}
	// Если буфер заполнился (собрали 8 бит), пишем его в поток
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

// Flush принудительно сбрасывает накопленные биты в поток.
// Если в буфере осталось меньше 8 бит, недостающие биты
// справа заполняются нулями (выравнивание до границы байта).
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

// BitReader выполняет побитовое чтение из базового потока.
// Читает данные по байтам и отдает их по одному биту (слева направо, MSB first).
type BitReader struct {
	reader io.Reader
	byte   byte
	count  uint8
	err    error
	buf    [1]byte // буфер: выделяется 1 раз, чтобы не грузить сборщик мусора (GC)
}

// NewBitReader создает новый BitReader для указанного потока ввода.
func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{reader: r}
}

// ReadBit читает и возвращает следующий бит из потока.
func (br *BitReader) ReadBit() (bool, error) {
	// Если на предыдущем шаге произошла ошибка (например, конец файла), сразу возвращаем её
	if br.err != nil {
		return false, br.err
	}

	// Если текущий байт исчерпан, читаем следующий байт из потока
	if br.count == 0 {
		_, br.err = br.reader.Read(br.buf[:])
		if br.err != nil {
			return false, br.err
		}
		br.byte = br.buf[0]
		br.count = 8
	}

	br.count--
	// Извлекаем нужный бит: сдвигаем байт так, чтобы нужный бит оказался в самой правой позиции (младший бит),
	// затем применяем маску & 1, чтобы отбросить все остальные биты.
	// Результат сравниваем с 0, чтобы преобразовать число в bool (true = 1, false = 0).
	bit := (br.byte & (1 << br.count)) != 0
	return bit, nil
}
