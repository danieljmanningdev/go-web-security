# go-web-security

[![CI](https://github.com/danieljmanningdev/go-web-security/actions/workflows/ci.yml/badge.svg)](https://github.com/danieljmanningdev/go-web-security/actions/workflows/ci.yml)

Reusable security middleware for Go web applications.

`go-web-security` provides a small set of security-focused packages for server-rendered Go applications without coupling them to authentication, sessions, or application-specific business logic.

It is designed to work alongside [`go-web-core`](https://github.com/danieljmanningdev/go-web-core), but can also be used independently.

## Features

* Cross-origin request protection using Go's standard library
* Secure HTTP response headers
* Panic recovery middleware
* Configurable trusted origins
* No application-managed CSRF secrets or tokens
* Automated tests for core security behaviour
* Compatible with Go's vulnerability tooling

## Installation

```bash
go get github.com/danieljmanningdev/go-web-security@v0.2.0
```

## Packages

### `headers`

Adds a small set of defensive HTTP response headers.

```go
handler := headers.Secure(mux)
```

Currently sets:

```text
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

Example:

```go
mux := http.NewServeMux()

handler := headers.Secure(mux)

if err := http.ListenAndServe(":8080", handler); err != nil {
	log.Fatal(err)
}
```

---

### `recovery`

Recovers from panics in downstream HTTP handlers, logs the failure, and returns an HTTP `500 Internal Server Error` response.

```go
handler := recovery.Middleware(
	logger,
	mux,
)
```

Logged request information includes:

* HTTP method
* request path
* recovered panic value

Example:

```go
logger := slog.New(
	slog.NewTextHandler(os.Stdout, nil),
)

handler := recovery.Middleware(
	logger,
	mux,
)
```

---

### `csrf`

Provides a small wrapper around Go's standard-library `http.CrossOriginProtection`.

The protection rejects unsafe cross-origin browser requests while allowing same-origin requests and safe methods such as `GET`, `HEAD`, and `OPTIONS`.

Unlike the previous implementation, applications do not need to manage:

```text
CSRF secret keys
CSRF cookies
hidden CSRF form fields
request tokens
```

Basic usage:

```go
handler, err := csrf.Protect(
	csrf.Config{},
	mux,
)
if err != nil {
	log.Fatal(err)
}
```

#### Trusted origins

Additional origins can be explicitly trusted when required:

```go
handler, err := csrf.Protect(
	csrf.Config{
		TrustedOrigins: []string{
			"http://localhost:8080",
		},
	},
	mux,
)
if err != nil {
	log.Fatal(err)
}
```

Trusted origins should only contain origins that the application intentionally allows to make unsafe requests.

For example, a local development configuration might use:

```go
csrfConfig := csrf.Config{
	TrustedOrigins: []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
	},
}
```

Production applications should configure only the origins they actually require.

#### Forms

No special hidden CSRF field is required.

A normal server-rendered form can remain:

```html
<form method="post" action="/account">
	<button type="submit">Save</button>
</form>
```

Cross-origin protection is applied by middleware before unsafe requests reach the application handler.

## Example

A minimal application using all three packages could look like:

```go
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/danieljmanningdev/go-web-security/csrf"
	"github.com/danieljmanningdev/go-web-security/headers"
	"github.com/danieljmanningdev/go-web-security/recovery"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, nil),
	)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		_, _ = w.Write([]byte("hello"))
	})

	var handler http.Handler = mux

	csrfHandler, err := csrf.Protect(
		csrf.Config{},
		handler,
	)
	if err != nil {
		log.Fatal(err)
	}

	handler = csrfHandler

	handler = headers.Secure(handler)

	handler = recovery.Middleware(
		logger,
		handler,
	)

	if err := http.ListenAndServe(
		":8080",
		handler,
	); err != nil {
		log.Fatal(err)
	}
}
```

## Middleware order

Middleware wraps from the outside in.

For example:

```go
var handler http.Handler = mux

csrfHandler, err := csrf.Protect(
	csrf.Config{},
	handler,
)
if err != nil {
	log.Fatal(err)
}

handler = csrfHandler
handler = headers.Secure(handler)

handler = recovery.Middleware(
	logger,
	handler,
)
```

The final `handler` is passed to the HTTP server.

The exact middleware stack can vary depending on the application.

## Design philosophy

This project is deliberately small.

It provides reusable security primitives without trying to become a complete authentication or security framework.

Authentication, users, sessions, password management, permissions, application-specific authorization, and business rules should live in separate packages or applications.

Where suitable security functionality exists in Go's standard library, this project prefers building on that rather than maintaining unnecessary custom security implementations.

The goal is simple:

> Give new Go web projects a small, tested security baseline without rebuilding the same middleware each time.

## Development

Format the code:

```bash
gofmt -w .
```

Run the tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Check for known vulnerabilities:

```bash
govulncheck ./...
```

Check whitespace errors:

```bash
git diff --check
```

A clean release should pass all four checks.

## Status

The project is currently in early development.

Until the API stabilises, releases should be considered pre-`v1.0.0` and may contain breaking changes.

Current release:

```text
v0.2.0
```

## License

See [LICENSE](LICENSE).
