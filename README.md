# [A2Y-D5L](https://github.com/a2y-d5l) / [serve](https://github.com/a2y-d5l/serve)

Go servers with a single function.

[![Go Reference](https://pkg.go.dev/badge/github.com/A2Y-D5L/serve.svg)](https://pkg.go.dev/github.com/A2Y-D5L/serve)
[![Go Report Card](https://goreportcard.com/badge/github.com/A2Y-D5L/serve)](https://goreportcard.com/report/github.com/A2Y-D5L/serve)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a2y-d5l/serve)](go.mod)
[![License](https://img.shields.io/github/license/a2y-d5l/serve)](LICENSE)
<!-- [![GitHub release (latest by date)](https://img.shields.io/github/v/release/a2y-d5l/serve)]() -->

## Install

```bash
go get -u github.com/A2Y-D5L/serve
```

## Usage

### HTTP Server

#### Run an HTTP server

```go
package main

import (
    "context"
    "log/slog"
    "net/http"

    "github.com/A2Y-D5L/serve/httpserver"
)

func handleGETRoot(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    if err := httpserver.Serve(
        context.Background(),
        httpserver.Routes(
            httpserver.Route{
                Pattern: "GET /",
                Handler: http.HandlerFunc(handleGETRoot),
            },
        ),
    ); err != nil {
        slog.Error("httpserver.Serve error:" + err)
    }
}
```

#### Configuring the HTTP Server

The `httpserver` package exposes several functional options for common server
configurations. If you need further customization, you can
define your own options. For example:

```go
package main

import (
    "context"
    "log/slog"
    "net/http"

    "github.com/A2Y-D5L/serve/httpserver"
)

func DisableGeneralOptionsHandler() httpserver.Option {
    return func(srv *http.Server) {
        srv.DisableGeneralOptionsHandler = true
    }
}

func main() {
    if err := httpserver.Serve(
        context.Background(),
        httpserver.Address(":8080"),    // provided by the httpserver package
        DisableGeneralOptionsHandler(), // custom option
    ); err != nil {
        slog.Error("httpserver.Serve error:" + err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
