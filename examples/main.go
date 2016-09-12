package main

import (
	"fmt"
	"net/http"

	filter "github.com/connect-asia/goji-before-filter"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := pat.Param(ctx, "name")
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func main() {
	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/hello/:name"), filter.Before(before).On(hello))
	http.ListenAndServe("localhost:8000", mux)
}

type BeforeMiddleware struct {
}

func (m *BeforeMiddleware) Handle(c context.Context, res http.ResponseWriter, req *http.Request) (context.Context, http.ResponseWriter, *http.Request, bool) {
	fmt.Fprintln(res, "Good morning!")
	return c, res, req, false
}

func before() filter.Handler {
	return &BeforeMiddleware{}
}
