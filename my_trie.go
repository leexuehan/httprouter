package httprouter

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrorDuplicatedRoute = errors.New("duplicated route")
)

// node:represent a new route
type node struct {
	key      string
	value    http.Handler
	children []*node
	indices  []byte
}

func (n *node) addRoute(key string, value http.Handler) error {
	if len(n.key) == 0 {
		// empty tree,insert it directly
		n.key = key
		n.value = value
		return nil
	}
	return doAdd(n, key, value)
}

func doAdd(n *node, key string, value http.Handler) error {
	// find the longest common prefix
	var i int
	for j := min(len(key), len(n.key)); i < j && key[i] == n.key[i]; i++ {
	}

	// split edge
	if i < len(n.key) {
		childKey := n.key[i:]
		n.children = []*node{
			{
				key:      childKey,
				value:    n.value,
				children: n.children,
			},
		}
		n.indices = []byte{n.key[i]}
		n.key = key[:i]
		n.value = nil
		fmt.Printf("split an edge,insert key[%s],root key[%s],child key[%s]\n", key, n.key, childKey)
	}

	// make new node a child of this node
	if i < len(key) {
		key = key[i:]
		c := key[0]
		for i, index := range n.indices {
			if c == index {
				n = n.children[i]
				return doAdd(n, key, value)
			}
		}
		n.indices = append(n.indices, c)
		child := &node{}
		n.children = append(n.children, child)
		fmt.Printf("insert a new child nod, parent key[%s],child key[%s],indices[%s]\n", n.key, key, n.indices)
		n = child
		n.key = key
		n.value = value
	} else if i == len(key) {
		if n.value != nil {
			return ErrorDuplicatedRoute
		}
		n.value = value
		return nil
	}
	return nil
}

func (n *node) insertRoute(key string, value http.HandlerFunc) error {
	n.key = key
	n.value = value
	return nil
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
