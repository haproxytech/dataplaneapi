package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/haproxytech/dataplaneapi/cmd/generate/parents"
)

func main() {
	specDir := flag.String("spec-dir", "specification", "Path to the specification directory")
	embeddedOut := flag.String("embedded-out", "handlers/dataplane/specification/dataplane_spec.gen.go", "Output path for the embedded Go spec file")
	oapiCodegen := flag.String("oapi-codegen", filepath.Join("bin", "oapi-codegen"), "Path to the pinned oapi-codegen binary")
	flag.Parse()

	// Only the version-pinned binary that `make specification` installs into bin/
	// may be used: an unpinned oapi-codegen from PATH produces different .gen.go
	// output across versions, silently rewriting every generated file.
	if _, err := os.Stat(*oapiCodegen); err != nil {
		fmt.Fprintf(os.Stderr, "error: pinned oapi-codegen not found at %s — generate via 'make specification'\n", *oapiCodegen)
		os.Exit(1)
	}

	pathsDir := filepath.Join(*specDir, "paths")

	var failed []string

	err := filepath.Walk(pathsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != "handlers.conf.yaml" {
			return nil
		}

		dir := filepath.Dir(path)
		folderName := filepath.Base(dir)
		oapiSpec := filepath.Join(dir, "oapi.yaml")

		if _, statErr := os.Stat(oapiSpec); os.IsNotExist(statErr) {
			// Skipping here would leave the resource's stale .gen.go in place.
			fmt.Fprintf(os.Stderr, "error: no oapi.yaml next to %s\n", path)
			failed = append(failed, path)
			return nil
		}

		fmt.Printf("generating %s\n", path)

		specPath, cleanup, prepErr := prepareSpec(oapiSpec, folderName)
		if prepErr != nil {
			fmt.Fprintf(os.Stderr, "error: preparing spec %s: %v\n", oapiSpec, prepErr)
			failed = append(failed, path)
			return nil
		}
		defer cleanup()

		cmd := exec.Command(*oapiCodegen, "--config="+path, specPath) //nolint:gosec // dev-time code generator; paths come from the trusted spec tree
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if runErr := cmd.Run(); runErr != nil {
			fmt.Fprintf(os.Stderr, "error: oapi-codegen failed for %s: %v\n", path, runErr)
			failed = append(failed, path)
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: walking %s: %v\n", pathsDir, err)
		os.Exit(1)
	}

	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "%d spec(s) failed to generate\n", len(failed))
		os.Exit(1)
	}

	if err := buildCombinedSpec(*specDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: building combined spec: %v\n", err)
		os.Exit(1)
	}

	if err := buildEmbeddedGoSpec(*specDir, *embeddedOut); err != nil {
		fmt.Fprintf(os.Stderr, "error: building embedded Go spec: %v\n", err)
		os.Exit(1)
	}
}

// prepareSpec returns the path to the spec file to pass to oapi-codegen and a cleanup function.
// For child resources with parents, it template-expands oapi.yaml into a temp file first.
func prepareSpec(oapiSpec, folderName string) (specPath string, cleanup func(), err error) {
	resourceParents := parents.Parents(folderName)
	if len(resourceParents) == 0 {
		return oapiSpec, func() {}, nil
	}

	tmplData, err := os.ReadFile(oapiSpec)
	if err != nil {
		return "", nil, err
	}

	tmpl, err := template.New("oapi").Funcs(template.FuncMap{
		"parents": func() []parents.Parent { return resourceParents },
	}).Parse(string(tmplData))
	if err != nil {
		return "", nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, nil); err != nil {
		return "", nil, fmt.Errorf("executing template: %w", err)
	}

	// Write temp file alongside oapi.yaml so relative $refs resolve correctly.
	tmpFile, err := os.CreateTemp(filepath.Dir(oapiSpec), "oapi-expanded-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp file: %w", err)
	}
	if _, err = tmpFile.Write(buf.Bytes()); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("writing temp file: %w", err)
	}
	tmpFile.Close()

	return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }, nil
}
