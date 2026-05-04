package core

import "testing"

func TestAssignCodes_MathRules(t *testing.T) {
	// Имитируем данные: 4 символа с длинами: 1, 2, 3, 3
	input := []CodeInfo{
		{Char: 'A', Len: 2},
		{Char: 'B', Len: 1},
		{Char: 'C', Len: 3},
		{Char: 'D', Len: 3},
	}

	assignCodes(input)

	// B: код 0       -> 0b0
	// A: код (0+1)<<1 -> 0b10
	// C: код (10+1)<<1 -> 0b110
	// D: код (110+1)   -> 0b111

	for _, ci := range input {
		switch ci.Char {
		case 'A':
			if ci.Code != 0b10 {
				t.Errorf("A code wrong, expected 0b10, got %b", ci.Code)
			}
		case 'B':
			if ci.Code != 0b0 {
				t.Errorf("B code wrong, expected 0b0, got %b", ci.Code)
			}
		case 'C':
			if ci.Code != 0b110 {
				t.Errorf("C code wrong, expected 0b110, got %b", ci.Code)
			}
		case 'D':
			if ci.Code != 0b111 {
				t.Errorf("D code wrong, expected 0b111, got %b", ci.Code)
			}
		}
	}
}

func TestGetLengths(t *testing.T) {
	// Создаём дерево:
	//      Root
	//     /    \
	//   Leaf(A) Leaf(B)
	root := &Node{
		Left:  &Node{Char: 'A'},
		Right: &Node{Char: 'B'},
	}

	lengths := make(map[byte]uint8)
	getLengths(root, 0, lengths)

	if lengths['A'] != 1 || lengths['B'] != 1 {
		t.Errorf("Expected lengths 1 for A and B, got %d and %d", lengths['A'], lengths['B'])
	}
}

func TestBuildCanonicalCodes_FullPipeline(t *testing.T) {
	// Создаём дерево для одного символа (крайний случай, длина должна стать 1)
	root := &Node{Left: &Node{Char: 'X', Freq: 100}}

	codes := BuildCanonicalCodes(root)
	if len(codes) != 1 || codes[0].Char != 'X' || codes[0].Len != 1 {
		t.Errorf("Failed for single char tree. Got: %v", codes)
	}

	// Создаём дерево из 2 символов
	root2 := &Node{
		Left:  &Node{Char: 'Y'},
		Right: &Node{Char: 'Z'},
	}
	codes2 := BuildCanonicalCodes(root2)

	// Y должен получить код 0, Z код 1.
	for _, c := range codes2 {
		if c.Char == 'Y' && c.Code != 0 {
			t.Errorf("Y should be 0")
		}
		if c.Char == 'Z' && c.Code != 1 {
			t.Errorf("Z should be 1")
		}
	}
}

func TestBuildDecodingTree(t *testing.T) {
	// Иммитируем данные из заголовка (символы и длины)
	input := []CodeInfo{
		{Char: 'A', Len: 2},
		{Char: 'B', Len: 1},
	}

	tree := BuildDecodingTree(input)
	node := tree

	// Идём по коду 'B' (0 -> налево)
	node = node.Left
	if node.Char != 'B' {
		t.Error("Left branch should be 'B'")
	}

	// Возвращаемся в корень и идем по коду 'A' (10 -> вправо, влево)
	node = tree
	node = node.Right
	node = node.Left
	if node.Char != 'A' {
		t.Error("Path Right->Left should be 'A'")
	}
}
