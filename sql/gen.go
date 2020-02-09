package main

//go:generate go run gen.go

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	tmpls := make(map[string][]string)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to access path %q: %v\n", path, err)
			return err
		}
		details := strings.SplitN(path, ".", 2)
		if len(details) != 2 || details[1] != "sql" {
			//fmt.Fprintf(os.Stderr, "Skipping non-template: %+v \n", info.Name())
			return nil
		}
		pathDetails := strings.SplitN(details[0], "/", 2)
		driver := pathDetails[0]
		//name := pathDetails[1]

		fmt.Fprintf(os.Stderr, "Processing file: %s\n", path)
		input, err := os.Open(path)
		if err != nil {
			return err
		}
		var out bytes.Buffer

		encoder := base64.NewEncoder(base64.StdEncoding, &out)
		compressor := zlib.NewWriter(encoder)
		_, err = io.Copy(compressor, input)

		input.Close()
		compressor.Close()
		encoder.Close()

		if err != nil {
			return err
		}
		tmpls[driver] = append(tmpls[driver], out.String())
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load templates: %s\n", err)
	}

	// Now generate our encoded template file
	metaTmpl, err := template.New("meta").Parse(tmplData)
	if err != nil {
		fmt.Printf("Failed to parse template: %s", err)
	} else {
		metaTmpl.Execute(os.Stdout, tmpls)
	}
}

var tmplData = `package main

// Automatically generated, don't bother editing.

import (
	"bytes"
	"strings"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)

// migrations holds a set of base64 encoded SQL scripts.
var migrations = map[string][]string{
{{- range $driver, $schema := . }}
	"{{$driver}}": []string{
{{- range $b := $schema }}
		"{{ $b }}",
{{ end }}
	},
{{ end }}
}

func decodeMigrations(driver string) ([]string, error) {
	data, ok := migrations[driver]
	if !ok {
		return nil, fmt.Errorf("no migrations for driver %q", driver)
	}

	out := make([]string, len(data))

	for i, b := range data {
		decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b))
		decompressor, err := zlib.NewReader(decoder)
		if err != nil {
			return nil, fmt.Errorf("unable to decode migration: %w", err)
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, decompressor)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress migration: %w", err)
		}
		out[i] = buf.String()
	}

	return out, nil
}
`
