package gee

import "strings"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(pattern string, parts []string, pointer int) {
	if len(parts) == pointer {
		n.pattern = pattern
		return
	}

	part := parts[pointer]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, pointer+1)
}

func (n *node) search(parts []string, pointer int) *node {
	if len(parts) == pointer || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}

		return n
	}

	part := parts[pointer]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, pointer+1)
		if result != nil {
			return result
		}
	}

	return nil
}
