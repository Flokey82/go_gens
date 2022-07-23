package genetics

import (
	"math/rand"
)

type Genes uint64

func (g *Genes) setBits(offs, n int, val uint64) {
	ge := uint64(*g)
	ge &^= ((1 << n) - 1) // clear bits
	ge |= (val << offs)
	*g = Genes(ge)
}

func (g *Genes) getBits(offs, n int) uint64 {
	return (uint64(*g) >> offs) & ((1 << n) - 1)
}

func Mix(a, b Genes) Genes {
	t := uint64(a) & uint64(b)        // common bits
	x := (uint64(a) | uint64(b)) &^ t // diff bits
	t |= x & uint64(rand.Int63())     // random diffs

	// Mutations (random bit flips)
	nMutations := 1
	for i := 0; i < nMutations; i++ {
		t ^= 1 << rand.Intn(63)
	}

	return Genes(t)
}

// Gene layout
//  _______________________ 2 gender
// || _____________________ 2 eye color
// ||||  __________________ 3 hair color         ___________________________ 4 Openness
// |||| ||| _______________ 4 complexion        ||||  ______________________ 4 Conscientiousness
// |||| |||| ||| __________ 3 height            |||| |||| __________________ 4 Extraversion
// |||| |||| |||| || ______ 3 mass              |||| |||| ||||  ____________ 4 Agreeableness
// |||| |||| |||| |||| | __ 3 growth            |||| |||| |||| ||||  _______ 4 Neuroticism
// |||| |||| |||| |||| ||||                     |||| |||| |||| |||| ||||
// xxxx xxxx|xxxx xxxx|xxxx xxxx|xxxx xxxx|xxxx xxxx|xxxx xxxx|xxxx xxxx|xxxx xxxx
//                          |||| |||| |||| ||||                          |||| ||||
// 4 strength _________________  |||| |||| ||||                           ________ unused
// 4 intelligence __________________  |||| ||||
// 4 dexterity __________________________  ||||
// 4 resilience ______________________________
