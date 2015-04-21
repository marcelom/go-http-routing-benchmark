// Copyright 2014 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/bmizerany/pat"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/rcrowley/go-tigertonic"
	goji "github.com/zenazn/goji/web"
)

type route struct {
	method string
	path   string
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

var nullLogger *log.Logger

func init() {
}

// Common
func httpHandlerFunc(w http.ResponseWriter, r *http.Request) {}

// goji
func gojiFuncWrite(c goji.C, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, c.URLParams["name"])
}

func loadGoji(routes []route) http.Handler {
	mux := goji.New()
	for _, route := range routes {
		switch route.method {
		case "GET":
			mux.Get(route.path, httpHandlerFunc)
		case "POST":
			mux.Post(route.path, httpHandlerFunc)
		case "PUT":
			mux.Put(route.path, httpHandlerFunc)
		case "PATCH":
			mux.Patch(route.path, httpHandlerFunc)
		case "DELETE":
			mux.Delete(route.path, httpHandlerFunc)
		default:
			panic("Unknown HTTP method: " + route.method)
		}
	}
	return mux
}

func loadGojiSingle(method, path string, handler interface{}) http.Handler {
	mux := goji.New()
	switch method {
	case "GET":
		mux.Get(path, handler)
	case "POST":
		mux.Post(path, handler)
	case "PUT":
		mux.Put(path, handler)
	case "PATCH":
		mux.Patch(path, handler)
	case "DELETE":
		mux.Delete(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return mux
}

// gorilla/mux
func gorillaHandlerWrite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	io.WriteString(w, params["name"])
}

func loadGorillaMux(routes []route) http.Handler {
	re := regexp.MustCompile(":([^/]*)")
	m := mux.NewRouter()
	for _, route := range routes {
		m.HandleFunc(
			re.ReplaceAllString(route.path, "{$1}"),
			httpHandlerFunc,
		).Methods(route.method)
	}
	return m
}

func loadGorillaMuxSingle(method, path string, handler http.HandlerFunc) http.Handler {
	m := mux.NewRouter()
	m.HandleFunc(path, handler).Methods(method)
	return m
}

// HttpRouter
func httpRouterHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {}

func httpRouterHandleWrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	io.WriteString(w, ps.ByName("name"))
}

func loadHttpRouter(routes []route) http.Handler {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.method, route.path, httpRouterHandle)
	}
	return router
}

func loadHttpRouterSingle(method, path string, handle httprouter.Handle) http.Handler {
	router := httprouter.New()
	router.Handle(method, path, handle)
	return router
}

// pat
func patHandlerWrite(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.URL.Query().Get(":name"))
}

func loadPat(routes []route) http.Handler {
	m := pat.New()
	for _, route := range routes {
		switch route.method {
		case "GET":
			m.Get(route.path, http.HandlerFunc(httpHandlerFunc))
		case "POST":
			m.Post(route.path, http.HandlerFunc(httpHandlerFunc))
		case "PUT":
			m.Put(route.path, http.HandlerFunc(httpHandlerFunc))
		case "DELETE":
			m.Del(route.path, http.HandlerFunc(httpHandlerFunc))
		default:
			panic("Unknow HTTP method: " + route.method)
		}
	}
	return m
}

func loadPatSingle(method, path string, handler http.Handler) http.Handler {
	m := pat.New()
	switch method {
	case "GET":
		m.Get(path, handler)
	case "POST":
		m.Post(path, handler)
	case "PUT":
		m.Put(path, handler)
	case "DELETE":
		m.Del(path, handler)
	default:
		panic("Unknow HTTP method: " + method)
	}
	return m
}

// Tiger Tonic
func tigerTonicHandlerWrite(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.URL.Query().Get("name"))
}

func loadTigerTonic(routes []route) http.Handler {
	re := regexp.MustCompile(":([^/]*)")
	mux := tigertonic.NewTrieServeMux()
	for _, route := range routes {
		mux.HandleFunc(route.method, re.ReplaceAllString(route.path, "{$1}"), httpHandlerFunc)
	}
	return mux
}

func loadTigerTonicSingle(method, path string, handler http.HandlerFunc) http.Handler {
	mux := tigertonic.NewTrieServeMux()
	mux.HandleFunc(method, path, handler)
	return mux
}

// Usage notice
func main() {
	fmt.Println("Usage: go test -bench=. -timeout=20m")
	os.Exit(1)
}
