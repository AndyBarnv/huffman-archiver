package core

import "sort"

// FileSignature — число, записываемое в начало архива для валидации формата.
const FileSignature = "HUF"

// CodeInfo хранит информацию о каноническом коде конкретного символа.
type CodeInfo struct {
	Char byte   // Символ (от 0 до 255)
	Len  uint8  // Длина кода
	Code uint32 // Код
}

// getLengths выполняет обход дерева (DFS) и собирает длины кодов для каждого символа.
// Результат записывается в мапу lengths, где ключ — символ, значение — длина его кода.
func getLengths(node *Node, length uint8, lengths map[byte]uint8) {
	if node == nil {
		return
	}

	// Найдя лист записываем длину символа
	if node.Left == nil && node.Right == nil {
		// Крайний случай: файл состоит из одного уникального символа.
		// В BuildTree мы сделали искусственный корень, поэтому длина до листа равна 0.
		// Хаффман не может иметь коды длиной 0 бит, поэтому принудительно задаем длину 1.
		if length == 0 {
			length = 1
		}
		lengths[node.Char] = length
		return
	}

	// Рекурсивно спускаемся влево (добавляем 1 бит к длине) и вправо
	getLengths(node.Left, length+1, lengths)
	getLengths(node.Right, length+1, lengths)
}

// assignCodes реализует алгоритм генерации значений канонических кодов по их длинам.
// Функция модифицирует переданный слайс codes, заполняя поле Code.
func assignCodes(codes []CodeInfo) {
	if len(codes) == 0 {
		return
	}

	// Сортировка кодов по правилам канонического Хаффмана:
	// Первичный ключ: длина кода (по возрастанию).
	// Вторичный ключ: значение символа (по алфавиту).
	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Len != codes[j].Len {
			return codes[i].Len < codes[j].Len
		}
		return codes[i].Char < codes[j].Char
	})

	// Математическое присвоение значений кодов
	code := uint32(0)
	codes[0].Code = code // Первому символу выделяется код из нулей
	for i := 1; i < len(codes); i++ {
		code++ // Увеличиваем код на 1

		// Если следующий символ имеет большую длину кода, освобождается место
		// для дополнительных бит. Это делается побитовым сдвигом влево.
		if codes[i].Len > codes[i-1].Len {
			code <<= (codes[i].Len - codes[i-1].Len)
		}
		codes[i].Code = code
	}
}

// BuildCanonicalCodes принимает построенное дерево и возвращает слайс канонических кодов
// для всех уникальных символов, встретившихся в данных.
func BuildCanonicalCodes(root *Node) []CodeInfo {
	lengthsMap := make(map[byte]uint8)
	getLengths(root, 0, lengthsMap)

	if len(lengthsMap) == 0 {
		return nil
	}

	// Подготавливаем слайс для передачи в функцию assignCodes
	codes := make([]CodeInfo, 0, len(lengthsMap))
	for char, length := range lengthsMap {
		codes = append(codes, CodeInfo{Char: char, Len: length})
	}

	assignCodes(codes)

	return codes
}

// BuildDecodingTree восстанавливает дерево Хаффмана из таблицы длин канонических кодов.
// Дерево строится специально для быстрого декодирования: каждый бит считанного потока
// — это переход налево (0) или направо (1), пока не будет достигнут лист с символом.
func BuildDecodingTree(codes []CodeInfo) *Node {
	if len(codes) == 0 {
		return nil
	}

	// Восстанавливаем значения кодов из длин
	assignCodes(codes)

	// Строим дерево по восстановленным кодам
	root := &Node{}
	for _, ci := range codes {
		node := root
		// Идем по битам кода от старшего к младшему
		for i := uint8(0); i < ci.Len; i++ {
			// Извлекаем текущий бит
			bit := (ci.Code >> (ci.Len - 1 - i)) & 1
			if bit == 0 {
				if node.Left == nil {
					node.Left = &Node{}
				}
				node = node.Left
			} else {
				if node.Right == nil {
					node.Right = &Node{}
				}
				node = node.Right
			}
		}
		// Конец кода — это лист. Записываем в него символ
		node.Char = ci.Char
	}
	return root
}
