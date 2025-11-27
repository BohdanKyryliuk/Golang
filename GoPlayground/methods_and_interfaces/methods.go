package methods_and_interfaces

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Scale method scales the Vertex by a given factor f.
// This method has a pointer receiver, so it can modify the original Vertex.
func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func Scale(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func ScaleFunc(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func AbsFunc(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func MethodsExample() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())

	v.Scale(10)
	fmt.Println(v.Abs())
	fmt.Println(v.X, v.Y)

	Scale(&v, 10)
	fmt.Println(Abs(v))

	v1 := Vertex{3, 4}
	v1.Scale(2)
	ScaleFunc(&v1, 10)

	p := &Vertex{4, 3}
	p.Scale(3)
	ScaleFunc(p, 8)

	fmt.Println(v1, p)

	v2 := Vertex{3, 4}
	fmt.Println(v2.Abs())
	fmt.Println(AbsFunc(v2))

	p2 := &Vertex{4, 3}
	fmt.Println(p2.Abs())
	fmt.Println(AbsFunc(*p2))
}
