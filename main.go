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

func (tree *BTree) delete(key int) {
	parent, childIndex, keyIndex := tree.search(key)
	if keyIndex == -1 {
		return
	}

	node := parent.children[childIndex]

	if !node.isLeaf {
		predecessor := tree.getPredecessor(node, keyIndex)
		node.keys[keyIndex] = predecessor
		tree.deleteFromSubtree(node, childIndex, predecessor)
		return
	}

	tree.deleteFromLeaf(node, parent, childIndex, keyIndex)
}

func (tree *BTree) deleteFromLeaf(node *BNode, parent *BNode, childIndex, keyIndex int) {
	node.keys = slices.Delete(node.keys, keyIndex, keyIndex+1)

	if len(node.keys) < MIN_KEYS && parent != nil {
		tree.rebalance(parent, childIndex)
	}
}

func (tree *BTree) deleteFromSubtree(node *BNode, childIndex int, key int) {
	child := node.children[childIndex]
	if len(child.keys) >= MIN_KEYS {
		tree.delete(key)
	} else {
		tree.rebalance(node, childIndex)
		tree.delete(key)
	}
}

func (tree *BTree) rebalance(parent *BNode, childIndex int) {
	node := parent.children[childIndex]
	leftSiblingIndex := childIndex - 1
	rightSiblingIndex := childIndex + 1

	if leftSiblingIndex >= 0 && len(parent.children[leftSiblingIndex].keys) > MIN_KEYS {
		leftSibling := parent.children[leftSiblingIndex]
		borrowedKey := leftSibling.keys[len(leftSibling.keys)-1]
		leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]

		node.keys = append([]int{parent.keys[leftSiblingIndex]}, node.keys...)
		parent.keys[leftSiblingIndex] = borrowedKey

		if !leftSibling.isLeaf {
			child := leftSibling.children[len(leftSibling.children)-1]
			leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
			node.children = append([]*BNode{child}, node.children...)
		}
		return
	}

	if rightSiblingIndex < len(parent.children) && len(parent.children[rightSiblingIndex].keys) > MIN_KEYS {
		rightSibling := parent.children[rightSiblingIndex]
		borrowedKey := rightSibling.keys[0]
		rightSibling.keys = rightSibling.keys[1:]

		node.keys = append(node.keys, parent.keys[childIndex])
		parent.keys[childIndex] = borrowedKey

		if !rightSibling.isLeaf {
			child := rightSibling.children[0]
			rightSibling.children = rightSibling.children[1:]
			node.children = append(node.children, child)
		}
		return
	}

	if leftSiblingIndex >= 0 {
		tree.mergeNodes(parent, leftSiblingIndex, childIndex)
	} else if rightSiblingIndex < len(parent.children) {
		tree.mergeNodes(parent, childIndex, rightSiblingIndex)
	}
}

func (tree *BTree) mergeNodes(parent *BNode, leftIndex, rightIndex int) {
	leftNode := parent.children[leftIndex]
	rightNode := parent.children[rightIndex]

	leftNode.keys = append(leftNode.keys, parent.keys[leftIndex])
	leftNode.keys = append(leftNode.keys, rightNode.keys...)
	if !leftNode.isLeaf {
		leftNode.children = append(leftNode.children, rightNode.children...)
	}

	parent.keys = slices.Delete(parent.keys, leftIndex, leftIndex+1)
	parent.children = slices.Delete(parent.children, rightIndex, rightIndex+1)

	if len(parent.keys) < MIN_KEYS && parent != tree.root {
		grandparent, grandparentIndex, _ := tree.search(parent.keys[0])
		tree.rebalance(grandparent, grandparentIndex)
	}

	if len(parent.keys) == 0 && parent == tree.root {
		tree.root = leftNode
	}
}

func (tree *BTree) getPredecessor(node *BNode, keyIndex int) int {
	current := node.children[keyIndex]
	for !current.isLeaf {
		current = current.children[len(current.children)-1]
	}
	return current.keys[len(current.keys)-1]
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
	debug   bool
	verbose bool
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

	cases := []struct {
		values []int
		key    int
	}{
		{values: []int{1, 2, 3, 4, 5, 6}, key: 5},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, key: 5},
		{values: []int{9, 8, 7, 6, 5, 4}, key: 8},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}, key: 8},
		{values: []int{10, 20, 30, 40, 50, 5, 15, 25, 35, 45, 55, 60, 70, 80, 90, 65}, key: 8},
		{values: []int{10, 20, 30, 40, 50, 5, 15, 25, 35, 45, 55, 60, 70, 80, 90, 65, 95, 85, 75}, key: 55},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8}, key: 7},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8}, key: 5},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, key: 16},
		{values: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}, key: 16},
	}

	for i, case_ := range cases {
		fmt.Printf("start case %d with data %v\n", i, case_.values)
		fmt.Printf("delete key %d\n", case_.key)
		bTree := newBTree(&r)

		for _, item := range case_.values {
			bTree.insert(item)
		}

		if r.verbose {
			fmt.Println("before delete")
			bTree.printTree()
			fmt.Println()
			fmt.Println()
		}
		bTree.delete(case_.key)
		fmt.Println("after delete")
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
	var verboseFlag bool
	flag.BoolVar(&debugFlag, "debug", false, "to show debug log statements")
	flag.BoolVar(&verboseFlag, "verbose", false, "to show verbose mode in test")
	flag.Parse()

	runner := Runner{
		debug:   debugFlag,
		verbose: verboseFlag,
	}

	runner.main()
}
