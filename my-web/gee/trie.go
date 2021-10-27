package gee

import (
	"log"
	"strings"
)

type node struct {
	pattern string
	value string
	children []*node
	variable bool
}

func NewNode() *node {
	return &node{}
}

func (n *node) matchChild(value string) *node {
	for _, child := range n.children {
		if value == child.value || child.variable {
			return child
		}
	}

	return nil
}

func (n *node) matchChildren(value string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if value == child.value || child.variable {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) Insert(pattern string, values []string, pointer int) {
	if len(values) == pointer {
		n.pattern = pattern
		return
	}

	value := values[pointer]
	child := n.matchChild(value)
	if child == nil {
		child = &node{
			value: value,
			variable: value[0] == ':' || value[0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.Insert(pattern, values, pointer+1)
}

func (n *node) Print() {
	log.Println(n.pattern, n.value)
	for _, child := range n.children {
		child.Print()
	}
}

func (n *node) Search(values []string, pointer int) *node {
	if len(values) == pointer || strings.HasPrefix(n.value, "*") {
		if n.pattern == "" {
			return nil
		}

		return n
	}
	value := values[pointer]
	children := n.matchChildren(value)
	for _, child := range children {
		if result := child.Search(values, pointer+1); result != nil {
			return result
		}
	}

	return nil
}

