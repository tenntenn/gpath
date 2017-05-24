package gpath_test

import (
	"reflect"
	"testing"

	. "github.com/tenntenn/gpath"
)

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
					"foo": []int{100, 200},
				},
			},
			p: `A["foo"][1]`,
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
			t.Errorf("got %v want %v", a, data[i].e)
		}
	}
}
