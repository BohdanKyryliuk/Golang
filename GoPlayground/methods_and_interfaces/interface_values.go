package methods_and_interfaces

import (
	"fmt"
	"math"
)

type In interface {
	M()
}

type Type struct {
	S string
}

func (t *Type) M() {
	fmt.Println(t.S)
}

type F float64

func (f F) M() {
	fmt.Println(f)
}

func InterfaceValuesExample() {
	var i In

	i = &Type{"Hello"}
	describe(i)
	i.M()

	i = F(math.Pi)
	describe(i)
	i.M()
}

func describe(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}
