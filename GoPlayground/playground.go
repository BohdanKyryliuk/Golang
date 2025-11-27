package GoPlayground

import (
	"Golang/GoPlayground/conditions"
	"Golang/GoPlayground/data_types"
	"Golang/GoPlayground/defer_example"
	"Golang/GoPlayground/functions"
	"Golang/GoPlayground/iota"
	"Golang/GoPlayground/loops"
	"Golang/GoPlayground/methods_and_interfaces"
	"Golang/GoPlayground/more_types"
	"Golang/GoPlayground/switch"
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

	printLabel("Infinite Loop Example:")
	loops.InfiniteLoopExample()

	// iota
	printLabel("iota:")
	printLabel("Byte Sizes using iota:")
	iota.PrintIotaBytes()

	printLabel("Days of the Week using iota:")
	iota.PrintWeekdays()

	// Conditions
	printLabel("Conditions:")
	fmt.Println(conditions.IfExample(5))
	fmt.Println(conditions.IfWithShortStatement(6))
	fmt.Println(conditions.NestedIfExample(-3))
	fmt.Println(conditions.IfWithoutElse(0))
	fmt.Println(conditions.IfMultipleConditions(15))
	fmt.Println(conditions.IfWithLogicalOperators(50))
	fmt.Println(conditions.Sqrt(-4))

	printLabel("Switch Example:")
	switch_example.SwitchExample()

	printLabel("When is Saturday?:")
	switch_example.PrintWhenIsSaturday()

	printLabel("Greetings based on Time of Day:")
	switch_example.Greetings()

	printLabel("Defer example:")
	defer_example.DeferExample()

	printLabel("Stacking defers:")
	defer_example.StackingDefers()

	printLabel("Pointers example:")
	more_types.PointersExample()

	printLabel("Structs example:")
	more_types.StructsExample()

	printLabel("Struct Literal example:")
	more_types.StructLiteralExample()

	printLabel("Slice Literals example:")
	more_types.SliceLiteralsExample()

	printLabel("Slice Defaults example:")
	more_types.SliceDefaultsExample()

	printLabel("Slice Length and Capacity example:")
	more_types.SliceLengthAndCapacityExample()

	printLabel("Nil Slice example:")
	more_types.NilSliceExample()

	printLabel("Slices of Slices example:")
	more_types.SlicesOfSlicesExample()

	printLabel("Range example:")
	more_types.RangeExample()

	printLabel("Maps example:")
	more_types.MapsExample()

	printLabel("Mutating Maps example:")
	more_types.MutatingMapsExample()

	printLabel("Word Count example:")
	more_types.WordCountExample()

	printLabel("Function Values example:")
	more_types.FunctionValuesExample()

	printLabel("Function Closures example:")
	more_types.FunctionClosuresExample()

	printLabel("Methods example:")
	methods_and_interfaces.MethodsExample()
}

func printLabel(label string) {
	fmt.Println()
	fmt.Println(label)
}
