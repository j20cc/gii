package gii

import "strings"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool //是否含有 : 或 *
}

// 匹配第一个，用于插入
func (o *node) matchChild(part string) *node {
	for _, c := range o.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

// 匹配所有，用于搜索
func (o *node) matchChildren(part string) []*node {
	var nodes []*node
	for _, c := range o.children {
		if c.part == part || c.isWild {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

func (o *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 最后一层设置parttern
		o.pattern = pattern
		return
	}

	part := parts[height]
	child := o.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == '*' || part[0] == ':'}
		o.children = append(o.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (o *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(o.part, "*") {
		if o.pattern == "" {
			return nil
		}
		return o
	}

	part := parts[height]
	children := o.matchChildren(part)
	for _, child := range children {
		r := child.search(parts, height+1)
		if r != nil {
			return r
		}
	}

	return nil
}
