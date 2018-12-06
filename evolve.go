package main

import (
	"fmt"
	"math/rand"
)

////////////////////////////////////// Evolve Structure ////////////////////////////////////
type Evolver struct {
	//parameters
	steps int
	snaps int
	glen  int
	gwid  int

	//bookkeepring
	imagecount int

	//rand key current
	rkey      *Keyz
	rinterval int

	//add and remove utility
	swapslice []int

	//image creation
	imageint int
	pcopy    []int
}

////////////////////////////////////// New evolver ////////////////////////////////////
func Evolve(sys *Sys, steps int, snaps int) {
	sys.init()
	ev := (&Evolver{steps: steps, snaps: snaps, imagecount: 1, glen: sys.glen, gwid: sys.gwid})
	ev.pcopy = make([]int, ev.glen*ev.gwid)
	ev.swapslice = make([]int, 0, ev.glen*ev.gwid)

	if snaps == 0 {
		ev.imageint = steps + 1
	} else if snaps == 1 {
		ev.imageint = steps + 1
		//make first image now
	} else {
		ev.imageint = (steps)/(snaps) + 1
		//make first image now
	}
	sys.evolver = ev

	//pre evolve distribution
	if sys.tree.smallend.next.key == 0 {
		ev.rkey, _ = sys.tree.Get(sys.den(sys.sites[rand.Intn(sys.tree.smallend.next.end+1)]))
	} else {
		ev.rkey, _ = sys.tree.Get(sys.den(sys.sites[rand.Intn(sys.tree.smallend.next.end+1)]))

	}
	ev.rinterval = ev.rkey.end
	fmt.Println(ev.rkey)
	ev.evolve(sys)
}

////////////////////////////////////// Evolution Loop //////////////////////////////////////

func (ev *Evolver) evolve(sys *Sys) {

	step := 1

	for step < ev.steps+1 {
		//fmt.Printf("\r %v   ", step)
		if sys.max() == 0 {
			fmt.Println("\nTotal death at step", step)
			ev.snaps = ev.imagecount
			break
		}
		if step%ev.imageint == 0 {
			ev.imagecount++
			sys.tree.PrintTree("key")
			//copy pop and image
		}

		evolutionstep(sys)
		//sys.tree.PrintTree("key")
		//sys.keychain()
		step++
	}

}

////////////////////////////////////// Evolution Step //////////////////////////////////////

//Main Thread
func evolutionstep(sys *Sys) {

	interval := sys.evolver.rinterval
	randomkey := 0
	for randomkey == 0 {

		sys.evolver.rkey = sys.tree.Select(rand.Intn(sys.tree.Rank(sys.den(sys.sites[rand.Intn(interval+1)])) + 1))
		randomkey = sys.evolver.rkey.key
	}
	//fmt.Println("interval", interval)
	// if sys.tree.smallend.next.key == 0 {
	// 	in := sys.tree.Rank(sys.den(sys.sites[rand.Intn(interval+1)]))
	// 	if in <= 0 {
	// 		sys.tree.PrintTree("key")
	// 		fmt.Println("fuck",in,interval)
	// 		sys.keychain()
	// 	}
	// 	sys.evolver.rkey = sys.tree.Select(rand.Intn(in)+1)
	// } else {
	// 	in := sys.tree.Rank(sys.den(sys.sites[rand.Intn(interval+1)]))
	// 	if in <= 0 {
	// 		fmt.Println("fuck!")
	// 	}
	// 	sys.evolver.rkey = sys.tree.Select(rand.Intn(in))

	// }
	sys.evolver.rinterval = sys.evolver.rkey.end
	//fmt.Println("end", sys.evolver.rinterval)
	test := rand.Intn(sys.evolver.rkey.end + 1)
	//fmt.Println("test", test)
	randomsite := sys.sites[test]
	sum := 0
	for _, p := range sys.pl {
		sum += p.den(randomsite)
		if sum >= sys.evolver.rkey.key {
			p.do(randomsite, p.atom, sys)
			//fmt.Println(o.atom.name)
			return
		}
	}
	//sys.pl[0].do(randomsite, sys.pl[0].atom, sys)
	//sys.pl[0].atom[0].Print()
	fmt.Println("Process Not Found")

}

////////////////////////////////////// Add and Remove //////////////////////////////////////
//changes only to atom populations and maps
//remove(site)
// a.pop--
// send newden=sys.den(site)
// receive swaplist
// do swaps
// return to process

func Remove(site int, a *Atom, sys *Sys) {
	//fmt.Println("Remove")
	ev := sys.evolver
	deletekey := false

	//current density
	dnow := sys.den(site)
	currentposition := sys.smap[site]

	//change pop
	a.pop[site]--
	a.tpop--

	//new density
	dnew := sys.den(site)

	//current key
	keynow := sys.tree.GetKnown(dnow)
	if keynow.start == keynow.end {
		deletekey = true
	}

	////Swap Down
	ev.swapslice = append(ev.swapslice, currentposition, keynow.end)
	//move back starting position
	keynow.end--

	//swap loop
	keynext := keynow.previous
	for keynext.key > dnew {
		//record positions
		ev.swapslice = append(ev.swapslice, keynext.end)
		//move up start and ending
		keynext.start--
		keynext.end--
		keynext = keynext.previous
	}
	//ev.swapslice = append(ev.swapslice, keynext.start)
	//simultaneous slice swap

	for i := 0; i < len(ev.swapslice)-1; i++ {
		sys.smap[sys.sites[ev.swapslice[i]]] = ev.swapslice[i+1]
		sys.sites[ev.swapslice[i]] = sys.sites[ev.swapslice[i+1]]
	}
	sys.smap[sys.sites[ev.swapslice[len(ev.swapslice)-1]]] = currentposition
	sys.sites[ev.swapslice[len(ev.swapslice)-1]] = site

	//delete and add if needed
	if keynext.key != dnew {
		//sys.tree.PrintTree("key")
		//keynext is largest val smaller than target
		//fmt.Println("add", dnew, "from key", keynext.key)
		//fmt.Println(ev.swapslice)
		sys.numkeys++
		//fmt.Println(keynext.key)
		sys.tree.DirInsert(keynext.next, dnew, ev.swapslice[len(ev.swapslice)-1], ev.swapslice[len(ev.swapslice)-1])
		// if keynext!=sys.tree.smallend{
		// sys.tree.DirInsert(keynext, dnew, keynext.start-1, keynext.start-1)
		// }else{
		// 	sys.tree.DirInsert(keynext.next, dnew, keynext.start-1, keynext.start-1)
		// }
	} else {
		keynext.start--
	}

	if deletekey {
		//sys.tree.PrintTree("key")
		sys.numkeys--
		//fmt.Println("delete", keynow.key, dnow)
		sys.tree.DirDelete(keynow)
	}

	//clear swap slice
	ev.swapslice = ev.swapslice[0:0]
}

func Add(site int, a *Atom, sys *Sys) {
	ev := sys.evolver
	deletekey := false

	//current density
	dnow := sys.den(site)
	currentposition := sys.smap[site]

	//change pop
	a.pop[site]++
	a.tpop++

	//new density
	dnew := sys.den(site)

	//current key
	keynow := sys.tree.GetKnown(dnow)
	if keynow.start == keynow.end {
		deletekey = true
	}

	////Swap Up
	ev.swapslice = append(ev.swapslice, currentposition, keynow.start)
	//move back starting position
	keynow.start++

	//swap loop
	keynext := keynow.next
	for keynext.key < dnew {
		//record positions
		ev.swapslice = append(ev.swapslice, keynext.start)
		//move back start and ending
		keynext.start++
		keynext.end++
		keynext = keynext.next
	}
	//ev.swapslice = append(ev.swapslice, keynext.start)
	//simultaneous slice swap
	for i := 0; i < len(ev.swapslice)-1; i++ {
		sys.smap[sys.sites[ev.swapslice[i]]] = ev.swapslice[i+1]
		sys.sites[ev.swapslice[i]] = sys.sites[ev.swapslice[i+1]]
	}
	sys.smap[sys.sites[ev.swapslice[len(ev.swapslice)-1]]] = currentposition
	sys.sites[ev.swapslice[len(ev.swapslice)-1]] = site

	//delete and add if needed

	if keynext.key != dnew {
		//make keynext largest val smaller than targer
		keynext = keynext.previous
		sys.numkeys++
		sys.tree.DirInsert(keynext, dnew, keynext.start-1, keynext.start-1)
	} else {
		keynext.end++
	}
	if deletekey {
		sys.numkeys--
		sys.tree.DirDelete(keynow)
	}

	//clear swap slice
	ev.swapslice = ev.swapslice[0:0]

}

///////////////////////////////////// Tree Manipulation ////////////////////////////////////

//update loop
// random key
// send key to edit
// send density to evolve
// queue++
////edit loop while queue >0, or maybe choose to use select statement
//// swap= append(swap,key.end)
//// if key multiplicity == one ---> start delete
//// if key multiplicity > one ---> change endpoint
//// receive newkey
//// for all key</>newkey
//////change endpoint, swap= append(swap,key.end)
//// if key=newkey change start
//// else add newkey=false
//// send swaplist
//// finish delete
//// add newkey

// func (sys *Sys) editTree(){

// }
