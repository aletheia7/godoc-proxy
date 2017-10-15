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
	fp, err := os.Open("style.css")
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
	override, err := ioutil.ReadFile("override.css")
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

var go_file = template.Must(template.New("").Parse(`
package main

import (
	"time"
)

var style_css_time, _ = time.Parse(time.RFC3339Nano, "{{.modtime}}")

const style_css = {{.delim}}{{.css}}{{.override_css}}{{.delim}}`))
