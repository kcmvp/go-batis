package sql

import (
	"fmt"
	"github.com/antonmedv/expr"
	"reflect"
)

type VM struct {
	Foo int
	Bar *int
}

func main() {

	env := map[string]interface{}{
		"foo": 1,
		"bar": 2,
	}
	vm := VM{
		Foo: 2,
	}
	//out, err := expr.Eval("foo + bar > 1", env)
	out, err := expr.Eval("Bar", vm)


	//env := VM{Foo: 1, Bar: 2}
	//out, err := expr.Eval("Foo", env)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

	s := reflect.ValueOf(env)

	fmt.Println(s)
	fmt.Println(s.Len())
	fmt.Println(s.Index(0))





}
