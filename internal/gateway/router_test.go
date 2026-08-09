package gateway

import "testing"

func TestRouterMatch(t *testing.T) {
	routes := []Route{
		{
			Path:       "/api/users",
			Target:     "http://localhost:4000",
			TargetPath: "/users",
		},
		{
			Path:       "/api/orders",
			Target:     "http://localhost:5000",
			TargetPath: "/orders",
		},
	}

	router, err := NewRouter(routes)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "users",
			path: "/api/users",
			want: "/api/users",
		},
		{
			name: "user by id",
			path: "/api/users/123",
			want: "/api/users",
		},
		{
			name: "orders",
			path: "/api/orders",
			want: "/api/orders",
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
			route := router.Match(tt.path)

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
