package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

const MAX_KEYS = 4
const MIN_KEYS = 1

type BNode struct {
	keys     []int
	children []*BNode
	isLeaf   bool
	runner   *Runner
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
	root *BNode
}

func newBTree(runner *Runner) *BTree {
	return &BTree{
		root: &BNode{
			keys:     []int{},
			children: nil,
			isLeaf:   true,
			runner:   runner,
		},
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
		if tree.root.runner.debug {
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

func (tree *BTree) search(key int) (*BNode, int, bool) {
	if tree.root.runner.debug {
		fmt.Println("start find")
	}
	prev := tree.root
	curr := tree.root

	childIndex := 0
	for {
		i := 0
		if tree.root.runner.debug {
			fmt.Printf("curr at start %v\n", curr)
		}
		for i < len(curr.keys) && key > curr.keys[i] {
			i++
		}

		if i < len(curr.keys) && key == curr.keys[i] {
			break
		}

		if i < len(curr.children) {
			prev = curr
			childIndex = i
			curr = curr.children[i]
			if tree.root.runner.debug {
				fmt.Printf("i %d curr at end %v\n", i, curr)
			}
		} else {
			break
		}
	}

	return prev, childIndex, false
}

func (tree *BTree) delete(key int) {
	// TODO: Implement deletion logic
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

	if tree.root.runner.debug {
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
	last := make([]int, 50)
	for i := range 50 {
		last[i] = i + 1
	}
	cases := [][]int{
		last,
	}

	for i, case_ := range cases {
		fmt.Printf("start case %d with data %v\n", i, case_)
		bTree := newBTree(&r)

		for _, item := range case_ {
			bTree.insert(item)
		}

		bTree.printTree()
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
