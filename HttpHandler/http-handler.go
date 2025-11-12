package HttpHandler

import (
	"fmt"
	"net/http"
	"strconv"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Step 1: Write a simple Hello World response
	//_, _ = w.Write([]byte("<h1>Hello, World!</h1>"))

	// Step 2: Read query parameter and respond accordingly
	//if r.FormValue("q") != "" {
	//	fmt.Fprintf(w, `
	//            <body>
	//                <h1>Hello, %s!</h1>
	//            </body>
	//        `, r.FormValue("q"))
	//	return
	//}

	//fmt.Fprintf(w, `
	//        <body>
	//            <form action="/" method="GET">
	//                <label>Enter your name</label>
	//                <input name="q">
	//                <button type="submit">Submit</button>
	//            </form>
	//        </body>
	//    `)

	// Step 3: Handle counter with POST method
	if r.Method == "POST" {
		count, _ := strconv.Atoi(r.FormValue("counter"))
		count++
		fmt.Fprintf(w, `
            <body>
                <form action="/" method="POST">
                    <label>Counter</label>
                    <input name="counter" value="%d" readonly>
                    <button type="submit">Add</button>
                </form>
                <a href="/">Reset</a>
            </body>
        `, count)
		return
	}

	fmt.Fprintf(w, `
            <body>
                <form action="/" method="POST">
                    <label>Counter</label>
                    <input name="counter" value="1" readonly>
                    <button type="submit">Add</button>
                </form>
            </body>
            <a href="/">Reset</a>
        `)
}
