package csrf

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProtectAllowsSameOriginPost(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler, err := Protect(Config{}, next)
	if err != nil {
		t.Fatalf("create protection: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/login",
		nil,
	)

	req.Header.Set("Origin", "http://example.com")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected same-origin request to be allowed")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestProtectRejectsCrossOriginPost(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		nextCalled = true
	})

	handler, err := Protect(Config{}, next)
	if err != nil {
		t.Fatalf("create protection: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/login",
		nil,
	)

	req.Header.Set(
		"Origin",
		"http://evil.example",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if nextCalled {
		t.Fatal("expected cross-origin request to be rejected")
	}

	if rec.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			rec.Code,
		)
	}
}

func TestProtectAllowsSafeCrossOriginGet(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler, err := Protect(Config{}, next)
	if err != nil {
		t.Fatalf("create protection: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.com/",
		nil,
	)

	req.Header.Set(
		"Origin",
		"http://evil.example",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected safe GET request to be allowed")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestProtectAllowsTrustedOrigin(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler, err := Protect(
		Config{
			TrustedOrigins: []string{
				"http://localhost:8080",
			},
		},
		next,
	)
	if err != nil {
		t.Fatalf("create protection: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"http://example.com/login",
		nil,
	)

	req.Header.Set(
		"Origin",
		"http://localhost:8080",
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("expected trusted origin to be allowed")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
