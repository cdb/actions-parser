package expressions

import (
	"fmt"
	"regexp"
)

// Context of the running action which acts basically like a few different JSON objects.Context
// This is _not_ golang context.Context
type Context map[string]map[string]interface{}

func (e *Expression) Evaluate(context Context) (interface{}, error) {
	if e.Op == nil {
		return e.LHS.Evaluate(context)
	}

	lhs, err := e.LHS.Evaluate(context)
	if err != nil {
		return nil, fmt.Errorf("Error on LSH evaluation: %w", err)
	}

	rhs, err := e.RHS.Evaluate(context)
	if err != nil {
		return nil, fmt.Errorf("Error on RHC evaluation: %w", err)
	}

	switch *e.Op {
	case "==":
		return lhs == rhs, nil
	case "!=":
		return lhs != rhs, nil
	case "<":
		l, _ := lhs.(int64)
		r, _ := rhs.(int64) // TODO: Handle errors, and floats
		return l < r, nil
	default:
		panic(fmt.Sprintf("op '%s' not implemented", *e.Op))
	}
}

var propCleaner = regexp.MustCompile(`[\.\[\]']`)

func (o *Object) Evaluate(context Context) (interface{}, error) {
	var ret interface{}
	ret, ok := context[o.Head]
	if !ok {
		return nil, fmt.Errorf("No root context named %s", o.Head)
	}

	for _, prop := range o.Props {
		prop = propCleaner.ReplaceAllString(prop, "") // Strip the leading . or wrapping [' ']
		switch t := ret.(type) {
		case map[string]interface{}:
			if v, ok := t[prop]; ok {
				ret = v
			}
		default:
			panic("Context not structured as expected")
		}
	}
	return ret, nil
}

func (l *Literal) Evaluate(context Context) (interface{}, error) {
	switch {
	case l.Nil != nil:
		return nil, nil
	case l.Str != nil:
		return *l.Str, nil
	case l.Int != nil:
		return *l.Int, nil
	case l.Float != nil:
		return *l.Float, nil
	case l.Bool != nil:
		return bool(*l.Bool), nil
	case l.Obj != nil:
		o := *l.Obj
		return o.Evaluate(context)
	default:
		panic("empty literal")
	}
}
