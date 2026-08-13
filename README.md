# go-web-security

Reusable security middleware for Go web applications.

`go-web-security` provides a small set of security-focused packages that can be added to server-rendered Go applications without coupling them to authentication, sessions, or application-specific business logic.

It is designed to work alongside [`go-web-core`](https://github.com/danieljmanningdev/go-web-core), but can also be used independently.

## Features

* CSRF protection using `github.com/gorilla/csrf`
* Secure HTTP response headers
* Panic recovery middleware
* HTMX-compatible CSRF usage
* Configurable secure-cookie behaviour
* Automated tests for core security behaviour

## Installation

```bash
go get github.com/danieljmanningdev/go-web-security@v0.1.0
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

Provides a small wrapper around `github.com/gorilla/csrf`.

```go
handler := csrf.Protect(
	csrf.Config{
		Key:    csrfKey,
		Secure: true,
	},
	mux,
)
```

The CSRF key should be a persistent 32-byte secret and should not be hardcoded into application source code.

For example:

```go
csrfKey := []byte(os.Getenv("CSRF_KEY"))
```

In production:

```go
csrf.Config{
	Key:    csrfKey,
	Secure: true,
}
```

For local development over plain HTTP:

```go
csrf.Config{
	Key:    csrfKey,
	Secure: false,
}
```

#### Tokens

Retrieve the current request token:

```go
token := csrf.Token(r)
```

Retrieve a ready-to-render hidden form field:

```go
field := csrf.TemplateField(r)
```

This can be passed into a Go template and included inside state-changing forms.

For example:

```html
<form method="post" action="/account">
	{{ .CSRFField }}

	<button type="submit">Save</button>
</form>
```

Unsafe requests such as `POST`, `PUT`, `PATCH`, and `DELETE` are rejected when they do not contain a valid CSRF token.

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	})

	handler := headers.Secure(mux)

	handler = recovery.Middleware(
		logger,
		handler,
	)

	handler = csrf.Protect(
		csrf.Config{
			Key:    []byte(os.Getenv("CSRF_KEY")),
			Secure: true,
		},
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
handler := headers.Secure(mux)

handler = recovery.Middleware(
	logger,
	handler,
)

handler = csrf.Protect(
	csrf.Config{
		Key:    csrfKey,
		Secure: true,
	},
	handler,
)
```

The final `handler` is then passed to the HTTP server.

The exact middleware stack can vary depending on the application.

## Design philosophy

This project is deliberately small.

It provides reusable security primitives without trying to become a complete authentication or security framework.

Authentication, users, sessions, password management, permissions, application-specific authorization, and business rules should live in separate packages or applications.

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

Check whitespace errors:

```bash
git diff --check
```

## Status

The project is currently in early development.

Until the API stabilises, releases should be considered pre-`v1.0.0` and may contain breaking changes.

Current release:

```text
v0.1.0
```

## License

See [LICENSE](LICENSE).
