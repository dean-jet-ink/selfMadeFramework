package framework

import (
	"strings"
)

type node struct {
	param    string
	children map[string]*node
	handler  func(*MyContext)
	parent   *node
}

type TrieTree struct {
	root *node
}

func NewTrieTree() *TrieTree {
	root := &node{
		param:    "",
		children: make(map[string]*node),
		handler:  nil,
	}

	return &TrieTree{
		root: root,
	}
}

func (tt *TrieTree) Insert(path string, handler func(*MyContext)) {
	currentNode := tt.root
	params := strings.Split(path, "/")

	for _, param := range params {
		child, ok := currentNode.children[param]

		if !ok {
			child = &node{
				param:    param,
				children: make(map[string]*node),
				handler:  nil,
				parent:   currentNode,
			}
			currentNode.children[param] = child
		}

		currentNode = child
	}

	currentNode.handler = handler
}

func (tt *TrieTree) Search(path string) *node {
	params := strings.Split(path, "/")

	child := tt.depthFirstSearch(params)

	return child
}

func (tt *TrieTree) depthFirstSearch(params []string) *node {
	return tt.depthFirstSearchHelper(params, tt.root)
}

func (tt *TrieTree) depthFirstSearchHelper(params []string, node *node) *node {
	param := params[0]
	isLast := len(params) == 1

	for k, child := range node.children {
		isPathParam := IsPathParam(k)

		if isLast {
			if isPathParam || param == k {
				return child
			}
			continue
		}

		if isPathParam || param == k {
			return tt.depthFirstSearchHelper(params[1:], child)
		}
	}

	return nil
}

func (tt *TrieTree) ParsePath(path string, node *node) map[string]string {
	params := strings.Split(path, "/")
	pathParams := make(map[string]string)

	for i := len(params) - 1; i >= 0; i-- {
		if IsPathParam(node.param) {
			pathParams[node.param] = params[i]
		}

		node = node.parent
	}

	return pathParams
}

func IsPathParam(param string) bool {
	return strings.HasPrefix(param, ":")
}
