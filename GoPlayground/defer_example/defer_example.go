package defer_example

import "fmt"

func DeferExample() {
	defer fmt.Println("world")

	fmt.Println("hello")
}

func StackingDefers() {
	fmt.Println("Counting")

	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}

	fmt.Println("Done")
}
