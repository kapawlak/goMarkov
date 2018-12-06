package main

import (
	"fmt"
	"sort"
)

////////////////////////////////////// Atoms //////////////////////////////////////
//Atom structure
type Atom struct {
	glen int
	gwid int

	pop  []int
	tpop int
	name string
}

//Internal Set
func (a *Atom) set(grid []int, name string, dim ...int) {
	//change this to be general!
	if dim != nil {
		a.glen = dim[0]
		a.gwid = dim[1]
	}
	a.pop = make([]int, len(grid))
	copy(a.pop, grid)
	for _, q := range a.pop {
		a.tpop += q
	}
	a.name = name
}

//External Set
func DefineAtom(grid []int, name string, dim ...int) *Atom {
	//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Allow for list of inputs
	var x Atom
	x.set(grid, name, dim...)
	return &x
}

////////////////////////////////////// Processes //////////////////////////////////////
//Process structure
type Process struct {
	name string                   // name of process
	c    float64                  // coefficient
	atom []*Atom                  // list of atoms involved
	dist int                      // spatial range of process
	res  float64                  // resolution for calculating density. Inherited from Sys
	df   func(int, []*Atom) int   // density function for density method
	do   func(int, []*Atom, *Sys) // Implimentation of process

}

//Internal set
func (p *Process) set(name string, c float64, dist int, atomlist []*Atom, df func(int, []*Atom) int, do func(int, []*Atom, *Sys)) {
	//Sets Process Parameters
	p.name = name
	p.c = c
	p.atom = atomlist

	p.dist = dist

	p.df = df
	p.do = do

}

//Quick Set Default Processes
func setdefaultobProcess(name string, atomlist interface{}, clist interface{}, dist int, df func(int, []*Atom) int, do func(int, []*Atom, *Sys)) []*Process {
	xL := make([]*Process, 0)
	switch a := atomlist.(type) {
	case *Atom:
		var x Process
		switch c := clist.(type) {
		case float64:
			x.set(name, c, dist, []*Atom{a}, df, do)
		case []float64:
			fmt.Println("Warining: Coefficent list provided but only one Atom given. Using first entry")
			co := clist.([]float64)[0]
			x.set(name, co, dist, []*Atom{a}, df, do)
		default:
			fmt.Println("Error, invalid Coef format, not set")
		}
		xL = append(xL, &x)
	case []*Atom:
		x := make([]Process, len(a))
		switch c := clist.(type) {
		case float64:
			for i := 0; i < len(a); i++ {
				x[i].set(name, c, dist, []*Atom{a[i]}, df, do)
				xL = append(xL, &x[i])
			}
			fmt.Println("Warning: Single Coefficent Provided. Using for all Processes in list")
		case []float64:
			for i := 0; i < len(a); i++ {
				x[i].set(name, c[i], dist, []*Atom{a[i]}, df, do)
				xL = append(xL, &x[i])
			}
		default:
			fmt.Println("Error, invalid Coef format")
		}
	default:
		fmt.Println("Error, invalid format")
	}
	return xL

}
func setdefaulttbProcess(name string, atomlist interface{}, clist interface{}, dist int, df func(int, []*Atom) int, do func(int, []*Atom, *Sys)) []*Process {
	xL := make([]*Process, 0)
	switch a := atomlist.(type) {
	case []*Atom:
		var x Process
		switch c := clist.(type) {
		case float64:
			x.set(name, c, dist, a, df, do)
		case []float64:
			fmt.Println("Warining: Coefficent list provided but only one Atom given. Using first entry")
			co := clist.([]float64)[0]
			x.set(name, co, dist, a, df, do)
		default:
			fmt.Println("Error, invalid Coef format")
		}
		xL = append(xL, &x)
	case [][]*Atom:
		x := make([]Process, len(a))
		switch c := clist.(type) {
		case float64:
			for i := 0; i < len(a); i++ {
				x[i].set(name, c, dist, a[i], df, do)
				xL = append(xL, &x[i])
			}
			fmt.Println("Warning: Single Coefficent Provided. Using for all Processes in list")
		case []float64:
			for i := 0; i < len(a); i++ {
				x[i].set(name, c[i], dist, a[i], df, do)
				xL = append(xL, &x[i])
			}
		default:
			fmt.Println("Error, invalid Coef format")
		}
	default:
		fmt.Println("Error, invalid format")
	}
	return xL

}

//Custom Set with pointer return
func SetProcess(atomlist interface{}, clist interface{}, dist int, namelist interface{}, dflist interface{}, dolist interface{}) *Process {
	var x Process
	x.set(namelist.(string), clist.(float64), dist, atomlist.([]*Atom), dflist.(func(int, []*Atom) int), dolist.(func(int, []*Atom, *Sys)))
	return &x

}

//Methods
//Density Method
func (p *Process) den(site int) int {
	return int(p.c/p.res) * p.df(site, p.atom)
}

////////////////////////////////////// System //////////////////////////////////////
type Sys struct {
	//Data dimension
	glen int
	gwid int
	res  float64

	//relevant process list
	pl []Process

	//Tracker maps are manually manipulated for efficieny
	smap    []int
	sites   []int
	numkeys int

	//Tree
	tree *KeyTree

	//Evolver
	evolver *Evolver
}

//Initialize system
func NewSys(length, width int, Resolution float64) *Sys {
	var sys Sys
	sys.glen = length
	sys.gwid = width
	sys.res = Resolution
	sys.pl = make([]Process, 0)
	return &sys
}

var KT *KeyTree

func (sys *Sys) init() {
	length := sys.glen
	width := sys.gwid
	//CopyProcesses
	for i := range sys.pl {
		sys.pl[i].res = sys.res
	}

	//Create System Grid and Tracking
	sys.sites = make([]int, sys.glen*sys.gwid)
	sys.smap = make([]int, sys.glen*sys.gwid)
	for i := range sys.sites {
		sys.sites[i] = i
	}
	sort.Sort(sys)
	for i, s := range sys.sites {
		sys.smap[s] = i
	}

	////Generate Tree
	sys.tree = SysTree(sys.den(sys.sites[0]), 0, 0, length, width)
	KT = sys.tree

	//obtain density keys
	keys := make([][3]int, 1, length*width)
	keys[0] = [3]int{sys.den(sys.sites[0]), 0, 0}
	for i := 1; i < len(sys.sites); i++ {
		keys[len(keys)-1][2] = i - 1
		if sys.den(sys.sites[i]) != sys.den(sys.sites[i-1]) {
			keys = append(keys, [3]int{sys.den(sys.sites[i]), i, i})
		}
		if i == len(sys.sites)-1 {
			keys[len(keys)-1][2] = i
		}

	}
	fmt.Println(keys)

	//set inital number of keys
	sys.numkeys = len(keys)

	//insert all into tree
	kl := len(keys) - 1
	for i := range keys {
		sys.tree.TopInsert(keys[kl-i][0], keys[kl-i][1], keys[kl-i][2])
	}

	sys.tree.Check()

}

//Density Method for System
func (sys *Sys) den(site int) int {
	sum := 0
	for _, p := range sys.pl {
		sum += p.den(site)
	}
	return sum
}

//Max density method
func (sys *Sys) max() int {
	return sys.tree.bigend.previous.key
}
