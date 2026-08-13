package csrf

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProtectAllowsSafeRequest(t *testing.T) {
	key := []byte("01234567890123456789012345678901")

	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := Protect(
		Config{
			Key:    key,
			Secure: false,
		},
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestProtectRejectsUnsafeRequestWithoutToken(t *testing.T) {
	key := []byte("01234567890123456789012345678901")

	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := Protect(
		Config{
			Key:    key,
			Secure: false,
		},
		next,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/submit",
		strings.NewReader("name=Daniel"),
	)

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if nextCalled {
		t.Fatal("expected request without CSRF token to be rejected")
	}

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestTokenReturnsToken(t *testing.T) {
	key := []byte("01234567890123456789012345678901")

	var token string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token = Token(r)
		w.WriteHeader(http.StatusOK)
	})

	handler := Protect(
		Config{
			Key:    key,
			Secure: false,
		},
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/form",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if token == "" {
		t.Fatal("expected CSRF token to be generated")
	}
}

func TestTemplateFieldReturnsHiddenInput(t *testing.T) {
	key := []byte("01234567890123456789012345678901")

	var field string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		field = string(TemplateField(r))
		w.WriteHeader(http.StatusOK)
	})

	handler := Protect(
		Config{
			Key:    key,
			Secure: false,
		},
		next,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/form",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !strings.Contains(field, `type="hidden"`) {
		t.Fatalf(
			"expected hidden CSRF input, got %q",
			field,
		)
	}

	if !strings.Contains(field, "gorilla.csrf.Token") {
		t.Fatalf(
			"expected CSRF field name, got %q",
			field,
		)
	}
}
