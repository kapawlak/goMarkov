package main

import (
	"fmt"
	"math/rand"
)

////////////////////////////////////// Default One-Body Density //////////////////////////////////////
func obdDefault(site int, a []*Atom) int {
	return a[0].pop[site]
}

////////////////////////////////////// Default Two-Body Density //////////////////////////////////////
func tbdDefault(site int, a []*Atom) int {
	return a[0].pop[site] * a[1].pop[site]
}

////////////////////////////////////// Default Hop //////////////////////////////////////
// func SetHopProcess(atomlist interface{}, clist interface{}, processlist *[]*Process) {
// 	newprocesses := setdefaultobProcess("Hop", atomlist, clist, 0, obdDesfault, dohopDefault)
// 	*processlist = append(*processlist, newprocesses[0])
// }
func (sys *Sys) SetHopProcess(atomlist interface{}, clist interface{}) {
	newprocesses := setdefaultobProcess("Hop", atomlist, clist, 0, obdDefault, dohopDefault)
	for _, p := range newprocesses {
		sys.pl = append(sys.pl, *p)
	}
}
func dohopDefault(site int, a []*Atom, sys *Sys) {
	//fmt.Println("hop")
	length := sys.glen
	width := sys.gwid
	//sys.queue++

	w := rand.Intn(4)
	coord := [2]int{(site - site%width) / width, site % width}
	switch w {
	case 0:
		coord[1] = CyclicBoundary(coord[1]+1, width)
	case 1:
		coord[1] = CyclicBoundary(coord[1]-1, width)
	case 2:
		coord[0] = CyclicBoundary(coord[0]+1, length)
	case 3:
		coord[0] = CyclicBoundary(coord[0]-1, length)

	}
	nsite := width*coord[0] + coord[1]
	//fmt.Println("Hop from site ", site, "with", sys.den(site), "to site", nsite, "with", sys.den(nsite))
	// sys.denchan <- t.den(nsite)
	// //fmt.Println("s")
	Remove(site, a[0], sys)
	Add(nsite, a[0], sys)

	// t.rchan <- 0

}

////////////////////////////////////// Default Death //////////////////////////////////////
// func SetDeathProcess(atomlist interface{}, clist interface{}, processlist *[]*Process) {
// 	newprocesses := setdefaultobProcess("Death", atomlist, clist, 0, obdDefault, dodeathDefault)
// 	*processlist = append(*processlist, newprocesses...)
//}
func (sys *Sys) SetDeathProcess(atomlist interface{}, clist interface{}) {
	newprocesses := setdefaultobProcess("Death", atomlist, clist, 0, obdDefault, dodeathDefault)
	for _, p := range newprocesses {
		sys.pl = append(sys.pl, *p)
	}
}
func dodeathDefault(site int, a []*Atom, sys *Sys) {
	//fmt.Println("Death")
	// Remove(site, a[0], t, false)
	// t.rchan <- 0
}

////////////////////////////////////// Default Consumer //////////////////////////////////////
// type interaction interface {
// 	SetConsumerProcess
// 	SetDeathProcess
// 	SetHopProcess
// }

// func SetConsumerProcess(atomlist interface{}, clist interface{}, processlist *[]*Process) {
// 	newprocesses := setdefaulttbProcess("Consume", atomlist, clist, 0, tbdDefault, dodeathDefault)
// 	*processlist = append(*processlist, newprocesses...)
// }
func (sys *Sys) SetConsumerProcess(atomlist interface{}, clist interface{}) {
	newprocesses := setdefaulttbProcess("Consume", atomlist, clist, 0, tbdDefault, dodeathDefault)
	for _, p := range newprocesses {
		sys.pl = append(sys.pl, *p)
	}
}

func doconsumerDefault(site int, a []*Atom, sys *Sys) {
	//fmt.Println("consume site ", site)
	// sys.queue++
	 Remove(site, a[1], sys)
		Add(site, a[0], sys)
	// t.rchan <- 0
}

func null2() {
	fmt.Print("")
}
