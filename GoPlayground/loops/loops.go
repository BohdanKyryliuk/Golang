package loops

import "fmt"

func Loops() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println("Sum:", sum)
}

func LoopWithOptionalInitAndPost() {
	sum := 1
	for sum < 1000 {
		sum += sum
	}
	fmt.Println("Sum with optional init and post:", sum)
}

func InfiniteLoopExample() {
	count := 0
	for {
		if count >= 5 {
			break
		}
		fmt.Println("Count:", count)
		count++
	}
}
