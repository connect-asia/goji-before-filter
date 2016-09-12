/*
The MIT License (MIT)

Copyright (c) 2016 kunihiko-t.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package goji_before_filter

import (
	"golang.org/x/net/context"
	"net/http"
)

type Handlers struct {
	BeforeMiddlewares []interface{}
}

func Before(middlewares ...interface{}) *Handlers {
	h := &Handlers{BeforeMiddlewares: middlewares}
	return h
}

func (h *Handlers) Before(middlewares ...interface{}) *Handlers {
	h.BeforeMiddlewares = middlewares
	return h
}

func (hs *Handlers) On(h func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
	f := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		halt := false
		if hs.BeforeMiddlewares != nil {
			ctx, w, r, halt = applyMiddlewares(ctx, w, r, hs.BeforeMiddlewares)
		}

		if halt {
			return
		} else {
			h(ctx, w, r)
		}

	}
	return f
}

func applyMiddlewares(ctx context.Context, w http.ResponseWriter, r *http.Request, middlewares []interface{}) (context.Context, http.ResponseWriter, *http.Request, bool) {
	halt := false
	for _, mw := range middlewares {
		switch t := mw.(type) {
		case func() Handler:
			mw := t()
			ctx, w, r, halt = mw.Handle(ctx, w, r)
		}
		if halt {
			break
		}
	}
	return ctx, w, r, halt
}

type Handler interface {
	Handle(c context.Context, res http.ResponseWriter, req *http.Request) (context.Context, http.ResponseWriter, *http.Request, bool)
}
