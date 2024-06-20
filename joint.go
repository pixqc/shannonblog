package main

import (
	"fmt"
	"math"
)

type JointCharGen struct {
	distribution map[[2]byte]float64
}

func NewJointCharGen() JointCharGen {
	// randomly-generated distributions
	d1 := map[byte]float64{
		0: 0.2119,
		1: 0.0364,
		2: 0.2641,
		3: 0.4876,
	}
	d2 := map[byte]float64{
		'a': 0.5577,
		'b': 0.0402,
		'c': 0.2180,
		'd': 0.1841,
	}
	jd := make(map[[2]byte]float64)
	for k1, p1 := range d1 {
		for k2, p2 := range d2 {
			joint := [2]byte{k1, k2}
			p := p1 * p2
			jd[joint] = p
		}
	}
	// validate distribution sums to 1
	var sum float64
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

func JointRun() {
	// joint probability: if i join two distributions, what's the entropy?
	// mutual information: how correlated/dependent are the two distributions?
	// conditional entropy: how much information is in one distribution given the other?
	cg := NewJointCharGen()
	fmt.Println(cg)
}
