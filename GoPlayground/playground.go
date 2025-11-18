package GoPlayground

import (
	"Golang/GoPlayground/data_types"
	"Golang/GoPlayground/functions"
	"Golang/GoPlayground/iota"
	"Golang/GoPlayground/loops"
	"Golang/GoPlayground/variables"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Playground() {
	printLabel("Welcome to the Go Playground!")
	fmt.Println("The time is", time.Now())

	// Packages
	printLabel("Packages:")
	fmt.Println("My favorite number is", rand.Intn(10))

	// Imports
	printLabel("Imports:")
	fmt.Printf("Now you have %g problems.\n", math.Sqrt(7))

	// Exported names
	printLabel("Exported Names:")
	// fmt.Println(math.pi) // This will cause a compile-time error: math.pi is not exported
	fmt.Println(math.Pi)

	// Functions
	printLabel("Functions:")
	fmt.Println(functions.Add(42, 13))

	// Multiple return values
	printLabel("Multiple Return Values:")
	a, b := functions.Swap("hello", "world")
	fmt.Println(a, b)

	// Named return values
	printLabel("Named Return Values:")
	fmt.Println(functions.Split(17))

	// Variables
	printLabel("Variables:")
	variables.Variables()

	fmt.Println("Email is", variables.Email) // Can use the exported variable 'Email' in the current package
	// fmt.Println(variables.password) // Cannot use the unexported variable 'password' in the current package

	// Data Types
	printLabel("Data Types:")
	data_types.DataTypes()

	// Type Conversions
	printLabel("Type Conversions:")
	data_types.TypeConversions()

	// Type Inference
	printLabel("Type Inference:")
	data_types.TypeInference()

	// Numeric Constants
	printLabel("Numeric Constants:")
	variables.NumericConstants()

	// Loops
	printLabel("Loops:")
	loops.Loops()

	printLabel("Loop with Optional Init and Post:")
	loops.LoopWithOptionalInitAndPost()

	// iota
	printLabel("iota:")
	printLabel("Byte Sizes using iota:")
	iota.PrintIotaBytes()

	printLabel("Days of the Week using iota:")
	iota.PrintWeekdays()
}

func printLabel(label string) {
	fmt.Println()
	fmt.Println(label)
}
