package variables

import "fmt"

// Predefined variables
var c, python, java, php, golang = true, false, "no!", `yes!`, "Golang!"
var i, j int = 1, 2
var Email string = "test@test.com" // email is exported because it starts with a capital letter
var password string                // password is unexported because it starts with a lowercase letter

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

func Variables() {
	var y int
	cPlusPlus, pythonDjango, javaScript := true, false, "no!"

	fmt.Println(i, j, y, c, python, java, php, golang, cPlusPlus, pythonDjango, javaScript)

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
}
