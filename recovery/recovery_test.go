package recovery

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareRecoversFromPanic(t *testing.T) {
	var output bytes.Buffer

	logger := slog.New(
		slog.NewTextHandler(
			&output,
			nil,
		),
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	handler := Middleware(
		logger,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/test",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	if !strings.Contains(
		rec.Body.String(),
		http.StatusText(http.StatusInternalServerError),
	) {
		t.Fatalf(
			"expected response body to contain %q, got %q",
			http.StatusText(http.StatusInternalServerError),
			rec.Body.String(),
		)
	}

	logOutput := output.String()

	expectedValues := []string{
		`msg="panic recovered"`,
		"error=boom",
		"method=GET",
		"path=/test",
	}

	for _, expected := range expectedValues {
		if !strings.Contains(logOutput, expected) {
			t.Fatalf(
				"expected log output to contain %q, got %q",
				expected,
				logOutput,
			)
		}
	}
}

func TestMiddlewareCallsNextHandlerNormally(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(
			&bytes.Buffer{},
			nil,
		),
	)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	handler := Middleware(
		logger,
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}
}
