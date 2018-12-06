package main

import (
	"fmt"
	"math"
	"time"
	"github.com/pkg/profile"
	"runtime"
)

func main() {
	_ = runtime.SetBlockProfileRate
	defer profile.Start(profile.ProfilePath(".")).Stop()

	steps, length, width := int(math.Pow(10, 5)), 1000, 1000
	//Initialize a system
	sys := NewSys(length, width, 1)

	//Set Atoms
	gridb := RandomPop(0, 5, length, width)
	gridf := Circular(0.5,0.25,0.3,10,length,width)

	b := DefineAtom(gridb, "bacteria", length, width)
	f := DefineAtom(gridf, "food", length, width)

	sys.SetHopProcess(b, 1.0)
	//sys.SetHopProcess(f, 1.0)

	sys.SetConsumerProcess([]*Atom{b,f},0.1)
	fmt.Println("At Start")
	//b.Print()
	start := time.Now()
	Evolve(sys, steps, 10)
	elapsed := time.Since(start)
	fmt.Println("\nreg time:", elapsed, "over", steps, "steps", "so av time per step:", float64(elapsed/1000000000)/float64(steps))
	//sys.tree.PrintTree("child")

	fmt.Println("Done")
	//b.Print()
}
