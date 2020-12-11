package main

import (
	"fmt"
	"github.com/antonmedv/expr"
)

type VM struct {
	Foo int
	Bar int
}

func main() {

	/*
	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}
	out, err := expr.Eval("foo + bar > 1", env)
	 */

	env := VM{Foo: 1, Bar: 2}
	out, err := expr.Eval("Foo", env)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
}
