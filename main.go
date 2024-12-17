package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

const MAX_KEYS = 4
const MIN_KEYS = 2

type BNode struct {
	keys     []int
	children []*BNode
	isLeaf   bool
}

func (node BNode) string(indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sNode:\n", indent))
	builder.WriteString(fmt.Sprintf("%s  keys: %v\n", indent, node.keys))
	builder.WriteString(fmt.Sprintf("%s  is leaf: %v\n", indent, node.isLeaf))

	if len(node.children) > 0 {
		builder.WriteString(fmt.Sprintf("%s  children:\n", indent))
		for _, child := range node.children {
			builder.WriteString(child.string(indent + "  "))
		}
	}

	return builder.String()
}

func (node BNode) String() string {
	return node.string("")
}

type BTree struct {
	root   *BNode
	runner *Runner
}

func newBTree(runner *Runner) *BTree {
	return &BTree{
		root: &BNode{
			keys:     []int{},
			children: nil,
			isLeaf:   true,
		},
		runner: runner,
	}
}

func (tree *BTree) insert(key int) {
	parent, childIndex, _ := tree.search(key)
	if childIndex >= len(parent.children) {
		parent.keys = append(parent.keys, key)
		slices.Sort(parent.keys)
		if len(parent.keys) > MAX_KEYS {
			tree.splitNode(parent, -1)
		}
	} else {
		if tree.runner.debug {
			fmt.Printf("key %d child index %d children:%v\n", key, childIndex, parent.children)
		}
		curr := parent.children[childIndex]
		curr.keys = append(curr.keys, key)
		slices.Sort(curr.keys)
		if len(curr.keys) > MAX_KEYS {
			tree.splitNode(parent, childIndex)
		}
	}
}

func (tree *BTree) search(key int) (*BNode, int, int) {
	if tree.runner.debug {
		fmt.Println("start find")
	}
	prev := tree.root
	curr := tree.root

	childIndex := 0
	keyIndex := -1
	for {
		i := 0
		if tree.runner.debug {
			fmt.Printf("curr at start %v\n", curr)
		}
		for i < len(curr.keys) && key > curr.keys[i] {
			i++
		}

		if i < len(curr.keys) && key == curr.keys[i] {
			keyIndex = i
			break
		}

		if i < len(curr.children) {
			prev = curr
			childIndex = i
			curr = curr.children[i]
			if tree.runner.debug {
				fmt.Printf("i %d curr at end %v\n", i, curr)
			}
		} else {
			break
		}
	}

	return prev, childIndex, keyIndex
}

func (tree *BTree) delete(key int) {
	// TODO: Implement deletion logic
	parent, childIndex, keyIndex := tree.search(key)
	if keyIndex == -1 {
		return
	}

	// Root
	if childIndex >= len(parent.children) {

		return
	}

	node := parent.children[childIndex]
	if node.isLeaf {
		if len(node.keys) > MIN_KEYS {
			j := keyIndex + 1
			if j > len(node.keys) {
				j = len(node.keys)
			}
			node.keys = slices.Delete(node.keys, keyIndex, j)
			slices.Sort(node.keys)
			return
		}

		leftSibIndex := childIndex - 1
		rigthSibIndex := childIndex + 1

		if leftSibIndex > -1 && len(parent.children[leftSibIndex].keys) > MIN_KEYS {
			leftSib := parent.children[leftSibIndex]
			if len(leftSib.keys) == MIN_KEYS {
				tree.mergeNodes(parent, childIndex)
				return
			}

			parentKeyIndex := leftSibIndex
			temp := parent.keys[parentKeyIndex]
			parent.keys[parentKeyIndex] = leftSib.keys[len(leftSib.keys)-1]
			leftSib.keys = leftSib.keys[:len(leftSib.keys)-1]
			node.keys[keyIndex] = temp
			slices.Sort(node.keys)
			return
		}

		if rigthSibIndex < len(parent.children) && len(parent.children[rigthSibIndex].keys) > MIN_KEYS {
			rigthSib := parent.children[rigthSibIndex]
			if len(rigthSib.keys) == MIN_KEYS {
				tree.mergeNodes(parent, childIndex)
				return
			}

			parentKeyIndex := childIndex
			temp := parent.keys[parentKeyIndex]
			parent.keys[parentKeyIndex] = rigthSib.keys[0]
			rigthSib.keys = rigthSib.keys[1:]
			node.keys[keyIndex] = temp
			slices.Sort(node.keys)
			return
		}

		panic("node must have either left or rigth sibiling")
	}

}

func (tree *BTree) splitNode(parent *BNode, childIndex int) {
	var node *BNode
	if childIndex == -1 {
		node = parent
	} else {
		node = parent.children[childIndex]
	}

	midIndex := len(node.keys) / 2
	middleKey := node.keys[midIndex]

	left := &BNode{
		keys:     append([]int(nil), node.keys[:midIndex]...),
		isLeaf:   node.isLeaf,
		children: nil,
	}

	right := &BNode{
		keys:     append([]int(nil), node.keys[midIndex+1:]...),
		isLeaf:   node.isLeaf,
		children: nil,
	}

	if !node.isLeaf {
		left.children = append([]*BNode(nil), node.children[:midIndex+1]...)
		right.children = append([]*BNode(nil), node.children[midIndex+1:]...)
	}

	if tree.runner.debug {
		fmt.Printf("left %v rigth %v\n", left, right)
	}

	if childIndex == -1 {
		node.keys = []int{middleKey}
		node.children = []*BNode{left, right}
		node.isLeaf = false
	} else {
		parent.keys = append(parent.keys, middleKey)
		slices.Sort(parent.keys)

		parent.children[childIndex] = left
		parent.children = append(parent.children[:childIndex+1], append([]*BNode{right}, parent.children[childIndex+1:]...)...)
	}

	if len(parent.keys) > MAX_KEYS {
		tree.splitNode(parent, -1)
	}
}

func (tree *BTree) mergeNodes(parent *BNode, childIndex int) {
	// TODO: Implement node merging logic
}

func (tree *BTree) traverse(node *BNode) {
	// TODO: Implement in-order traversal
}

func (tree *BTree) printTree() {
	q := []*BNode{tree.root}
	var curr *BNode

	i := -1
	for len(q) != 0 {
		i++
		curr, q = shift(q)
		fmt.Printf("%v ", curr)

		for _, child := range curr.children {
			q = append(q, child)
		}

		if i == 5 || i == 0 {
			fmt.Println()
			i = 0
		}
	}
}

func shift[E any](arr []E) (E, []E) {
	if len(arr) == 0 {
		var zero E
		return zero, arr
	}

	ret := arr[0]
	for i := 0; i < len(arr)-1; i++ {
		arr[i] = arr[i+1]
	}

	return ret, arr[:len(arr)-1]
}

type Runner struct {
	debug bool
}

func (r Runner) main() {
	r.delete()
}

func (r Runner) delete() {
	fmt.Println("run delete")
	last := make([]int, 50)
	for i := range 50 {
		last[i] = i + 1
	}

	cases := [][]int{
		[]int{1, 2, 3, 4, 5, 6},
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		[]int{9, 8, 7, 6, 5, 4},
	}

	for i, case_ := range cases {
		fmt.Printf("start case %d with data %v\n", i, case_)
		bTree := newBTree(&r)

		for _, item := range case_ {
			bTree.insert(item)
		}

		if i == 2 {
			bTree.delete(8)
		} else {
			bTree.delete(5)
		}
		bTree.printTree()
		fmt.Println()
		fmt.Println()
	}
}

func (r Runner) insert() {
	fmt.Println("run insert")
	last := make([]int, 50)
	for i := range 50 {
		last[i] = i + 1
	}
	cases := [][]int{
		[]int{10},
		[]int{10, 20, 30},
		[]int{10, 20, 30, 40},
		[]int{10, 20, 30, 40, 50, 60},
		[]int{10, 20, 30, 40, 50, 60, 70, 80},
		[]int{50, 40, 60, 30, 70, 20, 80, 10, 90},
		[]int{50, 30, 70, 20, 40, 60, 80, 10, 90, 100},
		[]int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		[]int{100, 90, 80, 70, 60, 50, 40, 30, 20, 10},
		[]int{50, 20, 80, 40, 70, 30, 60, 10, 90, 100},
		[]int{10, 20, 20, 30},
		last,
	}

	for i, case_ := range cases {
		fmt.Printf("start case %d with data %v\n", i, case_)
		bTree := newBTree(&r)
		fmt.Printf("root %v\n", bTree.root)

		for _, item := range case_ {
			bTree.insert(item)
		}

		fmt.Printf("root %v\n", bTree.root)
		fmt.Println()
		fmt.Println()
	}
}

func main() {
	var debugFlag bool
	flag.BoolVar(&debugFlag, "debug", false, "to show debug log statements")
	flag.Parse()

	runner := Runner{
		debug: debugFlag,
	}

	runner.main()
}
