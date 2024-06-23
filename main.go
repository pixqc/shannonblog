package main

import (
	"fmt"
	"math/rand"
)

type CharGen struct {
	distribution map[string]float64
}

func NewCharGen(typ string) CharGen {
	charMap := make(map[string]float64)
	switch typ {
	case "3bit":
		for i := 'a'; i <= 'h'; i++ {
			charMap[string(i)] = 0.125
		}
	case "3bit-skewed":
		charMap["a"] = 0.5
		charMap["b"] = 0.25
		charMap["c"] = 0.125
		charMap["d"] = 0.0625
		charMap["e"] = 0.03125
		charMap["f"] = 0.01562
		charMap["g"] = 0.0078125
		charMap["h"] = 0.0078125
	case "nums":
		// random probs
		charMap["0"] = 0.2119
		charMap["1"] = 0.0364
		charMap["2"] = 0.2641
		charMap["3"] = 0.4876
	default:
		panic(fmt.Sprintf("unknown distribution type: %s", typ))
	}
	cg := CharGen{
		distribution: charMap,
	}
	cg.Normalize()
	return cg
}

func (cg *CharGen) Normalize() {
	sum := 0.0
	for _, v := range cg.distribution {
		sum += v
	}
	for k, v := range cg.distribution {
		cg.distribution[k] = v / sum
	}
}

func (cg CharGen) Sample(count int) string {
	var result string
	for i := 0; i < count; i++ {
		r := rand.Float64()
		sum := 0.0
		for k, v := range cg.distribution {
			if r >= sum && r < sum+v {
				result += k
				break
			}
			sum += v
		}
	}
	return result
}

func JoinCharGens(d1, d2 CharGen) CharGen {
	jd := make(map[string]float64)
	for k1, p1 := range d1.distribution {
		for k2, p2 := range d2.distribution {
			joint := k1 + k2
			p := p1 * p2
			jd[joint] = p
		}
	}
	cg := CharGen{distribution: jd}
	cg.Normalize()
	return cg
}

type Node struct {
	Value     string
	Left      *Node
	Right     *Node
	CharMap   map[string]string
	BinaryMap map[string]string
}

func NewTree(typ string) *Node {
	switch typ {
	case "3bit":
		// ```
		//       root
		//      /    \
		//     0      1
		//    / \    / \
		//   0   1  0   1
		//  / \ / \/ \ / \
		// a  bc  de  fg  h
		// ```
		root := &Node{}
		root.Left = &Node{}
		root.Left.Left = &Node{}
		root.Left.Left.Left = &Node{Value: "a"}
		root.Left.Left.Right = &Node{Value: "b"}
		root.Left.Right = &Node{}
		root.Left.Right.Left = &Node{Value: "c"}
		root.Left.Right.Right = &Node{Value: "d"}
		root.Right = &Node{}
		root.Right.Left = &Node{}
		root.Right.Left.Left = &Node{Value: "e"}
		root.Right.Left.Right = &Node{Value: "f"}
		root.Right.Right = &Node{}
		root.Right.Right.Left = &Node{Value: "g"}
		root.Right.Right.Right = &Node{Value: "h"}
		root.CharMap = root.MakeCharMap()
		root.BinaryMap = root.MakeBinaryMap()
		return root
	case "3bit-skewed":
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
		root := &Node{}
		root.Left = &Node{Value: "a"}
		root.Right = &Node{}
		node := root.Right
		for _, value := range []string{"b", "c", "d", "e", "f"} {
			node.Left = &Node{Value: value}
			node.Right = &Node{}
			node = node.Right
		}
		node.Left = &Node{Value: "g"}
		node.Right = &Node{Value: "h"}
		root.CharMap = root.MakeCharMap()
		root.BinaryMap = root.MakeBinaryMap()
		return root
	default:
		return nil
	}
}

func (root *Node) MakeCharMap() map[string]string {
	charMap := make(map[string]string)
	var traverse func(node *Node, path string)
	traverse = func(node *Node, path string) {
		if node == nil {
			return
		}
		if node.Value != "" {
			charMap[node.Value] = path
		}
		traverse(node.Left, path+"0")
		traverse(node.Right, path+"1")
	}
	traverse(root, "")
	return charMap
}

func (root *Node) MakeBinaryMap() map[string]string {
	charMap := root.MakeCharMap()
	binaryMap := make(map[string]string)
	for char, code := range charMap {
		binaryMap[code] = char
	}
	return binaryMap
}

func main() {
	cg1 := NewCharGen("3bit")
	fmt.Println(cg1.Sample(10))
	t1 := NewTree("3bit")
	fmt.Println(t1.CharMap)
}
