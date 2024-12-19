// Package test provides common functions for testing.
package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

// UpdateOut is a flag that can be set to update example files
var UpdateOut = flag.Bool("update", false, "Update the example files in test/data and test/data/out")

// LoadEnvelope loads a test file from the test/data folder as a GOBL envelope
// and will rebuild it if necessary to ensure any changes are accounted for.
func LoadEnvelope(file string) *gobl.Envelope {
	path := Path("test", "data", file)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		panic(err)
	}

	out, err := gobl.Parse(buf.Bytes())
	if err != nil {
		panic(err)
	}

	var env *gobl.Envelope
	switch doc := out.(type) {
	case *gobl.Envelope:
		env = doc
	default:
		env = gobl.NewEnvelope()
		if err := env.Insert(doc); err != nil {
			panic(err)
		}
	}

	if err := env.Calculate(); err != nil {
		panic(err)
	}

	if err := env.Validate(); err != nil {
		panic(err)
	}

	if *UpdateOut {
		data, err := json.MarshalIndent(env, "", "\t")
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}

	return env
}

// LoadInvoice grabs the gobl envelope and attempts to extract the invoice payload
func LoadInvoice(name string) *bill.Invoice {
	env := LoadEnvelope(name)

	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		panic("envelope does not contain an invoice")
	}

	return inv
}

// Path joins the provided elements to the project root
func Path(elements ...string) string {
	np := []string{RootPath()}
	np = append(np, elements...)
	return path.Join(np...)
}

// RootPath returns the project root, regardless of current working directory.
func RootPath() string {
	cwd, _ := os.Getwd()

	for !isRootFolder(cwd) {
		cwd = removeLastEntry(cwd)
	}

	return cwd
}

func isRootFolder(dir string) bool {
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}

	return false
}

func removeLastEntry(dir string) string {
	lastEntry := "/" + filepath.Base(dir)
	i := strings.LastIndex(dir, lastEntry)
	return dir[:i]
}
