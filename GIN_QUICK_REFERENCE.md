# Gin Framework Quick Reference for This Project

## Handler Function Patterns

### Simple Route Handler
```go
// Define handler
func MyHandler(c *gin.Context) {
    c.String(200, "Hello World")
}

// Register route
router.GET("/my-route", MyHandler)
```

### Handler with Query Parameters
```go
func MyHandler(c *gin.Context) {
    // Get query parameter
    page := c.Query("page")
    
    // Get with default value
    limit := c.DefaultQuery("limit", "10")
    
    // Get required parameter (returns 400 if missing)
    id := c.Param("id")
    
    c.JSON(200, gin.H{
        "page": page,
        "limit": limit,
    })
}

router.GET("/items/:id", MyHandler)
```

### Handler with Form Data (POST)
```go
func MyHandler(c *gin.Context) {
    // Get form field
    name := c.PostForm("name")
    
    // Get with default value
    email := c.DefaultPostForm("email", "")
    
    c.JSON(200, gin.H{
        "name": name,
        "email": email,
    })
}

router.POST("/submit", MyHandler)
```

### Handler with JSON Response
```go
type Response struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
}

func MyHandler(c *gin.Context) {
    data := Response{ID: 1, Name: "John"}
    c.JSON(200, data) // Automatic JSON marshaling
}

router.GET("/api/user", MyHandler)
```

### Error Handling
```go
func MyHandler(c *gin.Context) {
    if someError {
        // Return error with status code
        c.AbortWithStatusJSON(400, gin.H{
            "error": "Invalid request",
        })
        return
    }
    
    c.JSON(200, gin.H{"success": true})
}
```

### Struct Method Handler
```go
type MyHandler struct {
    db *sql.DB
}

func (h *MyHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    // Use h.db to fetch user
    c.JSON(200, gin.H{"user_id": id})
}

// Register
handler := &MyHandler{db: db}
router.GET("/users/:id", handler.GetUser)
```

## Route Registration Examples

### Individual Routes
```go
router := gin.Default()

router.GET("/hello", handler.Hello)
router.POST("/submit", handler.Submit)
router.DELETE("/items/:id", handler.DeleteItem)
```

### Grouped Routes
```go
router := gin.Default()

// Currency endpoints
currencyGroup := router.Group("/currency")
{
    currencyGroup.GET("/status", currencyHandler.Status)
    currencyGroup.GET("/currencies", currencyHandler.Currencies)
    currencyGroup.GET("/latest", currencyHandler.LatestRates)
}

// Rates endpoints
ratesGroup := router.Group("/rates")
{
    ratesGroup.GET("", ratesHandler.GetRate)
    ratesGroup.GET("/all", ratesHandler.GetAllRates)
    ratesGroup.GET("/status", ratesHandler.GetWorkerStatus)
}
```

### Routes with Middleware
```go
public := router.Group("/public")
{
    public.GET("/", handler.Index)
}

private := router.Group("/private")
private.Use(authMiddleware) // Apply middleware to group
{
    private.GET("/dashboard", handler.Dashboard)
    private.POST("/update", handler.Update)
}
```

## Common Gin Context Methods

### Request Data
| Method | Purpose |
|--------|---------|
| `c.Query(key)` | Get URL query parameter |
| `c.DefaultQuery(key, default)` | Get query parameter with default |
| `c.Param(key)` | Get URL path parameter |
| `c.PostForm(key)` | Get form field value |
| `c.DefaultPostForm(key, default)` | Get form field with default |
| `c.GetPostForm(key)` | Get form field with ok flag |
| `c.BindJSON(&obj)` | Parse JSON body into struct |
| `c.Request` | Access underlying *http.Request |

### Response Data
| Method | Purpose |
|--------|---------|
| `c.String(code, format, args...)` | Send string response |
| `c.JSON(code, obj)` | Send JSON response (auto-marshaled) |
| `c.HTML(code, template, data)` | Render HTML template |
| `c.XML(code, obj)` | Send XML response |
| `c.File(path)` | Send file |
| `c.Writer` | Access http.ResponseWriter for direct writing |

### Headers & Metadata
| Method | Purpose |
|--------|---------|
| `c.Header(key, value)` | Set response header |
| `c.SetCookie(...)` | Set cookie |
| `c.GetHeader(key)` | Get request header |
| `c.GetBool(key)` | Get context value as bool |
| `c.Set(key, value)` | Set context value (for middleware) |
| `c.Get(key)` | Get context value |

### Error Responses
| Method | Purpose |
|--------|---------|
| `c.AbortWithStatusJSON(code, obj)` | Abort with error JSON |
| `c.AbortWithError(code, err)` | Abort with error |
| `c.AbortWithStatus(code)` | Abort with status code |
| `c.Error(err)` | Add error to context |

## Converting Old net/http Code

### Old: Writing to ResponseWriter
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    fmt.Fprintf(w, `{"message":"hello"}`)
}
```

### New: Using Gin Context
```go
func Handler(c *gin.Context) {
    c.Header("Content-Type", "application/json")
    c.JSON(200, gin.H{"message": "hello"})
    // Or:
    c.String(200, `{"message":"hello"}`)
}
```

### Old: Reading Query Parameters
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
}
```

### New: Using Gin Context
```go
func Handler(c *gin.Context) {
    name := c.Query("name")
}
```

### Old: HTTP Error Response
```go
func Handler(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "Not found", http.StatusNotFound)
}
```

### New: Using Gin Context
```go
func Handler(c *gin.Context) {
    c.AbortWithStatusJSON(404, gin.H{"error": "Not found"})
}
```

## Testing with Gin

### Using net/http/httptest
```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    router := gin.Default()
    router.GET("/hello", handler.Hello)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/hello?q=World", nil)
    
    router.ServeHTTP(w, req)
    
    if w.Code != 200 {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Running the Server

```bash
# Start with default settings
go run main.go

# Or build and run binary
go build -o main
./main

# Server runs on http://localhost:3001
```

## Debugging Tips

1. **Check logs**: Gin's default middleware logs all requests
2. **Use structured responses**: Always use `c.JSON()` or `c.String()` for consistency
3. **Test routes**: Use curl or Postman to test endpoints
4. **Context values**: Use `c.Set()` and `c.Get()` to pass data through middleware/handlers

## References

- [Gin Official Documentation](https://gin-gonic.com/en/docs/)
- [Gin GitHub Repository](https://github.com/gin-gonic/gin)
- [HTTP Status Codes](https://httpwg.org/specs/rfc7231.html#status.codes)

