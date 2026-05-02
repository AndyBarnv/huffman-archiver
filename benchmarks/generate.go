package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	os.Mkdir("benchmarks/testdata", 0755)

	// Исходный текст для генерации
	baseText := "Алгоритм Хаффмана — это метод оптимального префиксного кодирования. " +
		"Эффективность сильно зависит от входных данных: на текстах сжатие отличное, а на случайных данных размер только вырастет. "

	sizes := []int{1024, 10 * 1024, 100 * 1024, 1024 * 1024} // 1KB, 10KB, 100KB, 1MB

	// Генерирация текстовых файлов
	for _, size := range sizes {
		repeats := size/len(baseText) + 1
		data := strings.Repeat(baseText, repeats)[:size]
		os.WriteFile(fmt.Sprintf("benchmarks/testdata/text_%dkb.txt", size/1024), []byte(data), 0644)
	}

	// Генерирация случайных бинарных файлов
	for _, size := range sizes {
		buf := make([]byte, size)
		rand.Read(buf)
		os.WriteFile(fmt.Sprintf("benchmarks/testdata/rand_%dkb.bin", size/1024), buf, 0644)
	}

	// Генерация файлов с максимальной энтропией, заполненных байтами 0-255 по кругу (имитация архива)
	for _, size := range sizes {
		buf := make([]byte, size)
		for i := range buf {
			buf[i] = byte(i % 256)
		}
		os.WriteFile(fmt.Sprintf("benchmarks/testdata/pseudo_zip_%dkb.bin", size/1024), buf, 0644)
	}

	fmt.Println("Test data generated in the folder benchmarks/testdata")
}
