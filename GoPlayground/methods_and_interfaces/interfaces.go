package methods_and_interfaces

import (
	"fmt"
	"math"
)

type Abser interface {
	Abs() float64
}

func InterfacesExample() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f // a MyFloat implements Abser
	fmt.Println(a.Abs())

	a = &v // a *Vertex implements Abser

	// a = v // This will cause a compile-time error: Vertex does not implement Abser (Abs method has pointer receiver)

	fmt.Println(a.Abs())
}
