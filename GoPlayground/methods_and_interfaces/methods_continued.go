package methods_and_interfaces

import (
	"fmt"
	"math"
)

type MyFloat float64

// Abs method calculates the absolute value of MyFloat.
// You can declare a method on non-struct types, too.
func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func MethodsContinued() {
	// Example usage of MyFloat and its Abs method
	f := MyFloat(-math.Sqrt2)
	fmt.Println(f.Abs()) // Call the Abs method
}
