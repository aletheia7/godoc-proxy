// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"
)

const fn = `g_style.go`

func main() {
	fp, err := os.Open("web/style.css")
	if err != nil {
		log.Fatalln(err)
	}
	defer fp.Close()
	st, err := fp.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	css, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatalln(err)
	}
	override, err := ioutil.ReadFile("web/override.css")
	if err != nil {
		log.Fatalln(err)
	}
	var f bytes.Buffer
	go_file.Execute(&f, map[string]interface{}{
		"delim":        "`",
		"css":          string(css),
		"override_css": string(override),
		"modtime":      st.ModTime().Format(time.RFC3339Nano),
	})
	if err := ioutil.WriteFile(fn, f.Bytes(), 0666); err != nil {
		log.Println(err)
	}
}

var go_file = template.Must(template.New("").Parse(`// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.
package main

import (
	"time"
)

var style_css_time, _ = time.Parse(time.RFC3339Nano, "{{.modtime}}")

const style_css = {{.delim}}{{.css}}{{.override_css}}{{.delim}}`))
