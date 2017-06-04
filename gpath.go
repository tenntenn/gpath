package gpath

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/parser"
	"reflect"
)

// At access a field of v by a path.
// v must be struct or pointer of struct.
// A path is represented by Go's expression which can be parsed by go/parser.ParseExpr.
// You can use selectors and indexes in a path.
// Indexes allow only string and int literals for maps.
func At(v interface{}, path string) (interface{}, error) {

	path = "v." + path

	expr, err := parser.ParseExpr(path)
	if err != nil {
		return nil, err
	}

	ev, err := at(reflect.ValueOf(v), expr)
	if err != nil {
		return nil, err
	}

	if ev == (reflect.Value{}) {
		return nil, nil
	}

	return ev.Interface(), nil
}

func at(v reflect.Value, expr ast.Expr) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return at(v.Elem(), expr)
	}

	switch expr := expr.(type) {
	case *ast.Ident:
		return v, nil
	case *ast.SelectorExpr:
		return atBySelector(v, expr)
	case *ast.IndexExpr:
		return atByIndex(v, expr)
	default:
		return reflect.Value{}, errors.New("does not support expr")
	}
}

func direct(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.Elem()
	default:
		return v
	}
}

func atBySelector(v reflect.Value, expr *ast.SelectorExpr) (reflect.Value, error) {
	ev, err := at(v, expr.X)
	if err != nil {
		return reflect.Value{}, err
	}

	ev = direct(ev)
	switch ev.Kind() {
	case reflect.Struct:
		fv := ev.FieldByName(expr.Sel.Name)
		if fv == (reflect.Value{}) {
			return reflect.Value{}, errors.New("cannot find field")
		}
		return fv, nil
	default:
		return reflect.Value{}, errors.New("does not support selector type")
	}
}

func atByIndex(v reflect.Value, expr *ast.IndexExpr) (reflect.Value, error) {
	ev, err := at(v, expr.X)
	if err != nil {
		return reflect.Value{}, err
	}
	ev = direct(ev)

	idx, err := evalExpr(expr.Index)
	if err != nil {
		return reflect.Value{}, err
	}

	switch ev.Kind() {
	case reflect.Slice, reflect.Array:
		i, ok := constant.Int64Val(idx)
		if !ok {
			return reflect.Value{}, errors.New("does not support index type")
		}
		return ev.Index(int(i)), nil
	case reflect.Map:
		switch idx.Kind() {
		case constant.Int:
			k, ok := constant.Int64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not support index type")
			}
			return ev.MapIndex(reflect.ValueOf(int(k))), nil
		case constant.Float:
			k, ok := constant.Float64Val(idx)
			if !ok {
				return reflect.Value{}, errors.New("does not support index type")
			}
			return ev.MapIndex(reflect.ValueOf(k)), nil
		case constant.String:
			k := constant.StringVal(idx)
			return ev.MapIndex(reflect.ValueOf(k)), nil
		case constant.Bool:
			k := constant.BoolVal(idx)
			return ev.MapIndex(reflect.ValueOf(k)), nil
		default:
			return reflect.Value{}, errors.New("does not support index type")
		}
	default:
		return reflect.Value{}, errors.New("does not support expr type")
	}
}
