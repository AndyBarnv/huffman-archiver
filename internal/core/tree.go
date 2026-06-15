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
func (queue PriorityQueue) Len() int { return len(queue) }

// Less определяет порядок сортировки в куче (требование интерфейса heap.Interface).
// Возвращает true, если элемент i должен быть "выше" элемента j (то есть его частота меньше).
func (queue PriorityQueue) Less(i, j int) bool {
	return queue[i].Freq < queue[j].Freq
}

// Swap меняет местами два элемента в очереди (требование интерфейса heap.Interface).
func (queue PriorityQueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

// Push добавляет элемент в кучу (требование интерфейса heap.Interface).
func (queue *PriorityQueue) Push(x interface{}) {
	*queue = append(*queue, x.(*Node))
}

// Pop извлекает элемент с минимальной частотой из кучи.
func (queue *PriorityQueue) Pop() interface{} {
	old := *queue
	num := len(old)
	item := old[num-1]
	*queue = old[:num-1]
	return item
}

// BuildTree строит дерево Хаффмана на основе карты частот символов.
// Возвращает корень дерева или nil, если карта частот пуста.
func BuildTree(freq map[byte]uint64) *Node {
	queue := make(PriorityQueue, 0, len(freq))
	heap.Init(&queue)

	// Создаем листовые узлы для каждого встретившегося символа и помещаем в кучу
	for char, count := range freq {
		heap.Push(&queue, &Node{Char: char, Freq: count})
	}

	if queue.Len() == 0 {
		return nil
	}

	// Крайний случай: файл состоит из одного уникального символа (например, "AAAA").
	// Чтобы алгоритм мог его закодировать (дать код длиной хотя бы 1 бит),
	// мы искусственно создаем корень и прячем символ в левое поддерево.
	if queue.Len() == 1 {
		node := heap.Pop(&queue).(*Node)
		return &Node{Freq: node.Freq, Left: node}
	}

	// Основной цикл построения дерева:
	// извлекаем два узла с наименьшими частотами, склеиваем их в новый узел
	// и кладем обратно в кучу, пока не останется ровно один узел (корень).
	for queue.Len() > 1 {
		left := heap.Pop(&queue).(*Node)
		right := heap.Pop(&queue).(*Node)
		parent := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&queue, parent)
	}

	return heap.Pop(&queue).(*Node)
}
