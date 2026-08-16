package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterMatch(t *testing.T) {
	routes := []Route{
		{
			Path:       "/api/users",
			Methods:    []string{"GET", "POST"},
			Target:     "http://localhost:4000",
			TargetPath: "/users",
		},
		{
			Path:       "/api/orders",
			Methods:    []string{"GET", "POST"},
			Target:     "http://localhost:5000",
			TargetPath: "/orders",
		},
	}

	router, err := NewRouter(routes)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	tests := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{
			name:   "users",
			method: "POST",
			path:   "/api/users",
			want:   "/api/users",
		},
		{
			name:   "user by id",
			method: "GET",
			path:   "/api/users/123",
			want:   "/api/users",
		},
		{
			name:   "Method not allowed",
			method: "DELETE",
			path:   "/api/users/123",
			want:   "",
		},
		{
			name:   "orders",
			method: "GET",
			path:   "/api/orders",
			want:   "/api/orders",
		},
		{
			name: "unknown route",
			path: "/api/products",
			want: "",
		},
		{
			name: "similar but invalid route",
			path: "/api/usersabc",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := router.Match(tt.method, tt.path)

			if tt.want == "" {
				if route != nil {
					t.Fatalf("expected no route, got %q", route.route.Path)
				}

				return
			}

			if route == nil {
				t.Fatalf("expected route %q, got nil", tt.want)
			}

			if route.route.Path != tt.want {
				t.Fatalf(
					"expected route %q, got %q",
					tt.want,
					route.route.Path,
				)
			}
		})
	}
}

func TestRouterProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123" {
			t.Errorf("expected path /users/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))

	defer backend.Close()

	routes := []Route{
		{
			Path:       "/api/users",
			Target:     backend.URL,
			TargetPath: "/users",
		},
	}

	router, err := NewRouter(routes)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"http://gateway.local/api/users/123",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	if recorder.Body.String() != "backend response" {
		t.Fatalf(
			"expected response %q, got %q",
			"backend response",
			recorder.Body.String(),
		)
	}
}
