package main

//go:generate go run gen.go

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
	"text/template"
)

func main() {
	f, err := os.Open("sws.min.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	compressor := gzip.NewWriter(encoder)
	defer compressor.Close()
	defer encoder.Close()

	_, err = io.Copy(compressor, f)

	tmpl, err := template.New("counter").Parse(tmplData)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(os.Stdout, struct{ B64 string }{buf.String()})
}

var tmplData = `package main

// Automatically generated, don't bother editing.

import (
	"bytes"
	"strings"
	"compress/gzip"
	"encoding/base64"
	"io"
)

const data = "{{ .B64 }}"

func GetCounter() []byte {
	b := GetCounterGzipped()
	decompressor, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return panic(err)
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, decompressor)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func GetCounterGzipped() []byte {
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	out, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return out
}
`
