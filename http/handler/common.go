package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Hello handles the hello world endpoint
func Hello(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	// Step 1: Write a simple Hello World response
	c.Writer.WriteString("<h1>Hello, World!</h1>")

	// Step 2: Read query parameter and respond accordingly
	q := c.Query("q")
	if q != "" {
		c.Writer.WriteString(fmt.Sprintf(`
           <body>
               <h1>Hello, %s!</h1>
           </body>
       `, q))
		return
	}

	c.Writer.WriteString(`
       <body>
           <form action="/" method="GET">
               <label>Enter your name</label>
               <input name="q">
               <button type="submit">Submit</button>
           </form>
       </body>
   `)
}

// Counter handles the counter endpoint
func Counter(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	// Step 3: Handle counter with POST method
	if c.Request.Method == "POST" {
		count, _ := strconv.Atoi(c.PostForm("counter"))
		count++
		c.Writer.WriteString(fmt.Sprintf(`
            <body>
                <form action="/count" method="POST">
                    <label>Counter</label>
                    <input name="counter" value="%d" readonly>
                    <button type="submit">Add</button>
                </form>
                <a href="/count">Reset</a>
            </body>
        `, count))
		return
	}

	c.Writer.WriteString(`
            <body>
                <form action="/count" method="POST">
                    <label>Counter</label>
                    <input name="counter" value="1" readonly>
                    <button type="submit">Add</button>
                </form>
            </body>
            <a href="/count">Reset</a>
        `)
}
