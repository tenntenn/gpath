package gpath

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// At access a field of v by a path.
// v must be struct or pointer of struct.
// A path is represented by Go's expression which can be parsed by go/parser.ParseExpr.
// You can use selectors and indexes in a path.
// Indexes allow only string and int literals for maps.
func At(v interface{}, path string) (interface{}, error) {

	if strings.HasPrefix(path, "[") {
		path = "v" + path
	} else {
		path = "v." + path
	}

	expr, err := parser.ParseExpr(path)
	if err != nil {
		return nil, err
	}

	ev, err := at(reflect.ValueOf(v), expr)
	if err != nil {
		return nil, err
	}

	return ev.Interface(), nil
}

func at(v reflect.Value, expr ast.Expr) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return at(v.Elem(), expr)
	}

	switch expr := expr.(type) {
	case nil:
		return v, nil
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

	bl, ok := expr.Index.(*ast.BasicLit)
	if !ok {
		return reflect.Value{}, errors.New("does not support index type")
	}

	switch ev.Kind() {
	case reflect.Slice, reflect.Array:
		i, err := intIndex(bl)
		if err != nil {
			return reflect.Value{}, err
		}
		return ev.Index(i), nil
	case reflect.Map:
		switch bl.Kind {
		case token.INT:
			k, err := intIndex(bl)
			if err != nil {
				return reflect.Value{}, err
			}
			return ev.MapIndex(reflect.ValueOf(k)), nil
		case token.STRING:
			k, err := stringIndex(bl)
			if err != nil {
				return reflect.Value{}, err
			}
			return ev.MapIndex(reflect.ValueOf(k)), nil
		default:
			return reflect.Value{}, errors.New("does not support index type")
		}
	default:
		return reflect.Value{}, errors.New("does not support expr type")
	}
}

func intIndex(bl *ast.BasicLit) (int, error) {
	if bl.Kind != token.INT {
		return 0, errors.New("does not support index type")
	}

	cv := constant.MakeFromLiteral(bl.Value, bl.Kind, 0)
	i, ok := constant.Int64Val(cv)
	if !ok {
		return 0, errors.New("does not support index type")
	}

	return int(i), nil
}

func stringIndex(bl *ast.BasicLit) (string, error) {
	if bl.Kind != token.STRING {
		return "", errors.New("does not support index type")
	}
	cv := constant.MakeFromLiteral(bl.Value, bl.Kind, 0)
	return constant.StringVal(cv), nil
}
