package router

import "net/http"

type Router struct {
	routes []RouteEntry
}

func (rtr *Router) Route(method, path string, handlerFunc http.HandlerFunc) {
	e := RouteEntry{
		Method:  method,
		Path:    path,
		Handler: handlerFunc,
	}
	rtr.routes = append(rtr.routes, e)
}

func (rtr *Router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	for _, e := range rtr.routes {
		match := e.Match(r)
		if !match {
			continue
		}

		// We have a match! Call the handler, and return
		e.Handler.ServeHTTP(w, r)
		return
	}

	// return 404 for every request.
	http.NotFound(w, r)
}

type RouteEntry struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func (re *RouteEntry) Match(r *http.Request) bool {
	if r.Method != re.Method {
		return false // Method mismatch
	}

	if r.URL.Path != re.Path {
		return false // Path mismatch
	}

	return true
}
