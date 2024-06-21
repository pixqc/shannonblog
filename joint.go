package main

import (
	"fmt"
	"math"
	"math/rand"
)

type JointCharGen struct {
	distribution map[[2]byte]float64
}

func NewJointCharGen(d1, d2 CharGen) JointCharGen {
	jd := make(map[[2]byte]float64)
	for k1, p1 := range d1.distribution {
		for k2, p2 := range d2.distribution {
			var p float64
			joint := [2]byte{k1, k2}
			if k1 == '0' && k2 == 'a' {
				p = p1*p2 + 0.3
			} else {
				p = p1 * p2
			}
			jd[joint] = p
		}
	}
	// normalize
	var sum float64
	for _, p := range jd {
		sum += p
	}
	for k, v := range jd {
		jd[k] = v / sum
	}
	// validate distribution sums to 1
	sum = 0.0
	for _, p := range jd {
		sum += p
	}
	// allow for some floating point error
	if math.Abs(sum-1.0) > 1e-9 {
		panic(fmt.Errorf("probabilities do not sum to 1: got %f", sum))
	}
	return JointCharGen{
		distribution: jd,
	}
}

func (jcg JointCharGen) Sample(count int) []byte {
	result := make([]byte, count*2)
	for i := 0; i < count; i++ {
		r := rand.Float64()
		sum := 0.0
		for k, v := range jcg.distribution {
			sum += v
			if r < sum {
				result[i*2] = k[0]
				result[i*2+1] = k[1]
				break
			}
		}
	}
	return result
}

func calculateEntropy(distribution map[byte]float64) float64 {
	var entropy float64
	for _, p := range distribution {
		entropy += p * math.Log2(p)
	}
	return -entropy
}

func jointEntropy(jcg JointCharGen) float64 {
	var entropy float64
	for _, p := range jcg.distribution {
		entropy += p * math.Log2(p)
	}
	return -entropy
}

func JointRun() {
	// joint entropy: if i join two distributions, what's the entropy? (symmetric)
	// mutual information: how correlated/dependent are the two distributions? (symmetric)
	// conditional entropy: if you know outcome of distribution A, how much does it change B's entropy?
	// if A and B is very correlated, then knowing A will make B almost deterministic
	// reducing entropy of B (asymmetric)
	// randomly-generated distributions
	cg1 := CharGen{distribution: map[byte]float64{
		'0': 0.2119,
		'1': 0.0364,
		'2': 0.2641,
		'3': 0.4876,
	}}
	cg2 := CharGen{distribution: map[byte]float64{
		'a': 0.5577,
		'b': 0.0402,
		'c': 0.2180,
		'd': 0.1841,
	}}
	jcg := NewJointCharGen(cg1, cg2)
	fmt.Println(calculateEntropy(cg1.distribution))
	fmt.Println(calculateEntropy(cg2.distribution))
	fmt.Println(jointEntropy(jcg))
	fmt.Printf("H(cg1) + H(cg2) = %f\n", calculateEntropy(cg1.distribution)+calculateEntropy(cg2.distribution))
	fmt.Printf("H(jcg) = %f\n", jointEntropy(jcg))
	fmt.Println(jcg)
	fmt.Println(string(jcg.Sample(10)))
}
