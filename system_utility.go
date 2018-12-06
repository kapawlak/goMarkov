package main

import (
	"fmt"
)

////////////////////////////////////// Printing //////////////////////////////////////

//Print Population Grid (For debugging: small populations only!)
func (a *Atom) Print() {
	width:=a.gwid
	length:=a.glen
	fmt.Print("[")
	for i, q := range a.pop {
		fmt.Printf(" %2d", q)
		if (i+1)%width == 0 && (i+1) < width*(length) {
			fmt.Print(" ]\n[")
		}
	}
	fmt.Print(" ]\n")
	fmt.Println("Total population is:", a.tpop)

}
func (sys *Sys) Print() {
	width:=sys.gwid
	length:=sys.glen
	fmt.Print("[")
	for i, q := range sys.sites {
		fmt.Printf(" %2d", sys.den(q))
		if (i+1)%width == 0 && (i+1) < width*(length) {
			fmt.Print(" ]\n[")
		}
	}
	fmt.Print(" ]\n")


}
func (sys *Sys) keychain() {

	k:=sys.tree.bigend.previous
	for k!=sys.tree.smallend{
		fmt.Printf("(%v,%v,%v), ",k.key,k.start,k.end)
		k=k.previous
	}
fmt.Println()

}

	

////////////////////////////////////// Sorting //////////////////////////////////////

func (sys *Sys) Len() int {
	return len(sys.sites)
}
func (sys *Sys) Swap(i, j int) {
	sys.sites[i], sys.sites[j] = sys.sites[j], sys.sites[i]
}
func (sys *Sys) Less(i, j int) bool {
	return sys.den(sys.sites[i]) > sys.den(sys.sites[j])
}