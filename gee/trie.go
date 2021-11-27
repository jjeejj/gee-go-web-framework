package gee

import (
	"strings"
)

// 路径节点
type node struct {
	pattern  string  // 待匹配的路由 /p/:lang
	part     string  // 路由中的一部分 :lang
	children []*node // [doc, intor]
	isWild   bool    // 是否精确匹配 路径中含有 : 或者 * 为精确匹配
}

// 第一个匹配到的节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配到的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 处理完成了路径阶段，给最后一个节点设置匹配的路径然后返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	// 获取当前节点的子节点中 是否存在待匹配的节点
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "'" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
