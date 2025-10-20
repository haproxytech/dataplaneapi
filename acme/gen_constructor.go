//go:build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

var (
	modListFile  = flag.String("i", "", "input file with the list of DNS provider modules")
	templateFile = flag.String("t", "", "template file")
	outputFile   = flag.String("o", "", "output filename")
)

func main() {
	var err error
	flag.Parse()

	in := os.Stdin
	if *modListFile != "" {
		in, err = os.Open(*modListFile)
		chk(err)
	}

	out := os.Stdout
	if *outputFile != "" {
		out, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		chk(err)
	}

	funcs := template.FuncMap{
		"basename": modname,
	}

	tmpl, err := template.New(path.Base(*templateFile)).Funcs(funcs).ParseFiles(*templateFile)
	chk(err)

	// Read the list of modules, 1 per line.
	modules := make([]string, 0, 32)
	lines, err := io.ReadAll(in)
	chk(err)
	bytes.SplitSeq(lines, []byte{'\n'})(func(line []byte) bool {
		line = bytes.Trim(line, " \"\n\t\r")
		if len(line) > 0 && line[0] != '/' && line[0] != '#' {
			modules = append(modules, string(line))
		}
		return true
	})

	chk(tmpl.Execute(out, modules))
}

func chk(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Like basename but for Go repository names.
func modname(s string) string {
	slash := strings.LastIndexByte(s, '/')
	// ignore "/v2"
	match, err := regexp.MatchString("^/v[1-9]+$", s[slash:])
	chk(err)
	if match {
		s = s[:slash]
		slash = strings.LastIndexByte(s, '/')
	}
	s = s[slash+1:]
	if dash := strings.IndexByte(s, '-'); dash != -1 {
		s = s[dash+1:]
	}
	if end := strings.IndexByte(s, '.'); end != -1 {
		s = s[:end]
	}
	return s
}
