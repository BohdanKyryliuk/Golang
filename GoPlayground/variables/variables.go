package variables

import "fmt"

// Predefined variables
var c, python, java, php, golang = true, false, "no!", `yes!`, "Golang!"
var o, j int = 1, 2
var Email string = "test@test.com" // email is exported because it starts with a capital letter
var password string                // password is unexported because it starts with a lowercase letter
var g = "global"
var num1 = 5

// Constants
// can be character, string, boolean, or numeric values.
const shark = "Sammy"
const (
	year     = 365        // untyped constant
	leapYear = int32(366) // typed constant
)
const PI = 3.1415926535
const (
	// Big Create a huge number by shifting a 1 bit left 100 places.
	// In other words, the binary number that is 1 followed by 100 zeroes.
	Big = 1 << 100
	// Small Shift it right again 99 places, so we end up with 1<<1, or 2.
	Small = Big >> 99
)

func Variables() {
	var y int
	cPlusPlus, pythonDjango, javaScript := true, false, "no!"

	fmt.Println(o, j, y, c, python, java, php, golang, cPlusPlus, pythonDjango, javaScript)

	// Short variable declaration
	x := 42
	fmt.Println(x)

	s := "Hello, World!"
	f := 45.06
	boolean := 5 > 9
	array := [4]string{"item_1", "item_2", "item_3", "item_4"}
	slice := []string{"one", "two", "three"}
	m := map[string]string{"letter": "g", "number": "seven", "symbol": "&"}

	fmt.Println(s, f, boolean, array, slice, m)

	zeroValues()

	printNames()

	multipleAssignment()

	// Global and local variables
	printLocal()
	fmt.Println(g)

	printNumbers()
	fmt.Println(num1)

	constants()
}

func zeroValues() {
	var a int
	var b string
	var float float64
	var d bool
	fmt.Printf("var a %T = %+v\n", a, a)
	fmt.Printf("var b %T = %q\n", b, b)
	fmt.Printf("var float %T = %+v\n", float, float)
	fmt.Printf("var d %T = %+v\n", d, d)
}

func printNames() {
	names := []string{"Mary", "John", "Bob", "Anna"}
	for i, n := range names {
		fmt.Printf("index: %d = %q\n", i, n)
	}
}

func multipleAssignment() {
	m, k, l := "shark", 2.05, 15
	fmt.Println(m, k, l)
}

func constants() {
	// Constant
	fmt.Println("Constant shark:", shark)
	//shark = "Toddy" // This will cause a compile-time error: cannot assign to shark

	hours := 24
	minutes := int32(60)
	fmt.Println(hours * year)
	fmt.Println(minutes * year)
	fmt.Println(minutes * leapYear)

	const World = "世界"
	fmt.Println("Hello", World)
	fmt.Println("Happy", PI, "Day")

	const Truth = true
	fmt.Println("Go rules?", Truth)
}

func printLocal() {
	l := "local"
	fmt.Println(l)
	fmt.Println(g)
}

func printNumbers() {
	num1 := 10
	num2 := 7

	fmt.Println(num1)
	fmt.Println(num2)
}

func needInt(x int) int {
	return x*10 + 1
}
func needFloat(x float64) float64 {
	return x * 0.1
}

func NumericConstants() {
	// Numeric Constants
	fmt.Println(needInt(Small))
	//fmt.Println(needInt(Big)) // cannot use Big (untyped int constant 1267650600228229401496703205376) as int value in argument to needInt (overflows)
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))
}
