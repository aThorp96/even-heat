package main

import (
	"flag"
	"math/rand"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	rand.Seed(time.Now().Unix())

}
