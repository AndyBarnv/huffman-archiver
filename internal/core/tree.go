// Package core содержит базовые структуры данных и алгоритмы
// для работы с алгоритмом Хаффмана: построение дерева частот,
// генерация канонических кодов и восстановление дерева для декодирования.
package core

import "container/heap"

// Node представляет узел двоичного дерева Хаффмана.
// Листовые узлы содержат конкретный символ (Char) и его частоту.
// Промежуточные узлы содержат суммарную частоту дочерних поддеревьев, поле Char не используется.
type Node struct {
	Char  byte
	Freq  uint64
	Left  *Node
	Right *Node
}

// PriorityQueue реализует интерфейс container/heap для создания очереди с приоритетом (мин-кучи).
// Приоритет определяется полем Freq: узлы с меньшей частотой находятся "выше" (быстрее извлекаются).
type PriorityQueue []*Node

// Len возвращает количество элементов в очереди (требование интерфейса heap.Interface).
func (pq PriorityQueue) Len() int { return len(pq) }

// Less определяет порядок сортировки в куче (требование интерфейса heap.Interface).
// Возвращает true, если элемент i должен быть "выше" элемента j (то есть его частота меньше).
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

// Swap меняет местами два элемента в очереди (требование интерфейса heap.Interface).
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push добавляет элемент в кучу (требование интерфейса heap.Interface).
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

// Pop извлекает элемент с минимальной частотой из кучи.
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// BuildTree строит дерево Хаффмана на основе карты частот символов.
// Возвращает корень дерева или nil, если карта частот пуста.
func BuildTree(freq map[byte]uint64) *Node {
	pq := make(PriorityQueue, 0, len(freq))
	heap.Init(&pq)

	// Создаем листовые узлы для каждого встретившегося символа и помещаем в кучу
	for char, count := range freq {
		heap.Push(&pq, &Node{Char: char, Freq: count})
	}

	if pq.Len() == 0 {
		return nil
	}

	// Крайний случай: файл состоит из одного уникального символа (например, "AAAA").
	// Чтобы алгоритм мог его закодировать (дать код длиной хотя бы 1 бит),
	// мы искусственно создаем корень и прячем символ в левое поддерево.
	if pq.Len() == 1 {
		node := heap.Pop(&pq).(*Node)
		return &Node{Freq: node.Freq, Left: node}
	}

	// Основной цикл построения дерева:
	// извлекаем два узла с наименьшими частотами, склеиваем их в новый узел
	// и кладем обратно в кучу, пока не останется ровно один узел (корень).
	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)
		parent := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*Node)
}
