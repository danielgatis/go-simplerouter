package simplerouter

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type contextKey string

const paramsKey = contextKey("urlParams")

type route struct {
	method string
	regex  *regexp.Regexp
	fn     http.HandlerFunc
}

type Router struct {
	Routes                  []*route
	NotFoundHandler         http.HandlerFunc
	MethodNotAllowedHandler http.HandlerFunc
}

func New() *Router {
	return new(Router)
}

func (router *Router) add(method, pattern string, fn http.HandlerFunc) {
	router.Routes = append(router.Routes, &route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
		fn,
	})
}

func (router *Router) Head(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodHead, pattern, fn)
}

func (router *Router) Connect(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodConnect, pattern, fn)
}

func (router *Router) Options(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodOptions, pattern, fn)
}

func (router *Router) Trace(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodTrace, pattern, fn)
}

func (router *Router) Get(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodGet, pattern, fn)
}

func (router *Router) Post(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodPost, pattern, fn)
}

func (router *Router) Put(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodPut, pattern, fn)
}

func (router *Router) Patch(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodPatch, pattern, fn)
}

func (router *Router) Delete(pattern string, fn http.HandlerFunc) {
	router.add(http.MethodDelete, pattern, fn)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allow := make([]string, 0)

	for _, route := range router.Routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			params := make(map[string]string)

			for i, name := range route.regex.SubexpNames() {
				if i != 0 && name != "" {
					params[name] = matches[i]
				}
			}

			route.fn(w, r.WithContext(context.WithValue(r.Context(), paramsKey, params)))
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))

		if router.MethodNotAllowedHandler == nil {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		} else {
			router.MethodNotAllowedHandler(w, r)
		}

		return
	}

	if router.NotFoundHandler == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		router.NotFoundHandler(w, r)
	}
}

func GetParam(r *http.Request, name string) (string, bool) {
	params := r.Context().Value(paramsKey).(map[string]string)

	if param, ok := params[name]; ok {
		return param, true
	}

	return "", false
}
