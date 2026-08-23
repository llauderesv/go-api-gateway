package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/llauderesv/go-api-gateway/internal/config"
)

// Eventually Route will probably contain more than just a target URL, for example
// Methods, Timeout, RateLimit, Middleware
type Route = config.Route // Convert to alias

type proxyRoute struct {
	route Route
	proxy *httputil.ReverseProxy
}

type Router struct {
	routes []proxyRoute
}

func NewRouter(routes []Route) (*Router, error) {
	router := &Router{}

	for _, route := range routes {
		target, err := url.Parse(route.Target)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid target %q: %w",
				route.Target,
				err,
			)
		}
		// Each route in the slice array, we create a dedicated ReverseProxy instance
		proxy := &httputil.ReverseProxy{
			Rewrite: func(req *httputil.ProxyRequest) {
				req.SetURL(target)

				incomingPath := req.In.URL.Path
				fmt.Printf("incoming Path %s\n", incomingPath)
				remainingPath := strings.TrimPrefix(
					incomingPath,
					route.Path,
				)

				fmt.Printf("remaining Path %s\n", remainingPath)
				req.Out.URL.Path = route.TargetPath + remainingPath
			},
		}

		router.routes = append(router.routes, proxyRoute{
			route: route,
			proxy: proxy,
		})
	}

	return router, nil
}

// func NewRouter(routes []Route) *Router {
// 	return &Router{
// 		routes: routes,
// 	}
// }

func (r *Router) MatchPath(path string) []*proxyRoute {
	var matches []*proxyRoute
	for i := range r.routes {
		route := &r.routes[i]

		if matchesPath(route.route.Path, path) {
			matches = append(matches, route)
		}
	}

	return matches
}

func (r *Router) Match(method, path string) *proxyRoute {
	for i := range r.routes {
		route := &r.routes[i]
		if !matchesPath(route.route.Path, path) {
			continue
		}

		if !matchesMethod(route.route.Methods, method) {
			continue
		}

		return route
	}

	return nil
}

func (r *Router) AllowedMethods(path string) []string {
	var methods []string

	for i := range r.routes {
		route := &r.routes[i]

		if !matchesPath(route.route.Path, path) {
			continue
		}

		methods = append(methods, route.route.Methods...)
	}

	return methods
}

func matchesPath(routePath, requestPath string) bool {
	return requestPath == routePath || strings.HasPrefix(requestPath, routePath+"/")
}

func matchesMethod(methods []string, method string) bool {
	for _, allowed := range methods {
		if strings.EqualFold(allowed, method) {
			return true
		}
	}

	return false
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route := r.Match(req.Method, req.URL.Path)

	if route != nil {
		route.proxy.ServeHTTP(w, req)
		return
	}

	allowedMethods := r.AllowedMethods(req.URL.Path)

	if len(allowedMethods) == 0 {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	http.Error(
		w,
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed,
	)
}

// func (r *Router) Match(path string) *Route {
// 	for i, _ := range r.routes {
// 		route := &r.routes[i]

// 		if path == route.Path || strings.HasPrefix(path, route.Path+"/") {
// 			return route
// 		}
// 	}

// 	return nil
// }
