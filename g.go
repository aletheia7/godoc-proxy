// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

var parse_http = regexp.MustCompile(`(?:/http(?:s*):)(.*$)`)

const site_css = `
html { background-color: whitesmoke; }
body { background-color: white; }
h4 { margin-top: 20px; }
.container { width: 100%; } 
/* .container { max-width: 728px; } */

#x-projnav {
    min-height: 20px;
    margin-bottom: 20px;
    background-color: #eee;
    padding: 9px;
    border-radius: 3px;
}

#x-footer {
    padding-top: 14px;
    padding-bottom: 15px;
    margin-top: 5px;
    background-color: #eee;
    border-top-style: solid;
    border-top-width: 1px;

}

.highlighted {
    background-color: #FDFF9E;
}

#x-pkginfo {
    margin-top: 25px;
    border-top: 1px solid #ccc;
    padding-top: 20px;
    margin-bottom: 15px;
}

code {
    background-color: inherit;
    border: none;
    color: #222;
    padding: 0;
}

pre {
	width: 100% !important;
    color: #222;
    white-space: pre;
	overflow: visible;
    word-break: normal;
    word-wrap: normal;
}

.funcdecl > pre {
    white-space: pre-wrap;
    word-break: break-all;
    word-wrap: break-word;
}

pre .com {
    color: #006600;
}

.decl {
    position: relative;
}

.decl > a {
    position: absolute;
    top: 0px;
    right: 0px;
    display: none;
    border: 1px solid #ccc;
    border-top-right-radius: 4px;
    border-bottom-left-radius: 4px;
    padding-left: 4px;
    padding-right: 4px;
}

.decl > a:hover {
    background-color: white;
    text-decoration: none;
}

.decl:hover > a {
    display: block;
}

a, .navbar-default .navbar-brand {
    color: #375eab;
}

.navbar-default, #x-footer {
    background-color: hsl(209, 51%, 92%);
    border-color: hsl(209, 51%, 88%);
}

.navbar-default .navbar-nav > .active > a,
.navbar-default .navbar-nav > .active > a:hover,
.navbar-default .navbar-nav > .active > a:focus {
    background-color: hsl(209, 51%, 88%);
}

.navbar-default .navbar-nav > li > a:hover,
.navbar-default .navbar-nav > li > a:focus {
    color: #000;
}

.panel-default > .panel-heading {
    color: #333;
    background-color: transparent;
}

a.permalink {
    display: none;
}

a.uses {
    display: none;
    color: #666;
    font-size: 0.8em;
}

h1:hover .permalink, h2:hover .permalink, h3:hover .permalink, h4:hover .permalink, h5:hover .permalink, h6:hover .permalink, h1:hover .uses, h2:hover .uses, h3:hover .uses, h4:hover .uses, h5:hover .uses, h6:hover .uses {
    display: inline;
}

@media (max-width : 768px) {
    .form-control {
        font-size:16px;
    }
}

.synopsis {
  opacity: 0.87;
}

.additional-info {
    display: block;
    opacity: 0.54;
    text-transform: uppercase;
    font-size: 0.75em;
}
`

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
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
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		},
	}

	http.Handle("/", p)
	http.HandleFunc("/-/site.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/css; charset=utf-8")
		io.WriteString(w, site_css)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
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

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
