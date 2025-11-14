package GoPlayground

import (
	"Golang/GoPlayground/functions"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Predefined variables
var c, python, java, php, golang = true, false, "no!", `yes!`, "Golang!"
var i, j int = 1, 2

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

	// Multiple return values
	a, b := functions.Swap("hello", "world")
	fmt.Println(a, b)

	// Named return values
	fmt.Println(functions.Split(17))

	// Variables
	var y int
	cPlusPlus, pythonDjango, javaScript := true, false, "no!"

	fmt.Println(i, j, y, c, python, java, php, golang, cPlusPlus, pythonDjango, javaScript)

	// Short variable declaration
	x := 42
	fmt.Println(x)
}
