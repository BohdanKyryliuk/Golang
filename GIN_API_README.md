# Golang Gin API - After Migration

This is a Go HTTP API application built with the **Gin Web Framework**. The application provides endpoints for basic HTTP operations and currency exchange rate tracking via external APIs.

## Quick Start

### Prerequisites
- Go 1.25 or higher
- Gin framework (automatically installed via go.mod)

### Build
```bash
cd /home/bkyryliuk/projects/Golang
go build -o main
```

### Run
```bash
./main
```

The server will start on `http://localhost:3001`

### Run Tests
```bash
go test ./... -v
```

## Project Structure

```
├── http/
│   └── handler/
│       ├── common.go       # Basic handlers (Hello, Counter) - Gin version
│       ├── currency.go     # Currency API handlers - Gin version
│       └── rates.go        # Cached rates handlers - Gin version
├── web/
│   └── web.go             # Server setup with Gin router
├── currencyapi/           # External API client
├── currency_converter/    # Currency conversion logic
├── worker/               # Background worker for rate updates
└── main.go              # Entry point
```

## API Endpoints

### Basic Routes

#### GET /
Display "Hello World" or greet a user via query parameter
```bash
curl http://localhost:3001/
curl http://localhost:3001/?q=John
```

#### GET /count
Display counter form (initial value: 1)
```bash
curl http://localhost:3001/count
```

#### POST /count
Increment counter
```bash
curl -X POST http://localhost:3001/count -d "counter=1"
```

### Currency Endpoints

#### GET /currency/status
Check if the currency API is operational
```bash
curl http://localhost:3001/currency/status
```

#### GET /currency/currencies
List all supported currencies
```bash
curl http://localhost:3001/currency/currencies
```

#### GET /currency/latest
Get latest exchange rates
```bash
# Get rates with default currency
curl http://localhost:3001/currency/latest?base=USD

# Get rates for specific currencies
curl "http://localhost:3001/currency/latest?base=USD&currencies=EUR,GBP,JPY"
```

### Rates Cache Endpoints

#### GET /rates
Get cached rates for a specific currency
```bash
curl "http://localhost:3001/rates?base=USD"
```

#### GET /rates/all
Get all cached rates
```bash
curl http://localhost:3001/rates/all
```

#### GET /rates/status
Get worker status
```bash
curl http://localhost:3001/rates/status
```

Response example:
```json
{
  "running": true,
  "currencies": ["USD", "EUR", "GBP"]
}
```

## Configuration

### Environment Variables
The application uses the following environment variables (optional):

- `CURRENCY_API_KEY` - API key for currency API service
- `CURRENCY_API_BASE_URL` - Base URL for currency API (if custom)

Load from `.env` file using:
```bash
export $(cat .env | xargs)
```

### Worker Configuration
Workers automatically fetch and cache exchange rates. Default configuration:
- Update interval: 5 minutes
- Currencies tracked: USD, EUR, GBP (configurable in code)

## Key Changes from net/http to Gin

### Handler Functions
```go
// Before (net/http)
func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, "...")
}

// After (Gin)
func Handler(c *gin.Context) {
    c.Header("Content-Type", "application/json")
    c.JSON(200, gin.H{"message": "..."})
}
```

### Error Handling
```go
// Before (net/http)
http.Error(w, "Not found", http.StatusNotFound)

// After (Gin)
c.AbortWithStatusJSON(404, gin.H{"error": "Not found"})
```

### Query Parameters
```go
// Before (net/http)
name := r.URL.Query().Get("name")

// After (Gin)
name := c.Query("name")
```

## Route Organization

The application uses Gin's route grouping feature for better organization:

```go
// Currency routes
router.Group("/currency")
  - GET /currency/status
  - GET /currency/currencies
  - GET /currency/latest

// Rates routes
router.Group("/rates")
  - GET /rates
  - GET /rates/all
  - GET /rates/status
```

This makes it easy to add middleware or rate limiting to entire groups of routes.

## Development Guide

### Adding a New Handler

1. Create handler in appropriate file (e.g., `http/handler/newfile.go`)
2. Use Gin context signature:
   ```go
   func MyHandler(c *gin.Context) {
       // Get request data
       param := c.Query("param")
       
       // Process
       result := process(param)
       
       // Send response
       c.JSON(200, result)
   }
   ```

3. Register in `web/web.go`:
   ```go
   router.GET("/my-path", handler.MyHandler)
   ```

### Adding Middleware

```go
// Global middleware
router.Use(someMiddleware)

// Group middleware
group := router.Group("/api")
group.Use(authMiddleware)
{
    group.GET("/protected", protectedHandler)
}
```

### Testing Handlers

Use Go's `net/http/httptest`:
```go
func TestHandler(t *testing.T) {
    router := gin.Default()
    router.GET("/test", handler.TestHandler)
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/test?param=value", nil)
    
    router.ServeHTTP(w, req)
    
    if w.Code != 200 {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

## Troubleshooting

### Port Already in Use
If port 3001 is already in use, modify `web/web.go`:
```go
router.Run(":8080")  // Use different port
```

### Currency API Not Working
Check that `CURRENCY_API_KEY` environment variable is set:
```bash
export CURRENCY_API_KEY="your-api-key"
./main
```

### Worker Not Starting
Check logs for error messages. Workers require a valid currency API client.

## Documentation

- `GIN_MIGRATION_SUMMARY.md` - Detailed migration guide
- `GIN_QUICK_REFERENCE.md` - Quick reference for Gin patterns
- `MIGRATION_CHANGES.md` - Summary of all file changes
- [Gin Official Docs](https://gin-gonic.com/en/docs/)

## Performance

- Binary size: ~13 MB (including Gin and dependencies)
- Memory usage: Optimized with Gin's efficient routing
- Request handling: Fast with O(1) route lookup (via httprouter)

## Graceful Shutdown

The server handles signals (SIGINT, SIGTERM) gracefully:
1. Stops accepting new requests
2. Waits for in-flight requests to complete
3. Stops background workers
4. Exits cleanly

## License

Same as original project

## Support

Refer to the migration documents in the project root for detailed information about working with the Gin version of the API.

