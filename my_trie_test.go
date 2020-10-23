package httprouter

import (
	"fmt"
	"net/http"
	"testing"
)

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
}

func TestAddRoute(t *testing.T) {
	tree := &node{}
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
	//printChildren(tree, "")
	// check the route
	checkRequests(t, tree, testRequests{
		{"/a", false, "/a"},
		{"/", true, ""},
		{"/hi", false, "/hi"},
		{"/contact", false, "/contact"},
		{"/co", false, "/co"},
		{"/con", true, ""},  // key mismatch
		{"/cona", true, ""}, // key mismatch
		{"/no", true, ""},   // no matching child
		{"/ab", false, "/ab"},
	})
}

func TestWildcardRoute(t *testing.T) {
	route := "/:id/"
	fmt.Printf(route[0:3])
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

var fakeHandlerValue string

func fakeHandler(route string) http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {
		fakeHandlerValue = route
	}
}

func checkRequests(t *testing.T, tree *node, requests testRequests) {
	for _, request := range requests {
		handler := GetValue(tree, request.path)
		if handler == nil {
			if !request.nilHandler {
				t.Errorf("Handler mismatch for route '%s': Expected non-nil handler", request.path)
			}
		} else if request.nilHandler {
			t.Errorf("Handler mismatch for route '%s': Expected nil handler", request.path)
		} else {
			handler(nil, nil)
			if fakeHandlerValue != request.route {
				t.Errorf("Handler mismatch for route '%s': Wrong handler (%s != %s)", request.path, fakeHandlerValue, request.route)
			}
		}
	}
}
