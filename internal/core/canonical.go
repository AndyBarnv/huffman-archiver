package core

import "sort"

const FileSignature = "HUF"

type CodeInfo struct {
	Char byte
	Len  uint8
	Code uint32
}

func getLengths(node *Node, length uint8, lengths map[byte]uint8) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		if length == 0 {
			length = 1
		}
		lengths[node.Char] = length
		return
	}

	getLengths(node.Left, length+1, lengths)
	getLengths(node.Right, length+1, lengths)
}

func assignCodes(codes []CodeInfo) {
	if len(codes) == 0 {
		return
	}

	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Len != codes[j].Len {
			return codes[i].Len < codes[j].Len
		}
		return codes[i].Char < codes[j].Char
	})

	code := uint32(0)
	codes[0].Code = code
	for i := 1; i < len(codes); i++ {
		code++
		if codes[i].Len > codes[i-1].Len {
			code <<= (codes[i].Len - codes[i-1].Len)
		}
		codes[i].Code = code
	}
}

func BuildCanonicalCodes(root *Node) []CodeInfo {
	lengthsMap := make(map[byte]uint8)
	getLengths(root, 0, lengthsMap)

	if len(lengthsMap) == 0 {
		return nil
	}

	codes := make([]CodeInfo, 0, len(lengthsMap))
	for char, length := range lengthsMap {
		codes = append(codes, CodeInfo{Char: char, Len: length})
	}

	assignCodes(codes)

	return codes
}

func BuildDecodingTree(codes []CodeInfo) *Node {
	if len(codes) == 0 {
		return nil
	}

	assignCodes(codes)

	root := &Node{}
	for _, ci := range codes {
		node := root
		for i := uint8(0); i < ci.Len; i++ {
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
		node.Char = ci.Char
	}
	return root
}
