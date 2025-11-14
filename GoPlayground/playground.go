package GoPlayground

import (
	"Golang/GoPlayground/functions"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Playground() {
	fmt.Println("Welcome to the Go Playground!")
	fmt.Println("The time is", time.Now())

	// Packages
	fmt.Println("My favorite number is", rand.Intn(10))

	// Imports
	fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))

	// Exported names
	// fmt.Println(math.pi) // This will cause a compile-time error: math.pi is not exported
	fmt.Println(math.Pi)

	// Functions
	fmt.Println(functions.Add(42, 13))
}
