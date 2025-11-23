package more_types

import "fmt"

type Vertex struct {
	X int
	Y int
}

func StructsExample() {
	// Structs are collections of fields

	// Create a new struct
	v := Vertex{X: 1, Y: 2}
	fmt.Println(v)
	fmt.Println(v.X)

	// Access struct fields
	_ = v.X
	_ = v.Y

	// You can also use a pointer to a struct
	p := &v
	_ = p.X   // Accessing field X through pointer p
	p.X = 1e9 // Modifying field X through pointer p
	fmt.Println(v)

	// Structs can be nested
	type Circle struct {
		Center Vertex
		Radius int
	}

	c := Circle{Center: Vertex{X: 0, Y: 0}, Radius: 5}
	_ = c.Center.X
	_ = c.Radius
}
