package main

import (
	"fmt"
	//"sync"
	"math"
)

//Structures
// KeyTree is a B tree that links nodes to the next largest value
type KeyTree struct {
	//Data dimensions
	glen int
	gwid int

	//Tree root and ends
	root     *Keyz
	bigend   *Keyz
	smallend *Keyz
	deeproot *Keyz

	//Node Count
	count int
}

// Keyz is a node in the tree carrying key, index range, next/prev nodes, left and right children, and a height/balance of the node.
type Keyz struct {
	//values
	key       int
	start     int
	end       int
	height    int8
	lchildren int
	//relationships
	left, right    *Keyz
	next, previous *Keyz
	parent         *Keyz
}

////////////////////////////////////// Initializations //////////////////////////////////////

//TREES///////////////////////////////
//Make New Key Tree
func SysTree(rootkey, rootstart, rootend, length, width int) *KeyTree {
	var kt KeyTree
	kt.init()
	kt.root.key = rootkey
	kt.root.start = rootstart
	kt.root.end = rootend

	kt.count = 1
	kt.glen = length
	kt.gwid = width
	return &kt
}

func (kt *KeyTree) init() {
	kt.root = (&Keyz{key: 0, start: 0, end: 0}).init()

	kt.bigend = (&Keyz{key: int(math.Pow(10, 9)), start: -1, end: -1}).init()
	kt.bigend.next = kt.bigend
	kt.bigend.previous = kt.root
	kt.root.next = kt.bigend

	kt.smallend = (&Keyz{key: -int(math.Pow(10, 9)), start: kt.glen*kt.gwid, end: kt.glen*kt.gwid}).init()
	kt.smallend.previous = kt.smallend
	kt.smallend.next = kt.root
	kt.root.previous = kt.smallend

	kt.deeproot = (&Keyz{key: -999, start: -1, end: -1}).init()
	kt.root.parent = kt.deeproot
	kt.deeproot.right = kt.root

}

//KEYS///////////////////////////////
//Init initializes the values of the node or clears the node and returns the node pointer

func (k *Keyz) init() *Keyz {
	k.height = 1
	k.left = nil
	k.right = nil
	k.lchildren = 0
	return k
}

////////////////////////////////////// Insert and Delete //////////////////////////////////////
// Insert inserts a new key into the tree
func (kt *KeyTree) TopInsert(key, start, end int) {
	added := false
	kt.root, added = kt.root.TopInsert(key, start, end, kt.root.parent)
	if added {
		//deeproot
		kt.deeproot.right = kt.root
		kt.count++
	} else {
		fmt.Println("Not Added")
	}

}
func (k *Keyz) TopInsert(key, start, end int, parent *Keyz) (*Keyz, bool) {
	added := false

	if k == nil { //found place
		//values
		k = (&Keyz{key: key, start: start, end: end}).init()
		k.key = key
		k.start = start
		k.end = end

		//relationships
		k.parent = parent

		if parent.key > key { //inserted to left of parent
			parent.previous.next = k
			k.previous = parent.previous

			k.next = k.parent
			k.next.previous = k
		} else {
			parent.next.previous = k
			k.next = parent.next

			k.previous = k.parent
			k.previous.next = k
		}

		added = true
		return k, added

	} else if k.key < key { //val lesser
		k.right, added = k.right.TopInsert(key, start, end, k)
	} else if k.key > key { //val greater
		k.left, added = k.left.TopInsert(key, start, end, k)
		if added {
			k.lchildren++ //adjust all left paths
		}
	} else { //key already exists, overwrite
		k.start = start
		k.end = end
	}

	if added {

		//adjust height
		k.height = k.maxheight() + 1
		//check balance
		bal := k.balance()
		if bal > 1 {
			if key < k.left.key {
				return k.rotateRight(), added
			} else if key > k.left.key {
				k.left = k.left.rotateLeft()
				return k.rotateRight(), added
			}
		} else if bal < -1 {
			if key > k.right.key {
				return k.rotateLeft(), added
			} else if key < k.right.key {
				k.right = k.right.rotateRight()
				return k.rotateLeft(), added
			}
		}
	}

	return k, added

}
func (kt *KeyTree) DirInsert(k *Keyz, key, start, end int) {
	//fmt.Println("DirInsert ", key, "from", k.key)
	
	kp := k.parent
	knew := k
	onleft := false
	added := false
	if kp.key > k.key {
		onleft = true
		knew, added = kp.left.TopInsert(key, start, end, kp)
	} else {
		knew, added = kp.right.TopInsert(key, start, end, kp)
	}
	if added {
		kt.count++
		kp.iAscend(knew, onleft, key)
		kt.root = kt.deeproot.right
	}
	
}
func (k *Keyz) iAscend(knew *Keyz, onleft bool, key int) {

	kp := k.parent

	if kp == nil {

		k.right = knew
		return
	}
	//	fmt.Println("On key", k.key, "parent", k.parent.key)
	//reset connections
	if onleft {
		k.left = knew
		//if child on left reduce children
		k.lchildren++
	} else {
		k.right = knew
	}
	if kp.key > k.key {
		onleft = true
	} else {
		onleft = false
	}

	//rebalacing
	//adjust height
	k.height = k.maxheight() + 1

	//check balance
	bal := k.balance()
	//	fmt.Println("for", k.key, "bal", bal)
	if bal > 1 {
		if key < k.left.key {
			kp.iAscend(k.rotateRight(), onleft, key)
		} else if key > k.left.key {
			k.left = k.left.rotateLeft()
			kp.iAscend(k.rotateRight(), onleft, key)
		}
	} else if bal < -1 {
		if key > k.right.key {
			kp.iAscend(k.rotateLeft(), onleft, key)
		} else if key < k.right.key {
			k.right = k.right.rotateRight()
			kp.iAscend(k.rotateLeft(), onleft, key)
		}
	} else {
		kp.iAscend(k, onleft, key)
	}

}

//Delete deletes the node from the tree
func (kt *KeyTree) TopDelete(key int) {
	deleted := false
	kt.root, deleted = kt.root.TopDelete(key)
	if deleted {
		kt.count--
	} else {
		fmt.Println("Not Deleted")
	}

}
func (k *Keyz) TopDelete(key int) (*Keyz, bool) {
	deleted := false
	if k == nil { //entry not found
		return k, deleted
	} else if k.key < key { //travel right
		k.right, deleted = k.right.TopDelete(key)
	} else if k.key > key { //travel left
		k.left, deleted = k.left.TopDelete(key)
		if deleted {
			k.lchildren--
		}
	} else {
		k.previous.next = k.next
		k.next.previous = k.previous
		if k.left == nil { //the current node replaced w/ right
			deleted = true

			if k.right != nil {
				k.right.parent = k.parent
			}
			knew := k.right
			k.init() //use pool allocator
			return knew, deleted

		} else if k.right == nil { //the current node replaced w/ left
			deleted = true
			k.left.parent = k.parent
			knew := k.left
			k.init() //use pool allocator
			return knew, deleted

		} else {
			//copy data
			knew := k.right.min()
			k.key, k.start, k.end = knew.key, knew.start, knew.end
			k.previous, k.next = knew.previous, knew.next
			//delete redundant node, rebal
			k.right, deleted = k.right.TopDelete(knew.key)
			//overwrite changes
			k.previous.next = k
			k.next.previous = k

		}

	}

	//rebalance
	if deleted {
		k.height = k.maxheight() + 1
		bal := k.balance()
		if bal > 1 {
			if k.left.balance() >= 0 {
				return k.rotateRight(), deleted
			}
			k.left = k.left.rotateLeft()
			return k.rotateRight(), deleted
		} else if bal < -1 {
			if k.right.balance() <= 0 {
				return k.rotateLeft(), deleted
			}
			k.right = k.right.rotateRight()
			return k.rotateLeft(), deleted
		}
	}
	return k, deleted

}

func (kt *KeyTree) DirDelete(k *Keyz) {
	//	fmt.Println("DirDelete ", k.key)

	kp := k.parent

	onleft := false
	if kp.key > k.key {
		onleft = true
	}
	//do delete
	knew, deleted := k.TopDelete(k.key)

	//ascend tree from current node
	if deleted {
		kp.dAscend(knew, onleft)
		kt.root = kt.deeproot.right

		kt.count--
	}
	//	kt.PrintTree("key")

}
func (k *Keyz) dAscend(knew *Keyz, onleft bool) {
	kp := k.parent
	if kp == nil {
		k.right = knew

	} else {
		//reset connections
		if onleft {
			k.left = knew
			//if child on left reduce children
			k.lchildren--
		} else {
			k.right = knew
		}
		if kp.key > k.key {
			onleft = true
		} else {
			onleft = false
		}

		//rebalacing

		k.height = k.maxheight() + 1
		bal := k.balance()

		if bal > 1 {
			if k.left.balance() >= 0 {
				kp.dAscend(k.rotateRight(), onleft)
			} else {
				k.left = k.left.rotateLeft()
				kp.dAscend(k.rotateRight(), onleft)
			}

		} else if bal < -1 {
			if k.right.balance() <= 0 {
				kp.dAscend(k.rotateLeft(), onleft)
			} else {
				k.right = k.right.rotateRight()
				kp.dAscend(k.rotateLeft(), onleft)
			}

		} else {
			kp.dAscend(k, onleft)
		}

	}

	return
	//recalculate heights and rebalance
	//if parent.key!=-999
	//parent.dAscend(onleft)

}

////////////////////////////////////// Structure Manipulation //////////////////////////////////////

func (k *Keyz) rotateRight() *Keyz {

	//track left children count
	l := k.left

	//fmt.Println("R Swap", k.key, l.key)

	if l != nil {
		k.lchildren -= l.lchildren + 1
	}

	// Rotation
	l.right, k.left = k, l.right

	//set parents
	l.parent = k.parent
	k.parent = l

	if k.left != nil {
		k.left.parent = k
	}

	//l.right.children, n.left.children = n.children, l.right.children

	// update heights
	k.height = k.maxheight() + 1
	l.height = l.maxheight() + 1
	//KT.PrintTree("key")
	return l
}
func (k *Keyz) rotateLeft() *Keyz {

	k.right.lchildren += k.lchildren + 1
	r := k.right
	//fmt.Println("L Swap", k.key, r.key)
	// Rotation
	r.left, k.right = k, r.left

	//set parents
	r.parent = k.parent
	k.parent = r
	if k.right != nil {
		k.right.parent = k
	}

	// update heights
	k.height = k.maxheight() + 1
	r.height = r.maxheight() + 1
	//KT.PrintTree("key")
	return r
}
func (k *Keyz) balance() int8 {
	if k != nil {
		return k.left.getheight() - k.right.getheight()
	}
	return 0
}
func (k *Keyz) min() *Keyz {
	kl := k
	if k.left != nil {
		kl = k.left.min()
	} else {
	}
	return kl
}

////////////////////////////////////// Probing //////////////////////////////////////

// Get returns the node associated with the search value
func (kt *KeyTree) Get(key int) (*Keyz, bool) {
	var node *Keyz
	if kt.root != nil {
		node = kt.root.get(key)
	}
	if node != nil {
		return node, true
	}

	return nil, false
}
func (kt *KeyTree) GetKnown(key int) (k *Keyz) {
	k, f := kt.Get(key)
	if !f {
		fmt.Println("Error:", key, "wasn't found")
		kt.PrintTree("key")
	}
	return k

}
func (n *Keyz) get(key int) *Keyz {
	var node *Keyz
	if key < n.key {
		if n.left != nil {
			node = n.left.get(key)
		}
	} else if key > n.key {
		if n.right != nil {
			node = n.right.get(key)
		}
	} else {
		node = n
	}
	return node
}

//Select returns node w/ i'th largest value
func (kt *KeyTree) Select(rand int) *Keyz {
	return kt.root.Select(rand, kt)
}
func (k *Keyz) Select(rand int, kt *KeyTree) *Keyz {
	l := k.lchildren

	if l == rand {
		return k
	} else if rand < l {
		// if k.left == nil {
		// 	fmt.Println("nil\n")
		// 	kt.PrintTreeInfo()
		// 	// dk, _ := kt.Get(143)
		// 	// fmt.Println("kk")
		// 	// fmt.Println(dk.key, dk.right.key)
		// 	// kt.printTreeFile()
		// }
		return k.left.Select(rand, kt)
	} else {
		// if k.right == nil {
		// 	fmt.Println("nil\n")
		// 	// fmt.Println("kk")
		// 	// kt.PrintKTree()
		// 	// kt.PrintTree("child")
		// 	// dk, _ := kt.Get(30)
		// 	// fmt.Println(dk.start, dk.end)
		// 	// kt.printTreeFile()
		// }
		return k.right.Select(rand-l-1, kt)

	}

}

func (kt *KeyTree) Rank(key int) int {
	//fmt.Println("Rank of", key)
	//fmt.Println("rank of ", key, kt.bigend.previous.key)
	return kt.root.Rank(key)
}

func (k *Keyz) Rank(key int) int {
	if k.key == key {
		return k.lchildren
	} else if k.key > key {
		return k.left.Rank(key)
	} else {
		return k.lchildren + 1 + k.right.Rank(key)
	}

}
