// Package genetics provides a simple framework for gene mixing / editing for procedural generation.
package genetics

import "math/rand"

// Pool provides means to use a set seed for gene generation/mixing.
type Pool struct {
	Rand *rand.Rand
}

// NewPool initializes a gene pool with the given seed.
func NewPool(seed int64) *Pool {
	return &Pool{
		Rand: rand.New(rand.NewSource(seed)),
	}
}

// NewRandom returns a randomized set of genes.
func (p *Pool) NewRandom() Genes {
	return Genes(p.Rand.Uint64())
}

// Mix combines the two sets of genes by using the bits in common in both sets
// as baseline and then randomly sets bits that are different in both sets.
// Additional a number of mutations can be applied which randomly flip bits.
func (p *Pool) Mix(a, b Genes, nMutations int) Genes {
	t := uint64(a) & uint64(b)        // common bits
	x := (uint64(a) | uint64(b)) &^ t // diff bits
	t |= x & uint64(p.Rand.Int63())   // random diffs

	// Mutations (random bit flips)
	for i := 0; i < nMutations; i++ {
		t ^= 1 << p.Rand.Intn(63)
	}

	return Genes(t)
}

type Genes uint64

// NewRandom returns a randomized set of genes.
func NewRandom() Genes {
	return Genes(rand.Uint64())
}

// Set the gene 'd' to the supplied value.
func (g *Genes) Set(d Gene, val int) {
	g.setBits(d.Offset, d.NumBits, uint64(val))
}

// Get the value of gene 'd'.
func (g *Genes) Get(d Gene) int {
	return int(g.getBits(d.Offset, d.NumBits))
}

func (g *Genes) setBits(offs, n int, val uint64) {
	mask := uint64(1<<n) - 1
	ge := uint64(*g)
	ge &^= mask << offs        // clear bits
	ge |= (val & mask) << offs // set masked bits
	*g = Genes(ge)
}

func (g *Genes) getBits(offs, n int) uint64 {
	return (uint64(*g) >> offs) & ((1 << n) - 1)
}

// Mix combines the two sets of genes by using the bits in common in both sets
// as baseline and then randomly sets bits that are different in both sets.
// Additional a number of mutations can be applied which randomly flip bits.
func Mix(a, b Genes, nMutations int) Genes {
	t := uint64(a) & uint64(b)        // common bits
	x := (uint64(a) | uint64(b)) &^ t // diff bits
	t |= x & uint64(rand.Int63())     // random diffs

	// Mutations (random bit flips)
	for i := 0; i < nMutations; i++ {
		t ^= 1 << rand.Intn(63)
	}

	return Genes(t)
}

// Example gene layout
//
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

type Gene struct {
	NumBits int
	Offset  int
}

func (g *Gene) MaxVal() int {
	return (1 << g.NumBits) - 1
}
