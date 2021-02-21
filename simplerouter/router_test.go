package simplerouter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(router *Router, method string, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func check(rr *httptest.ResponseRecorder, status int, body string) error {
	if rr.Code != status {
		return fmt.Errorf("handler returned wrong status code: got <%d> want <%d>", rr.Code, status)
	}

	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(body) {
		return fmt.Errorf("handler returned wrong body: got <%s> want <%s>", strings.TrimSpace(rr.Body.String()), strings.TrimSpace(body))
	}

	return nil
}

func TestHead(t *testing.T) {
	router := New()
	router.Head(`/head`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "head")
	})

	rr := do(router, http.MethodHead, "/head")
	if err := check(rr, http.StatusOK, "head"); err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	router := New()
	router.Connect(`/connect`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "connect")
	})

	rr := do(router, http.MethodConnect, "/connect")
	if err := check(rr, http.StatusOK, "connect"); err != nil {
		t.Error(err)
	}
}

func TestOptions(t *testing.T) {
	router := New()
	router.Options(`/options`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "options")
	})

	rr := do(router, http.MethodOptions, "/options")
	if err := check(rr, http.StatusOK, "options"); err != nil {
		t.Error(err)
	}
}

func TestTrace(t *testing.T) {
	router := New()
	router.Trace(`/trace`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "trace")
	})

	rr := do(router, http.MethodTrace, "/trace")
	if err := check(rr, http.StatusOK, "trace"); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	router := New()
	router.Get(`/get`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "get")
	})

	rr := do(router, http.MethodGet, "/get")
	if err := check(rr, http.StatusOK, "get"); err != nil {
		t.Error(err)
	}
}

func TestPost(t *testing.T) {
	router := New()
	router.Post(`/post`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "post")
	})

	rr := do(router, http.MethodPost, "/post")
	if err := check(rr, http.StatusOK, "post"); err != nil {
		t.Error(err)
	}
}

func TestPut(t *testing.T) {
	router := New()
	router.Put(`/put`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "put")
	})

	rr := do(router, http.MethodPut, "/put")
	if err := check(rr, http.StatusOK, "put"); err != nil {
		t.Error(err)
	}
}

func TestPatch(t *testing.T) {
	router := New()
	router.Patch(`/patch`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "patch")
	})

	rr := do(router, http.MethodPatch, "/patch")
	if err := check(rr, http.StatusOK, "patch"); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	router := New()
	router.Delete(`/delete`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "delete")
	})

	rr := do(router, http.MethodDelete, "/delete")
	if err := check(rr, http.StatusOK, "delete"); err != nil {
		t.Error(err)
	}
}

func TestDefaultNotFound(t *testing.T) {
	router := New()

	rr := do(router, http.MethodGet, "/ops")
	if err := check(rr, http.StatusNotFound, http.StatusText(http.StatusNotFound)); err != nil {
		t.Error(err)
	}
}

func TestCustomNotFound(t *testing.T) {
	router := New()
	router.NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "custom")
	}

	rr := do(router, http.MethodPost, "/ops")
	if err := check(rr, http.StatusOK, "custom"); err != nil {
		t.Error(err)
	}
}

func TestDefaultMethodNotAllowed(t *testing.T) {
	router := New()
	router.Get(`/ops`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ops")
	})

	rr := do(router, http.MethodPost, "/ops")
	if err := check(rr, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)); err != nil {
		t.Error(err)
	}
}

func TestCustomtMethodNotAllowed(t *testing.T) {
	router := New()
	router.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "custom")
	}

	router.Get(`/ops`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ops")
	})

	rr := do(router, http.MethodPost, "/ops")
	if err := check(rr, http.StatusOK, "custom"); err != nil {
		t.Error(err)
	}
}

func TestGetValidParams(t *testing.T) {
	router := New()
	router.Get(`/get/(?P<param1>\d+)/(?P<param2>\w+)`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := make(map[string]string)

		if value, ok := GetParam(r, "param1"); ok {
			data["param1"] = value
		}

		if value, ok := GetParam(r, "param2"); ok {
			data["param2"] = value
		}

		json.NewEncoder(w).Encode(data)
	})

	rr := do(router, http.MethodGet, "/get/123/abc")
	if err := check(rr, http.StatusOK, `{"param1":"123","param2":"abc"}`); err != nil {
		t.Error(err)
	}
}

func TestGetInvalidParams(t *testing.T) {
	router := New()
	router.Get(`/get/(?P<param1>\d+)/(?P<param2>\w+)`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := make(map[string]string)

		if value, ok := GetParam(r, "invalid1"); ok {
			data["param1"] = value
		}

		if value, ok := GetParam(r, "invalid2"); ok {
			data["param2"] = value
		}

		json.NewEncoder(w).Encode(data)
	})

	rr := do(router, http.MethodGet, "/get/123/abc")
	if err := check(rr, http.StatusOK, `{}`); err != nil {
		t.Error(err)
	}
}
