# URL Shortener Service

A simple URL shortening service that converts long URLs into short ones and redirects to the original URL when accessing the shortened URL.

## Features
- Context support (with cancellation handling)
- Structured logging (INFO, WARN, ERROR)
- Thread-safe in-memory store
- Detailed error handling

## How to Start the Service

1. Run the following command in the project root directory:

```bash
go run main.go
```

The server will start on port 8080.

## Usage Instructions

### 1. Shorten URL

You can convert a long URL into a shortened URL using the following curl command:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url/that/needs/shortening"}'
```

On success, you will receive a response like this:

```json
{
  "short_url": "http://localhost:8080/Ab3Cd9"
}
```

### 2. Verify Redirect

When you access the shortened URL, you will be redirected to the original URL:

```bash
curl -i http://localhost:8080/Ab3Cd9
```

The response will include a 302 status code and a Location header pointing to the original URL.

### Error Cases

1. For non-existent codes:
```bash
curl -i http://localhost:8080/nonexistent
```
Returns 404 Not Found and outputs a WARN level log.

2. For invalid requests:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'
```
Returns 400 Bad Request and outputs a WARN level log.

3. Context cancellation:
When a request is cancelled (e.g., timeout), an ERROR level log is output and an appropriate error response is returned.

### Log Levels
- INFO: Normal operations (URL registration, redirects, etc.)
- WARN: Invalid requests, non-existent URLs, etc.
- ERROR: Internal errors, context cancellations, etc.

## Running Tests

Run unit tests with the following command:

```bash
go test -v
```

To run tests including race condition checks:

```bash
go test -race -v
```

### Test Cases
1. URL Shortening Functionality
   - Normal URL registration and short URL generation
   - Invalid request handling

2. Redirect Functionality
   - Normal redirect operation
   - Non-existent code handling

3. Context and Cancellation Handling
   - Context cancellation behavior
   - Timeout handling

4. Log Output
   - Verification of each log level (INFO/WARN/ERROR)
   - Log message content validation
