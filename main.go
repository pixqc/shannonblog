package main

import (
	"fmt"
	"math/rand"
)

type CharGenNode struct {
	Char  byte
	Left  *CharGenNode
	Right *CharGenNode
}

func NewCharGenNode(typ string) *CharGenNode {
	switch typ {
	// ```
	//       root
	//      /    \
	//     0      1
	//    / \    / \
	//   0   1  0   1
	//  / \ / \/ \ / \
	// a  bc  de  fg  h
	// ```
	case "3bit":
		root := &CharGenNode{}
		root.Left = &CharGenNode{}
		root.Left.Left = &CharGenNode{}
		root.Left.Left.Left = &CharGenNode{Char: 'a'}
		root.Left.Left.Right = &CharGenNode{Char: 'b'}
		root.Left.Right = &CharGenNode{}
		root.Left.Right.Left = &CharGenNode{Char: 'c'}
		root.Left.Right.Right = &CharGenNode{Char: 'd'}
		root.Right = &CharGenNode{}
		root.Right.Left = &CharGenNode{}
		root.Right.Left.Left = &CharGenNode{Char: 'e'}
		root.Right.Left.Right = &CharGenNode{Char: 'f'}
		root.Right.Right = &CharGenNode{}
		root.Right.Right.Left = &CharGenNode{Char: 'g'}
		root.Right.Right.Right = &CharGenNode{Char: 'h'}
		return root
	case "3bit-skewed":
		// this kind of tree is usually constructed
		// by Huffman coding, this is manually constructed
		// to make things deterministic
		// ```
		//   root
		//  /    \
		// a      1
		//       / \
		//      b   1
		//         / \
		//        c   1
		//           / \
		//          d   1
		//             / \
		//            e   1
		//               / \
		//              f   1
		//                 / \
		//                g   h
		// ```
		root := &CharGenNode{}
		root.Left = &CharGenNode{Char: 'a'}
		root.Right = &CharGenNode{}
		node := root.Right
		for _, char := range []byte{'b', 'c', 'd', 'e', 'f'} {
			node.Left = &CharGenNode{Char: char}
			node.Right = &CharGenNode{}
			node = node.Right
		}
		node.Left = &CharGenNode{Char: 'g'}
		node.Right = &CharGenNode{Char: 'h'}
		return root
	default:
		return nil
	}
}

func (root *CharGenNode) MakeCharMap() map[byte]string {
	charMap := make(map[byte]string)
	var traverse func(node *CharGenNode, path string)
	traverse = func(node *CharGenNode, path string) {
		if node == nil {
			return
		}
		if node.Char != 0 {
			charMap[byte(node.Char)] = path
		}
		traverse(node.Left, path+"0")
		traverse(node.Right, path+"1")
	}
	traverse(root, "")
	return charMap
}

func (root *CharGenNode) MakeBinaryMap() map[string]byte {
	charMap := root.MakeCharMap()
	binaryMap := make(map[string]byte)
	for char, code := range charMap {
		binaryMap[code] = char
	}
	return binaryMap
}

func (root *CharGenNode) Encode(input []byte) string {
	charMap := root.MakeCharMap()
	res := ""
	for _, char := range input {
		res += charMap[char]
	}
	return res
}

func (root *CharGenNode) Decode(input string) []byte {
	binaryMap := root.MakeBinaryMap()
	res := make([]byte, 0)
	for i := 0; i < len(input); {
		for j := i + 1; j <= len(input); j++ {
			if char, ok := binaryMap[input[i:j]]; ok {
				res = append(res, char)
				i = j
				break
			}
		}
	}
	return res
}

type CharGen struct {
	distribution map[byte]float64
}

func NewCharGen(typ string) CharGen {
	charMap := make(map[byte]float64)
	switch typ {
	case "3bit":
		for i := 'a'; i <= 'h'; i++ {
			charMap[byte(i)] = 0.125
		}
	case "3bit-skewed":
		charMap['a'] = 0.5
		charMap['b'] = 0.25
		charMap['c'] = 0.125
		charMap['d'] = 0.0625
		charMap['e'] = 0.03125
		charMap['f'] = 0.015625
		charMap['g'] = 0.0078125
		charMap['h'] = 0.0078125
	default:
		panic(fmt.Errorf("unknown distribution type: %s", typ))
	}
	return CharGen{
		distribution: charMap,
	}
}

func (cg CharGen) Sample(count int) []byte {
	tmp := make([]byte, count)
	for i := 0; i < count; i++ {
		r := rand.Float64()
		sum := 0.0
		for k, v := range cg.distribution {
			if r >= sum && r < sum+v {
				tmp[i] = k
				break
			}
			sum += v
		}
	}
	return tmp
}

func main() {
	fmt.Println("Hello, World!")
	cg1 := NewCharGen("3bit")
	samples := cg1.Sample(10)
	tree := NewCharGenNode("3bit")
	fmt.Println(samples)
	fmt.Println(tree.Encode(samples))
	fmt.Println(tree.Decode(tree.Encode(samples)))
}
