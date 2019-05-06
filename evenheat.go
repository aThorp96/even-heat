package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	. "github.com/athorp96/graphs"
)

// starting parameters
var filepath string
var temperature float64
var coolRate float64
var numProcessors int
var chunksize int

// algorithm globals
var g *Undirected
var best []int
var bestFit float64

const ready int = 1
const done int = 2

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

// EvenHeat handles some of the overhead for the algorithm
// make the graph, plot the energy over time, printing, etc
func EvenHeat(fp string, initialTemperature float64) []int {
	filepath = fp
	rand.Seed(time.Now().Unix())

	temperature = initialTemperature
	coolRate = 0.0001
	numProcessors = 3
	chunksize = 8

	g = NewWeightedGraphFromFile(filepath) //O(n)
	startPath := makeRandomPath()

	c := make(chan float64)

	go energyTracker(c)
	endPath := anneal(startPath, c)

	fmt.Println("# Fitness =  ", fitness(endPath.path))
	//printPath(endPath.path)

	return endPath.path
}

func energyTracker(c chan float64) {
	for e := range c {
		fmt.Printf("# %f\n", e)
	}
	fmt.Println("# Done")
}

func anneal(p *Hamiltonian, energyChan chan float64) *Hamiltonian {
	path := make([]int, len(p.path))
	copy(path, p.path)

	best = make([]int, len(p.path))
	copy(best, p.path)
	bestFit = fitness(best)

	// create processors
	toProcessors := make(chan int)
	fromProcessors := make(chan int)
	n := len(path)
	for temperature > 0 {
		for i := 0; i < numProcessors; i++ {
			j := i * chunksize
			top := path[n-j-chunksize : n-j]
			bot := path[j : j+chunksize]
			go process(top, bot, toProcessors, fromProcessors)
		}

		// wait for all processes to finish
		for i := 0; i < numProcessors; i++ {
			_ = <-fromProcessors
		}

		energyChan <- fitness(path)
		if fitness(path) < bestFit {
			bestFit = fitness(path)
			copy(best, path)
		}
		path = quarterTurn(path)
		temperature -= coolRate
	}

	copy(p.path, best)
	close(energyChan)

	return p
}

func quarterTurn(p []int) []int {

	quarter := p[0:2] // take first two elements off an end
	rest := p[2:]     // and get the path without those elements
	turned := make([]int, len(rest))
	copy(turned, rest)
	turned = append(turned, quarter...)

	return turned
}

// a process performs the algorithm specified on a "top tier"
// path and "bottom tier" path.
//
// process are designed to run independantly. They perform
// their work on a portion of an array, and are managed by
// the anneal function.
//
// - perform swapping on the bottom tier
// - perform swapping on the top tier
// - permorm remote swapping between the tiers
// - return the results and where each tier starts
func process(top, bot []int, rx, tx chan int) {

	botChan := make(chan []int)
	topChan := make(chan []int)

	// run concurrent top and bottom
	// tier work concurrently
	go localMin(bot, botChan)
	go localMin(top, topChan)

	for i := 0; i < 2; i++ {
		select {
		case _ = <-botChan:
		case _ = <-topChan:
		}
	}
	// perform remote swapping
	remoteSwap(top, bot)
	tx <- done
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
func localMin(curpath []int, retChan chan []int) {

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
	}
	retChan <- newPath
}

// Remote swap opperates as follows:
// - swap two elements between the top and bottom tiers
// - accept the swap with the acceptNewPath function
// - repeat until equalibrium
//
// Currently equalibrium is defined as some number of new path
// rejections in a row, as this would mean there is little to no
// improvement happening.
func remoteSwap(top, bot []int) ([]int, []int) {
	numMax := 300
	numAccepts := 0
	numRejects := 0

	// count rejections and accepts seperately
	// for anticipated "equalibium" redefinition
	for numRejects+numAccepts < numMax {
		newTop, newBot := twoOptRemote(top, bot)

		if acceptNewRemotePath(top, bot, newTop, newBot) {
			copy(bot, newBot)
			copy(top, newTop)
			numAccepts++
		} else {
			numRejects++
		}
	}
	return top, bot
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
// in a path. **swaps elements i and j
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
// in a path. **swaps elements i and j
func twoOptRemote(top, bot []int) ([]int, []int) {
	newBot := make([]int, len(bot))
	newTop := make([]int, len(top))
	copy(newBot, bot)
	copy(newTop, top)

	i := rand.Intn(len(top))
	j := rand.Intn(len(bot))

	temp := newTop[i]
	newTop[i] = newBot[j]
	newBot[j] = temp

	return newTop, newBot
}

// two-opt performs a 2-opt reversal between two elements
// in a path. **reverses elements between i and j
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

// acceptNewPath always accepts a better path, and accepts
// a worse path with probability p(d, T) = e^(-d/T)
func acceptNewPath(curPath, newPath []int) bool {

	accepted := true

	p1 := fitness(curPath)
	p2 := fitness(newPath)

	if p1 < p2 {
		d := p2 - p1
		p := float64(math.Exp((-1 * d) / temperature))
		randNum := rand.Float64()
		accepted = randNum <= p
		fmt.Println("0")
	} else {
		fmt.Println("1")
	}

	return accepted
}

// acceptNewRemotePath always accepts a better path, and accepts
// a worse path with probability p(d, T) = e^(-d/T)
// it compares the difference of two paths and their respective
// new paths to determine acceptance
func acceptNewRemotePath(top, bot, newTop, newBot []int) bool {

	accepted := true

	t1 := fitness(top)
	t2 := fitness(newTop)
	b1 := fitness(bot)
	b2 := fitness(newBot)

	if t1 < t2 && b1 > b2 {
		//fmt.Printf("Worse path ")
		d := t2 - t1
		p := float64(math.Exp((-1 * d) / temperature))
		randNum := rand.Float64()
		accepted = randNum <= p
	} else if t1 > t2 && b1 < b2 {
		//fmt.Printf("Worse path ")
		d := b2 - b1
		p := float64(math.Exp((-1 * d) / temperature))
		randNum := rand.Float64()
		accepted = randNum <= p
	} else if t1 < t2 && b1 < b2 {
		d := t2 - t1 + b2 - b1
		d = d / 2
		p := float64(math.Exp((-1 * d) / temperature))
		randNum := rand.Float64()
		accepted = randNum <= p
		//fmt.Println("0")
	} else {
		//fmt.Println("1")
	}

	return accepted
}

func printPath(p []int) {
	//fmt.Printf("%v\n", p)
	fmt.Print("# Path = [")
	for i := 0; i < len(p)-1; i++ {
		fmt.Printf(" %d,", p[i])
	}
	fmt.Printf(" %d ]\n", p[-1+len(p)])
}
