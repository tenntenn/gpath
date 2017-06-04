package gpath_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/tenntenn/gpath"
)

type Hoge struct {
	Foo *Foo
}

type Foo struct {
	Bar *Bar
}

type Bar struct {
	N int
}

func ExampleAt() {
	type Bar struct {
		N []int
	}

	type Foo struct {
		Bar *Bar
	}

	f := &Foo{
		Bar: &Bar{
			N: []int{100},
		},
	}

	v, err := At(f, `Bar.N[0]`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(v)
	}
	// Output:
	// 100
}

func TestAt(t *testing.T) {
	data := []struct {
		d      interface{}
		p      string
		e      interface{}
		hasErr bool
	}{
		{
			d: struct {
				A int
			}{
				A: 100,
			},
			p: "A",
			e: 100,
		},
		{
			d: struct {
				A []int
			}{
				A: []int{100, 200},
			},
			p: "A[1]",
			e: 200,
		},
		{
			d: struct {
				A map[string][]int
			}{
				A: map[string][]int{
					"foo": {100, 200},
				},
			},
			p: `A["foo"][1]`,
			e: 200,
		},
		{
			d: struct {
				A map[int][]int
			}{
				A: map[int][]int{
					200: {100, 200},
				},
			},
			p: `A[200][1]`,
			e: 200,
		},
		{
			d: struct {
				A struct {
					B int
				}
			}{
				A: struct {
					B int
				}{
					B: 100,
				},
			},
			p: `A.B`,
			e: 100,
		},
		{
			d: struct {
				A struct {
					B int
				}
			}{
				A: struct {
					B int
				}{
					B: 100,
				},
			},
			p:      `A.C`,
			hasErr: true,
		},
		{
			d: &Hoge{
				Foo: &Foo{
					Bar: &Bar{
						N: 100,
					},
				},
			},
			p: `Foo.Bar`,
			e: &Bar{N: 100},
		},
		{
			d: struct {
				N []int
			}{
				N: []int{100},
			},
			p: `N[0]`,
			e: 100,
		},
		{
			p:      `import "fmt"`,
			hasErr: true,
		},
		{
			p:      `Call()`,
			hasErr: true,
		},
		{
			d: struct {
				N map[int]int
			}{
				N: map[int]int{-1: 100},
			},
			p: `N[-1]`,
			e: 100,
		},
		{
			d: struct {
				N map[int]int
			}{
				N: map[int]int{0: 100},
			},
			p: `N[1-1]`,
			e: 100,
		},
		{
			d: struct {
				N map[int]int
			}{
				N: map[int]int{0: 100},
			},
			p: `N[(0)]`,
			e: 100,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p: `N[true]`,
			e: 100,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{false: 100},
			},
			p: `N[false]`,
			e: 100,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p: `N[100 > 0]`,
			e: 100,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p: `N[true]`,
			e: 100,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p:      `N[1 + f()]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p:      `N[T]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p:      `N[-T]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p:      `N[T - 1]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[bool]int
			}{
				N: map[bool]int{true: 100},
			},
			p:      `N["key" + 1]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[int]int
			}{
				N: map[int]int{10: 100},
			},
			p:      `N[99999999999999999999999999999]`,
			hasErr: true,
		},
		{
			d: struct {
				N []int
			}{
				N: []int{100},
			},
			p:      `N[99999999999999999999999999999]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[float64]int
			}{
				N: map[float64]int{1.5: 100},
			},
			p: `N[1.5]`,
			e: 100,
		},
		{
			d: struct {
				N map[float64]int
			}{
				N: map[float64]int{1.5: 100},
			},
			p:      `N[99999999999999999999999999999999999.0]`,
			hasErr: true,
		},
		{
			d: struct {
				N map[int]int
			}{
				N: map[int]int{100: 100},
			},
			p: `N[0]`,
			e: nil,
		},
		{
			d:      100,
			p:      `N[0]`,
			hasErr: true,
		},
		{
			d:      100,
			p:      `A.B[0]`,
			hasErr: true,
		},
		{
			d: struct {
				N int
			}{
				N: 100,
			},
			p:      `N[0]`,
			hasErr: true,
		},
	}

	for i := range data {
		a, err := At(data[i].d, data[i].p)
		if data[i].hasErr {
			if err == nil {
				t.Errorf("expected error did not occur")
				continue
			}
		} else if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		if !reflect.DeepEqual(a, data[i].e) {
			t.Errorf("case %d: got %v want %v", i, a, data[i].e)
		}
	}
}
