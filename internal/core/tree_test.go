package core

import "testing"

func TestBuildTree_EdgeCases(t *testing.T) {
	// Передаём пустую мапу частот
	t.Run("Empty frequencies", func(t *testing.T) {
		freq := map[byte]uint64{}
		tree := BuildTree(freq)
		if tree != nil {
			t.Error("Tree should be nil for empty input")
		}
	})

	// Передаем только один символ
	t.Run("One unique character", func(t *testing.T) {
		freq := map[byte]uint64{'X': 100}
		tree := BuildTree(freq)

		// Для одного символа мы делаем искусственный узел,
		// где символ прячется в Left, а корень пустой
		if tree == nil || tree.Left == nil || tree.Left.Char != 'X' {
			t.Error("Tree for 1 char is structured wrong")
		}
		if tree.Right != nil {
			t.Error("Tree for 1 char should not have Right child")
		}
	})
}
