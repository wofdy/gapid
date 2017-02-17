// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The embed command is used to embed text files into Go executables as strings.
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/google/gapid/core/app"
	"github.com/google/gapid/core/log"
	"golang.org/x/tools/imports"
)

var (
	pkg    string
	output string
	root   string
)

const header = `
////////////////////////////////////////////////////////////////////////////////
// Do not modify!
// Generated by embed
////////////////////////////////////////////////////////////////////////////////

package %s

`

func main() {
	flag.StringVar(&pkg, "package", "", "the package to use, defaults to the dir name")
	flag.StringVar(&output, "out", "embed.go", "the file to generate")
	flag.StringVar(&root, "root", "", "The root path for embedding")
	app.ShortHelp = "embed: A tool to embed files into go source."
	app.Run(run)
}

type embed struct {
	path     string
	filename string
	name     string
	contents []byte
}

var (
	nameReplacer = strings.NewReplacer(
		"-", "_",
		".", "_",
	)
	contentReplacer = strings.NewReplacer(
		"`", "` + \"`\" + `",
	)
)

func run(ctx log.Context) error {
	args := flag.Args()
	entries := []*embed{}
	if len(args) == 0 {
		pwd, err := filepath.Abs(".")
		if err != nil {
			return err
		}
		files, err := ioutil.ReadDir(pwd)
		if err != nil {
			return err
		}
		if root == "" {
			root = pwd
		}
		for _, info := range files {
			if info.IsDir() {
				continue
			}
			extension := filepath.Ext(info.Name())
			if extension == ".go" {
				continue
			}
			path := filepath.Join(pwd, info.Name())
			// if Rel fails, we just fallback to basename in the loop below
			rel, _ := filepath.Rel(root, info.Name())
			entries = append(entries, &embed{filename: rel, path: path})
		}
	} else {
		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return err
			}
			rel := ""
			if root != "" {
				// if Rel fails, we just fallback to basename in the loop below
				rel, _ = filepath.Rel(root, path)
			} else {
				rel = filepath.Base(path)
			}
			entries = append(entries, &embed{filename: rel, path: path})
		}
	}
	var err error
	for _, entry := range entries {
		filename := filepath.Base(entry.path)
		if entry.filename == "" {
			// filename will be empty if path was not relative to root, so just use basename
			entry.filename = filename
		}
		entry.name = nameReplacer.Replace(filename)
		entry.contents, err = ioutil.ReadFile(entry.path)
		if err != nil {
			return err
		}
	}
	// write the header
	out, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	if pkg == "" {
		pkg = filepath.Base(filepath.Dir(out))
	}
	b := &bytes.Buffer{}
	fmt.Fprintf(b, header, pkg)
	// write the map
	fmt.Fprint(b, "var embedded = map[string]string{\n")
	for _, entry := range entries {
		fmt.Fprintf(b, "%s_file: %s,\n", entry.name, entry.name)
	}
	fmt.Fprint(b, "}\n")
	fmt.Fprint(b, "var embedded_utf8 = map[string]bool{\n")
	for _, entry := range entries {
		fmt.Fprintf(b, "%s_file: %s_utf8,\n", entry.name, entry.name)
	}
	fmt.Fprint(b, "}\n")
	// write the data lumps
	for _, entry := range entries {
		validUTF8 := utf8.Valid(entry.contents)
		encoded := ""
		if validUTF8 {
			encoded = contentReplacer.Replace(string(entry.contents))
		} else {
			encoded = base64.StdEncoding.EncodeToString(entry.contents)
		}
		fmt.Fprintf(b, "const %s_file = `%s`\n", entry.name, entry.filename)
		fmt.Fprintf(b, "const %s_utf8 = %v\n", entry.name, validUTF8)
		fmt.Fprintf(b, "const %s = `%s`\n", entry.name, encoded)
		ctx.Printf("Embed %s from %s\n", entry.name, entry.path)
	}
	// reformat the output
	result, err := imports.Process("", b.Bytes(), nil)
	if err != nil {
		result = b.Bytes()
	}
	if err := ioutil.WriteFile(out, result, 0666); err != nil {
		return err
	}
	return nil
}
