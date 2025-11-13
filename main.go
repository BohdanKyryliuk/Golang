package main

import (
	"Golang/HttpHandler"
	"Golang/greeter"
	"fmt"
	"log"
	"net/http"
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

	greeter.Hello("World")
	greeter.Hello("Bohdan")

	go greeter.Hello("Bohdan Kyryliuk")
	time.Sleep(time.Second * 3)

	// Setting up HTTP server with handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", HttpHandler.HelloHandler)
	mux.HandleFunc("/count", HttpHandler.CounterHandler)

	log.Println("Listening on :3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
