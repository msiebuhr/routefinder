package routefinder

import (
	"fmt"
	"testing"
)

func Example() {
	r, _ := NewRoutefinder("/shop/:item", "/shop/:item/rate", "/shop/:item/buy")

	fmt.Println(r.Lookup("/shop/gopher/rate"))
	// Output: /shop/:item/rate map[item:gopher]
}

func TestBasic(t *testing.T) {
	r, err := NewRoutefinder("/foo/:id/a", "/foo/:id/b", "/foo/:id")

	if err != nil {
		t.Fatal("Unexpected error creating routes", err)
	}

	tests := []struct {
		p  string
		t  string
		kv map[string]string
	}{
		{
			p:  "/",
			t:  "",
			kv: map[string]string{},
		}, {
			p:  "/foo/foo/a",
			t:  "/foo/:id/a",
			kv: map[string]string{"id": "foo"},
		},
	}

	for _, tt := range tests {

		templ, meta := r.Lookup(tt.p)

		if templ != tt.t {
			t.Errorf("Expected to get route %s, got %s", tt.t, templ)
		}

		for key, value := range tt.kv {
			if data, ok := meta[key]; !ok || data != value {
				t.Errorf("Expected to get %+v, got %+v", tt.kv, meta)
			}
		}
	}
}

func TestPrecedence(t *testing.T) {
	r, _ := NewRoutefinder("/:a", "/:b")

	templ, _ := r.Lookup("/xxxxx")

	if templ != "/:a" {
		t.Errorf("Expected to get route /:a, got %s", templ)
	}
}

func TestStringer(t *testing.T) {
	r, _ := NewRoutefinder("/:a", "/:b")

	if r.String() != "/:a,/:b" {
		t.Errorf("Expected to get `/:a,/:b`, got %s", r.String())
	}
}
