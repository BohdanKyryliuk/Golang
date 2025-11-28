package main

import (
	"Golang/GoPlayground"
	"Golang/currency_converter"
	"Golang/greeter"
	"fmt"
	"log"
	"time"

	"rsc.io/quote"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
	// to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!\n", s)

	for i := 1; i <= 5; i++ {
		//TIP <p>To start your debugging session, right-click your code in the editor and select the Debug option.</p> <p>We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>.</p>
		fmt.Println("i =", 100/i)
	}

	fmt.Println(quote.Go())

	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	message, err := greeter.Hello("World")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)

	message, err = greeter.Hello("Bohdan")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)

	message, err = greeter.Hello("Gladys")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)

	// Uncomment to see the error handling in action
	//message, err = greeter.Hello("")
	// If an error was returned, print it to the console and
	// exit the program.
	//if err != nil {
	//	log.Fatal(err)
	//}

	// If no error was returned, print the returned message
	// to the console.
	//fmt.Println(message)

	// A slice of names.
	names := []string{"Gladys", "Samantha", "Darrin"}

	// Request greeting messages for the names.
	messages, err := greeter.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned map of
	// messages to the console.
	fmt.Println(messages)

	// Uncomment to see the error handling in action
	//names = []string{"Bohdan", "Kyryliuk", ""}
	//messages, err = greeter.Hellos(names)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(messages)

	// Using goroutine to call Hello function concurrently
	hello, err := greeter.Hello("Bohdan Kyryliuk")
	if err != nil {
		log.Fatal(err)
	}

	go fmt.Println(hello)
	time.Sleep(time.Second) // Wait for the goroutine to finish

	// Using channels to communicate between goroutines
	ch := make(chan string)
	go func() {
		ch <- "Hello, Channel!"
	}()

	fmt.Println(<-ch)

	// The Go Playground
	GoPlayground.Playground()

	// Setting up HTTP server with handlers
	//web.StartServer()

	fmt.Println()

	// Currency API example
	status := currency_converter.CheckStatus()
	fmt.Println("Currency API Status:", status)

	fmt.Println()

	// Get and print supported currencies
	currencies := currency_converter.GetCurrencies()
	fmt.Println("Supported Currencies:", currencies)

	// Get and print latest currency rates
	latestRates := currency_converter.GetLatestRates()
	fmt.Println("Latest Currency Rates:", latestRates)
}
