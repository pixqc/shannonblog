package main

import (
	"container/heap"
	"fmt"
	"math/rand"
)

// SECTION: chargens

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

func (cg CharGen) Join(cg2 CharGen) CharGen {
	// independently join two distributions
	jd := make(map[string]float64)
	for k1, p1 := range cg.distribution {
		for k2, p2 := range cg2.distribution {
			joint := k1 + k2
			p := p1 * p2
			jd[joint] = p
		}
	}
	newCg := CharGen{distribution: jd}
	newCg.Normalize()
	return newCg
}

func (cg CharGen) Correlate(k string, correlationFactor float64) CharGen {
	// add dependence to a distribution
	nd := make(map[string]float64)
	for k1, p1 := range cg.distribution {
		nd[k1+k] = p1 * correlationFactor
	}
	cg.distribution = nd
	cg.Normalize()
	return cg
}

// SECTION: trees and huffman trees

type Node struct {
	Value     string
	Count     int
	Left      *Node
	Right     *Node
	CharMap   map[string]string
	BinaryMap map[string]string
}

func NewTree(typ string) *Node {
	switch typ {
	case "3bit":
		// graphviz
		// digraph BinaryTree {
		//     node [shape=square, fixedsize=true, width=0.4];
		//     root [label="root"];
		//     zero1 [label="0"];
		//     one1 [label="1"];
		//     zero2 [label="0"];
		//     one2 [label="1"];
		//     zero3 [label="0"];
		//     one3 [label="1"];
		//     a [label="a"];
		//     b [label="b"];
		//     c [label="c"];
		//     d [label="d"];
		//     e [label="e"];
		//     f [label="f"];
		//     g [label="g"];
		//     h [label="h"];
		//     root -> zero1;
		//     root -> one1;
		//     zero1 -> zero2;
		//     zero1 -> one2;
		//     one1 -> zero3;
		//     one1 -> one3;
		//     zero2 -> a;
		//     zero2 -> b;
		//     one2 -> c;
		//     one2 -> d;
		//     zero3 -> e;
		//     zero3 -> f;
		//     one3 -> g;
		//     one3 -> h;
		//     {rank=same; a; b; c; d; e; f; g; h;}
		// }
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
		root.CharMap = root.MakeStrMap()
		root.BinaryMap = root.MakeBinaryMap()
		return root
	case "3bit-skewed":
		// graphviz:
		// digraph Tree {
		//   node [shape=square, fixedsize=true, width=0.4];
		//   root [label="root"];
		//     a [label="a"];
		//     b [label="b"];
		//     c [label="c"];
		//     d [label="d"];
		//     e [label="e"];
		//     f [label="f"];
		//     g [label="g"];
		//     h [label="h"];
		//     one1 [label="1"];
		//     one2 [label="1"];
		//     one3 [label="1"];
		//     one4 [label="1"];
		//     one5 [label="1"];
		//     one6 [label="1"];
		//     root -> a;
		//     root -> one1;
		//     one1 -> b;
		//     one1 -> one2;
		//     one2 -> c;
		//     one2 -> one3;
		//     one3 -> d;
		//     one3 -> one4;
		//     one4 -> e;
		//     one4 -> one5;
		//     one5 -> f;
		//     one5 -> one6;
		//     one6 -> g;
		//     one6 -> h;
		// }
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
		root.CharMap = root.MakeStrMap()
		root.BinaryMap = root.MakeBinaryMap()
		return root
	default:
		return nil
	}
}

func (root *Node) MakeStrMap() map[string]string {
	// eg. {"a": "000", "b": "001", ...}
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
	// eg. {"000": "a", "001": "b", ...}
	charMap := root.MakeStrMap()
	binaryMap := make(map[string]string)
	for char, code := range charMap {
		binaryMap[code] = char
	}
	return binaryMap
}

func (root *Node) Encode(input string) string {
	result := ""
	for _, char := range input {
		result += root.CharMap[string(char)]
	}
	return result
}

func (root *Node) EncodeChunks(input []string) string {
	result := ""
	for _, char := range input {
		result += root.CharMap[char]
	}
	return result
}

func (root *Node) Decode(binary string) string {
	result := ""
	for i := 0; i < len(binary); {
		node := root
		for node.Value == "" {
			if string(binary[i]) == "0" {
				node = node.Left
			} else {
				node = node.Right
			}
			i++
		}
		result += node.Value
	}
	return result
}

type MinHeap []Node

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Node))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func NewHuffmanTree(input string, chunkSize int) Node {
	// 1. make frequency map
	frequencyMap := make(map[string]int)
	for i := 0; i < len(input); i += chunkSize {
		end := i + chunkSize
		if end > len(input) {
			end = len(input)
		}
		chunk := input[i:end]
		frequencyMap[chunk]++
	}
	// 2. turn map to min heap
	h := &MinHeap{}
	heap.Init(h)
	for k, v := range frequencyMap {
		heap.Push(h, Node{k, v, nil, nil, nil, nil})
	}
	// 3. build huffman tree
	for h.Len() > 1 {
		left := heap.Pop(h).(Node)
		right := heap.Pop(h).(Node)
		heap.Push(h, Node{"", left.Count + right.Count, &left, &right, nil, nil})
	}
	node := h.Pop().(Node)
	node.CharMap = node.MakeStrMap()
	node.BinaryMap = node.MakeBinaryMap()
	return node
}

func splitString(src string) []string {
	var res []string
	for i := 0; i < len(src); i += 2 {
		end := i + 2
		if end > len(src) {
			end = len(src)
		}
		res = append(res, src[i:end])
	}
	return res
}

func main() {
	// TODO: fact check the numbers
	// as sampleSize -> Inf, the calculated values will asymptotically
	// approach entropy, all of the values below are approximations
	sampleSize := 10000
	cg1 := NewCharGen("3bit")
	src1 := cg1.Sample(sampleSize)
	cg2 := NewCharGen("3bit-skewed")
	src2 := cg2.Sample(sampleSize)
	fmt.Printf("src1: %s...\n", src1[:50])
	fmt.Printf("src2: %s...\n", src2[:50])
	fmt.Println()
	// t1 is tree optimized to encode cg1, t2 for cg2
	// t2(src2) will produce shorter binary strings
	t1 := NewTree("3bit")
	t2 := NewTree("3bit-skewed")
	encodedSrc1T1 := t1.Encode(src1)
	encodedSrc2T2 := t2.Encode(src2)
	encodedSrc2T1 := t1.Encode(src2)
	encodedSrc1T2 := t2.Encode(src1)
	h1 := float64(len(encodedSrc1T1)) / float64(sampleSize)
	h2 := float64(len(encodedSrc2T2)) / float64(sampleSize)
	h12 := float64(len(encodedSrc2T1)) / float64(sampleSize)
	h21 := float64(len(encodedSrc1T2)) / float64(sampleSize)
	printEncodingInfo := func(label string, encoded string, h float64) {
		fmt.Printf("%s: %s...\n", label, encoded[:40])
		fmt.Printf("len(%s): %d\n", label, int(h*float64(sampleSize)))
		fmt.Printf("avg bit length per char (%s): %f\n", label, h)
		fmt.Println()
	}
	printEncodingInfo("t1.Encode(src1)", encodedSrc1T1, h1)
	printEncodingInfo("t2.Encode(src2)", encodedSrc2T2, h2)
	printEncodingInfo("t1(src2)", encodedSrc2T1, h12)
	printEncodingInfo("t2(src1)", encodedSrc1T2, h21)
	// KL divergence: how much extra bits required
	// KL(cg1||cg2): extra bits when using t1 (optimized for cg1) to encode src2 (from cg2)
	// KL(cg2||cg1): extra bits when using t2 (optimized for cg2) to encode src1 (from cg1)
	fmt.Printf("approximated KL(cg1||cg2): %f\n", h12-h2)
	fmt.Printf("approximated KL(cg2||cg1): %f\n", h21-h1)
	fmt.Println()
	fmt.Println()
	// joint entropy: if we join two distributions, what's the entropy
	cg3 := cg1.Join(NewCharGen("nums"))
	src3 := splitString(cg3.Sample(sampleSize))
	t3 := NewHuffmanTree(cg3.Sample(sampleSize), 2)
	fmt.Printf("t3.Encode(src3): %s...\n", t3.EncodeChunks(src3)[:50])
	// H(cg3) = H(cg1) + H(cg2) if cg1 and cg2 are independent
	// now lets add dependence/correlation to cg3

	// Commented out as it wasn't being used
	// h3 := float64(len(t3.Encode(cg3.Sample(sampleSize))) / sampleSize)
	// fmt.Println(t3.Encode(src3))
}
