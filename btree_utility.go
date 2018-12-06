package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

////////////////////////////////////// Basic //////////////////////////////////////

// Len return the number of nodes in the tree
func (kt *KeyTree) Count() int {
	return kt.count
}
func (k *Keyz) getheight() int8 {
	if k != nil {
		return k.height
	}
	return 0
}
func (k *Keyz) maxheight() int8 {
	hl := k.left.getheight()
	hr := k.right.getheight()
	if k != nil {
		if hl >= hr {
			return hl
		} else {
			return hr
		}
	}
	return 0
}

////////////////////////////////////// Testing

func (kt *KeyTree) Check() {
	//ascending in-order traversal
	key := kt.smallend.next
	for key != kt.bigend {
		//in order?
		if key.key >= key.next.key {
			fmt.Println("Next out of order on", key.key)
			log.Fatal()
		}
		if key.key <= key.previous.key {
			fmt.Println("Previous out of order on", key.key)
			log.Fatal()
		}
		//parents agree?
		kl := key.left
		kr := key.right
		if kl != nil && kl.parent != key {
			fmt.Println(key.key, "key's left has wrong parent")
			log.Fatal()
		}
		if kr != nil && kr.parent != key {
			fmt.Println(key.key, "key's right has wrong parent")
			log.Fatal()
		}

		key = key.next
	}
	fmt.Println("Checked")

}

////////////////////////////////////// Printing //////////////////////////////////////

//Printing Utils
//////////////////////////////////////

var DASH string
var LEND string
var REND string
var PIPE string

//Print to terminal
func (kt *KeyTree) PrintTree(data string) {
	DASH, LEND, REND, PIPE = "-", "|", "|","|"
	//DASH,LEND,REND, PIPE="─","┌","┐","│"
	//height of tree *2 to account for "|"
	height := int((kt.root.height))
	//make a [][]interface{} to hold strings and tree data
	treearray := make([][]string, 2*height+1)

	//calculate the number of maximum elements in the last row
	width := int(math.Pow(2, float64(kt.root.height-1))) + 2

	//calculate spacing:
	space := 1
	if data == "data" {
		space = len(strconv.Itoa(kt.glen * kt.gwid))
	} else {
		space = len(strconv.Itoa(kt.bigend.previous.key))
	}
	spacing := "%" + strconv.Itoa(space) + "v"

	//make an array with width double the max elements + the height to allow for padding
	for i := range treearray {
		treearray[i] = make([]string, height+2*width+2)
		for j := range treearray[i] {
			//the array is initialized with only space
			treearray[i][j] = " "
		}
	}
	//fmt.Println("dim", len(treearray), len(treearray[0]))
	//CoordPrint take a pointer to the array and overwrites the spaces where appropriate
	kt.root.CoordPrint(data, 0, height/2+width+1, 0, &treearray, (height+width)-1, space)
	//Print the array element by element
	for _, ta := range treearray {
		for _, r := range ta {
			//I set to %3v because I have 3 digit numbers. If this is changed to e.g. %2v
			//you also have to change --- to -- etc in CoordPrint. I could automate but
			//I'm lazy
			fmt.Printf(spacing, r)
		}
		fmt.Println()
	}

}

//Print to .txt file
func (kt *KeyTree) printTreeFile(data string, name string) error {

	f, err := os.Create(name + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()

	height := int(2 * (kt.root.height))
	treearray := make([][]string, height+1)
	width := int(math.Pow(2, float64(kt.root.height-1)))

	space := 1
	if data == "data" {
		space = len(strconv.Itoa(kt.glen * kt.gwid))
	} else if data == "child" {
		space = len(strconv.Itoa(kt.root.lchildren))

	} else {
		space = len(strconv.Itoa(kt.bigend.previous.key))
	}
	spacing := "%" + strconv.Itoa(space) + "v"

	for i := range treearray {
		treearray[i] = make([]string, height+2*width)
		for j := range treearray[i] {
			//the array is initialized with only space
			treearray[i][j] = strings.Repeat(" ", space)
		}
	}

	kt.root.CoordPrint(data, 0, (height+2*width)/2+1, 0, &treearray, (height+2*width)/2, space)

	for _, ta := range treearray {
		for _, r := range ta {
			//I set to %3v because I have 3 digit numbers. If this is changed to e.g. %2v
			//you also have to change --- to -- etc in CoordPrint. I could automate but
			//I'm lazy
			fmt.Fprintf(f, spacing, r)
		}
		fmt.Fprintln(f, " ")
	}

	return nil
}

func (k *Keyz) CoordPrint(data string, y, x, px int, treearray *[][]string, width, space int) {
	//Coord print takes the Node, data, the y level, x position, parent's x position,
	//array pointer, and the current spacing
	//fmt.Println(x, y, len((*treearray)[0]), len(*treearray))
	//check that you aren't accessing a nil node
	if k != nil {
		//print the desired data
		switch data {
		case "key":
			(*treearray)[y][x] = strconv.Itoa(k.key)
		case "data":
			(*treearray)[y][x] = strconv.Itoa(k.key)
			(*treearray)[y+1][x] = strconv.Itoa(k.start)
			(*treearray)[y+2][x] = strconv.Itoa(k.end)
		case "child":
			(*treearray)[y][x] = strconv.Itoa(k.lchildren)
		}
		//for all children nodes, put characters above.
		if y > 0 {
			for len((*treearray)[y-2][px]) < space {
				(*treearray)[y-2][px] = DASH + (*treearray)[y-2][px]
			}
			//if this is a left child, put left chars
			if px > x {
				(*treearray)[y-2][x] = LEND //+ strings.Repeat(DASH, space-1)
				(*treearray)[y-1][x] = PIPE
				for i := x + 1; i < px; i++ {
					(*treearray)[y-2][i] = strings.Repeat(DASH, space)
				}
			} else { //else put right chars
				(*treearray)[y-2][x] = strings.Repeat(DASH, space-1) + REND
				(*treearray)[y-1][x] = PIPE
				for i := px + 1; i < x; i++ {
					(*treearray)[y-2][i] = strings.Repeat(DASH, space)
				}

			}
		}
		//lazy way to avoid index error
		if y+2 < len(*treearray) {
			//reduce width by a factor of 2
			w := width / 2
			//pass to left and right nodes with new coordinates
			k.left.CoordPrint(data, y+2, x-w, x, treearray, w, space)
			k.right.CoordPrint(data, y+2, x+w, x, treearray, w, space)
		}
	}

}

//Print keys+data to terminal
func (kt *KeyTree) PrintTreeInfo() {
	fmt.Println()
	key := kt.bigend.previous

	for key != kt.smallend {
		fmt.Printf("%2v  ", key.key)
		key = key.previous
	}
	fmt.Println()
	key = kt.bigend.previous

	for key != kt.smallend {
		fmt.Printf("%2v  ", key.start)
		key = key.previous
	}
	fmt.Println()
	key = kt.bigend.previous

	for key != kt.smallend {
		fmt.Printf("%2v  ", key.end)
		key = key.previous
	}
	fmt.Println()

}

///junk
//Print to .txt file (implement as go routine)
// func (kt *KeyTree) goprintTreeFile() error {

// 	f, err := os.Create("TreeNotDefered.txt")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	for ch := range kt.writetree {
// 		if ch[0] == 0 {
// 			fmt.Fprintf(f, "!!Added Key %v\n\n", ch[1])
// 		} else if ch[0] == 1 {
// 			fmt.Fprintf(f, "!!Deleted Key %v\n\n", ch[1])

// 		}

// 		height := int(2 * (kt.root.height))
// 		treearray := make([][]string, height+1)
// 		width := int(math.Pow(2, float64(kt.root.height-1)))
// 		for i := range treearray {
// 			treearray[i] = make([]string, height+2*width)
// 			for j := range treearray[i] {
// 				//the array is initialized with only space
// 				treearray[i][j] = " "
// 			}
// 		}
// 		space := len(strconv.Itoa(kt.bigend.previous.key))
// 		spacing := "%" + strconv.Itoa(space) + "v"
// 		kt.root.CoordPrint("key", 0, (height+2*width)/2+1, 0, &treearray, (height+2*width)/2, space)
// 		for _, ta := range treearray {
// 			for _, r := range ta {
// 				//I set to %3v because I have 3 digit numbers. If this is changed to e.g. %2v
// 				//you also have to change --- to -- etc in CoordPrint. I could automate but
// 				//I'm lazy
// 				fmt.Fprintf(f, spacing, r)
// 			}
// 			fmt.Fprintln(f, " ")
// 		}
// 		fmt.Fprintln(f, " ")
// 		fmt.Fprintln(f, " ")
// 		for i := range treearray {
// 			for j := range treearray[i] {
// 				//re-initialize
// 				treearray[i][j] = " "
// 			}
// 		}
// 		space = len(strconv.Itoa(kt.root.children))
// 		spacing = "%" + strconv.Itoa(space) + "v"
// 		kt.root.CoordPrint("child", 0, (height+2*width)/2+1, 0, &treearray, (height+2*width)/2, space)
// 		for _, ta := range treearray {
// 			for _, r := range ta {
// 				//I set to %3v because I have 3 digit numbers. If this is changed to e.g. %2v
// 				//you also have to change --- to -- etc in CoordPrint. I could automate but
// 				//I'm lazy
// 				fmt.Fprintf(f, spacing, r)
// 			}
// 			fmt.Fprintln(f, " ")
// 		}
// 		fmt.Fprintln(f, " ")
// 		fmt.Fprintln(f, " ")
// 		fmt.Fprintln(f, " ")
// 		fmt.Fprintln(f, " ")
// 		fmt.Fprintln(f, " ------------------------------------------------------------------------------ ")
// 		kt.wg.Done()

// 	}

// 	return nil
// }
