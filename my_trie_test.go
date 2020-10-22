package httprouter

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAddRoute(t *testing.T) {
	tree := node{}
	routes := []string{
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
	}
	for _, route := range routes {
		if err := tree.addRoute(route, fakeHandler(route)); err != nil {
			t.Fatalf("Error inserting route '%s': %s", route, err.Error())
		}
	}
	printChildren(&tree,"")
}

func printChildren(n *node, prefix string) {
	fmt.Printf("%s%s\n", prefix, n.key)
	for l := len(n.key); l > 0; l-- {
		prefix += " "
	}
	for _, child := range n.children {
		printChildren(child, prefix)
	}
}

func fakeHandler(route string) http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
	}
}
