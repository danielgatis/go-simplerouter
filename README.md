# go-simplerouter

[![Go Report Card](https://goreportcard.com/badge/github.com/danielgatis/go-simplerouter?style=flat-square)](https://goreportcard.com/report/github.com/danielgatis/go-simplerouter)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/danielgatis/go-simplerouter/master/LICENSE)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/danielgatis/go-simplerouter)

This package is a simple request router.

### How to use

```bash
go get -u github.com/danielgatis/go-simplerouter
```

And then import the package in your code:

```go
import "github.com/danielgatis/go-simplerouter/simplerouter"
```

#### CRUD Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgatis/go-simplerouter/simplerouter"
	"github.com/spf13/cast"
)

var db map[int]*book = make(map[int]*book)

func genID() int {
	return int(time.Now().UnixNano() / int64(time.Millisecond))
}

type book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func index(w http.ResponseWriter, r *http.Request) {
	books := make([]*book, 0)

	for k := range db {
		books = append(books, db[k])
	}

	json.NewEncoder(w).Encode(books)
}

func create(w http.ResponseWriter, r *http.Request) {
	var newBook book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newBook.ID = genID()
	db[newBook.ID] = &newBook

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func show(w http.ResponseWriter, r *http.Request) {
	paramID, ok := simplerouter.GetParam(r, "id")
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id := cast.ToInt(paramID)
	book, ok := db[id]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

func update(w http.ResponseWriter, r *http.Request) {
	paramID, ok := simplerouter.GetParam(r, "id")
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id := cast.ToInt(paramID)
	oldBook, ok := db[id]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	var newBook book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	oldBook.Title = newBook.Title
	oldBook.Author = newBook.Author

	json.NewEncoder(w).Encode(oldBook)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	paramID, ok := simplerouter.GetParam(r, "id")
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id := cast.ToInt(paramID)
	delete(db, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := simplerouter.New()
	router.Get(`/books`, index)
	router.Post(`/books`, create)
	router.Get(`/books/(?P<id>\d+)`, show)
	router.Put(`/books/(?P<id>\d+)`, update)
	router.Delete(`/books/(?P<id>\d+)`, destroy)

	log.Fatal(http.ListenAndServe(":8080", router))
}

```

#### Middleware Example

```go
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/danielgatis/go-simplerouter/simplerouter"
	"github.com/justinas/alice"
)

func use(chain alice.Chain, fn http.HandlerFunc) http.HandlerFunc {
	return chain.ThenFunc(fn).ServeHTTP
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("log middleware")
		h.ServeHTTP(w, r)
	})
}

func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth middleware")
		h.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "index")
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "home")
}

func main() {
	public := alice.New(logMiddleware)
	private := alice.New(logMiddleware, authMiddleware)

	router := simplerouter.New()
	router.Get(`/`, use(public, index))
	router.Get(`/home`, use(private, home))

	log.Fatal(http.ListenAndServe(":8080", router))
}
```

#### Subdomain Example

```go
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/danielgatis/go-simplerouter/simplerouter"
)

type Routres map[string]http.Handler

func (routers Routres) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler := routers[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}

func index1(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "index1")
}

func index2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "index2")
}

func main() {
	router1 := simplerouter.New()
	router1.Get(`/`, index1)

	router2 := simplerouter.New()
	router2.Get(`/`, index2)

	routers := make(Routres)
	routers["site1.localhost:8080"] = router1
	routers["site2.localhost:8080"] = router2

	log.Fatal(http.ListenAndServe(":8080", routers))
}
```

### License

Copyright (c) 2021-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE)

### Buy me a coffee
Liked some of my work? Buy me a coffee (or more likely a beer)

<a href="https://www.buymeacoffee.com/danielgatis" target="_blank"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;"></a>
