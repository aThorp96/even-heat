package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	. "github.com/athorp96/graphs"
)

// if somehow I need to change how this is checked, this provides an interface

// A rerecombination is a function that somehow constructs a child Hamiltonian
// from two other Hamiltonians.

// Hamiltonian is a hamiltonian cycle.
// it consists of a cycle and a fitness grade
// the lower the fitness, the shorter the cycle

var filepath string
var temperature float64
var coolRate float64
var g *Undirected
var best []int
var bestFit float64

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

func EvenHeat(fp string, initialTemperature float64) []int {
	filepath = fp
	rand.Seed(time.Now().Unix())

	temperature = initialTemperature
	coolRate = 0.000001

	g = NewWeightedGraphFromFile(filepath) //O(n)
	startPath := makeRandomPath()

	endPath := anneal(startPath)

	fmt.Println("Fitness =  ", fitness(endPath.path))
	printPath(endPath.path)

	return endPath.path
}

func printPath(p []int) {
	//fmt.Printf("%v\n", p)
	fmt.Print("Path = [")
	for i := 0; i < len(p)-1; i++ {
		fmt.Printf(" %d,", p[i])
	}
	fmt.Printf(" %d ]\n", p[-1+len(p)])
}

// acceptNewPath allways accepts a better path, and accepts
// a worse path with probability p(d, T) = e^(-d/T)
func acceptNewPath(curPath, newPath []int) bool {
	p1 := fitness(curPath)
	p2 := fitness(newPath)

	if p1 > p2 {
		return true
	} else {
		d := p2 - p1
		p := float64(math.Exp((-1 * d) / temperature))
		//fmt.Printf("%f = e^-(%f-%f) / %f\n", p, p2, p1, temperature)
		randNum := rand.Float64()
		if randNum <= p {
			//fmt.Println("accepted")
		} else {
			//fmt.Println("rejected")
		}
		return randNum <= p
	}
}

func anneal(p *Hamiltonian) *Hamiltonian {
	path := make([]int, len(p.path))
	copy(path, p.path)

	best = make([]int, len(p.path))
	copy(best, p.path)

	bestFit = fitness(best)

	for temperature > 0 {
		path = phase1(path)
		temperature -= coolRate
	}
	copy(p.path, best)

	return p
}

// phase 1 opperates as follows:
// - Compute new state, p' using 2 opting on the portion of
//   the tour given.
// - accept the new path with the acceptNewPath funtion
// - repeat until equalibrium
//
// Currently equalibrium is defined as some number of new path
// rejections in a row, as this would mean there is little to no
// improvement happening.
func phase1(curpath []int) []int {

	numMax := 300
	numAccepts := 0
	numRejects := 0

	newPath := make([]int, len(curpath))

	// count rejections and accepts seperately
	// for anticipated "equalibium" redefinition
	for numRejects+numAccepts < numMax {
		newPath = twoOpt(curpath)

		if acceptNewPath(curpath, newPath) {
			copy(curpath, newPath)
			numAccepts++
		} else {
			numRejects++
		}

		fit := fitness(curpath)
		if fit < bestFit {
			copy(best, curpath)
			bestFit = fit
		}
	}

	return newPath
}

// two-opt performs a 2-opt switch of two elements
// in a path.
func twoOpt(path []int) []int {
	i := rand.Intn(len(path) - 1)
	j := rand.Intn(len(path))

	for j <= i {
		j = rand.Intn(len(path))
	}
	//fmt.Printf("%v\ni: %d\nj: %d\n", path, i, j)
	newPath := []int{}
	newPath = append(newPath, path[:i]...)
	// reverse all elements between i and j
	for h := j - 1; h >= i; h-- {
		newPath = append(newPath, path[h])
		//fmt.Printf("h: %d\npath[h]: %d\n%v\n", h, path[h], newPath)
	}
	newPath = append(newPath, path[j:]...)
	//fmt.Printf("%v\n", newPath)

	return newPath
}

// two-opt performs a 2-opt switch of two elements
// in a path.
func twoOptSwap(path []int) []int {
	i := rand.Intn(len(path) - 1)
	j := (i + 1) % len(path)

	//fmt.Printf("%v\ni: %d\nj: %d\n", path, i, j)
	newPath := make([]int, len(path))
	copy(newPath, path)

	temp := newPath[i]
	newPath[i] = newPath[j]
	newPath[j] = temp

	return newPath
}

// two-opt performs a 2-opt switch of two elements
// in a path.
func twoOptSwitch(path []int) []int {
	i := rand.Intn(len(path) - 1)
	j := rand.Intn(len(path))

	for j == i {
		j = rand.Intn(len(path))
	}
	//fmt.Printf("%v\ni: %d\nj: %d\n", path, i, j)
	newPath := make([]int, len(path))
	copy(newPath, path)

	temp := newPath[i]
	newPath[i] = newPath[j]
	newPath[j] = temp

	return newPath
}

// makeRandomPath creates a random permutation of
// a connected graph's vertices and returns that
// permutation in the form of a *Hamiltonian
func makeRandomPath() *Hamiltonian {
	tour := new(Hamiltonian)
	tour.path = rand.Perm(g.Order())
	tour.fitness = fitness(tour.path)
	return tour
}

// A fitness evaluator
// Returns the sum weight of the walk
func fitness(walk []int) float64 {
	length := 0.0
	for i := 0; i < len(walk); i++ {
		n := (i + 1) % len(walk)
		length += g.Weight(walk[i], walk[n])
	}
	return length
}
