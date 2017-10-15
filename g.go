// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.

//go:generate go run cmd/gen_css/g.go

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aletheia7/gogroup"
	"gogitver"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"runtime"
	"strings"
)

var parse_http = regexp.MustCompile(`(?:/http(?:s*):)(.*$)`)

func main() {
	setup_log()
	phttp := flag.String("http", "127.0.0.1:80", "address:port")
	ver := flag.Bool("v", false, "version")
	gver := flag.Bool("gv", false, "go version")
	flag.Parse()
	switch {
	case *ver:
		log.Println("version:", gogitver.Git())
		return
	case *gver:
		log.Println("go version:", runtime.Version())
		return
	}
	p := new(Proxy)
	host := "godoc.org"
	u, err := url.Parse(fmt.Sprintf("https://%v/", host))
	if err != nil {
		log.Printf("Error parsing URL")
	}
	targetQuery := u.RawQuery
	p.proxy = &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = host
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.URL.Path = singleJoiningSlash(u.Path, req.URL.Path)
			req.URL.Path = parse_http.ReplaceAllString(req.URL.Path, `$1`)
			// if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery //+ req.URL.RawQuery
			// } else {
			// req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			// }
		},
	}
	srv := &http.Server{
		Addr:    *phttp,
		Handler: p,
	}
	gg := gogroup.New(nil)
	gg.Add_signals(gogroup.Unix)
	go func() {
		k := gg.Register()
		defer gg.Unregister(k)
		log.Println("listen:", *phttp)
		srv.ListenAndServe()
	}()
	defer log.Println("stopped")
	defer gg.Wait()
	<-gg.Ctx.Done()
	srv.Close()
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

type Proxy struct {
	proxy *httputil.ReverseProxy
}

var style_css_r = bytes.NewReader([]byte(style_css))

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == `/-/site.css` {
		http.ServeContent(w, r, "site.css", style_css_time, style_css_r)
		return
	}
	p.proxy.ServeHTTP(w, r)
}
