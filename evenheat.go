package main

import (
	"math/rand"

	. "github.com/athorp96/graphs"
)

// if somehow I need to change how this is checked, this provides an interface

// A rerecombination is a function that somehow constructs a child Hamiltonian
// from two other Hamiltonians.

// Hamiltonian is a hamiltonian cycle.
// it consists of a cycle and a fitness grade
// the lower the fitness, the shorter the cycle

var path string
var temperature float64

// Hamiltonian is a hamiltonian cycle.
// it consists of a cycle and a fitness grade
// the lower the fitness, the shorter the cycle
type Hamiltonian struct {
	// the path
	path []int

	// The optimitality of the solution
	// The lower the number the shorter the path
	fitness float64
}

func EvenHeat(filepath string, initialTemperature,  float64) []int {
	path = filepath

	mutationRate = mutRate
	graph := NewWeightedGraphFromFile(filepath) //O(n)
	startPath = makeRandomPath(graph)

	endPath = anneal(startPath, coolRate)

	return int[1, 2]
}

func anneal(p * Hamiltonian, coolRate float64) * Hamiltonian {

}

func makeRandomPath(g *Undirected) *Hamiltonian {
	tour := new(Hamiltonian)
	tour.path = rand.Perm(g.Order())
	tour.fitness = fitness(g, tour.path)
	return tour
}

// A fitness evaluator
// Returns the sum weight of the walk
func fitness(g *Undirected, walk []int) float64 {
	length := 0.0
	for i := 0; i < len(walk); i++ {
		n := (i + 1) % len(walk)
		length += g.Weight(walk[i], walk[n])
	}
	return length
}
