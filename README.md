# [A2Y-D5L](https://github.com/a2y-d5l) / [serve](https://github.com/a2y-d5l/serve)

Go servers with a single function.

[![Go Reference](https://pkg.go.dev/badge/github.com/A2Y-D5L/serve.svg)](https://pkg.go.dev/github.com/A2Y-D5L/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/A2Y-D5L/serve)](https://goreportcard.com/report/github.com/A2Y-D5L/serve)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a2y-d5l/serve)](go.mod)
[![License](https://img.shields.io/github/license/a2y-d5l/serve)](LICENSE)
<!-- [![GitHub release (latest by date)](https://img.shields.io/github/v/release/a2y-d5l/serve)]() -->

## Usage

### Run an HTTP server

```go
package main

import (
    "github.com/A2Y-D5L/serve/httpserver"
)

func main() {
    if err := httpserver.Serve(
        context.Background(),
        []httpserver.Route{{
            Pattern: "GET /",
            Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Write([]byte("Hello, World!"))
            }),
        }},
    ); err != nil {
        log.Fatal(err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
