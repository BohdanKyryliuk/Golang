package more_types

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	words := strings.Fields(s)

	for _, w := range words {
		m[w] = m[w] + 1
	}

	return m
}

func WordCountExample() {
	sentences := []string{
		"I am learning Go!",
		"The quick brown fox jumped over the lazy dog.",
		"I ate a donut. Then I ate another donut.",
		"A man a plan a canal panama.",
	}

	for _, sentence := range sentences {
		fmt.Println(WordCount(sentence))
	}
}
