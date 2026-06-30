package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func buildCombinedSpec(specDir string) error {
	pathsDir := filepath.Join(specDir, "paths")
	modelsDir := filepath.Join(specDir, "models")
	buildDir := filepath.Join(specDir, "build")

	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("creating build dir: %w", err)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	ctx := context.Background()

	combined := &openapi3.T{
		OpenAPI: "3.0.1",
		Info:    &openapi3.Info{Title: "HAProxy Data Plane API", Version: "3.4"},
		Servers: openapi3.Servers{{URL: "/v3"}},
		Paths:   openapi3.NewPathsWithCapacity(256),
		Components: &openapi3.Components{
			Schemas:   make(openapi3.Schemas),
			Responses: make(openapi3.ResponseBodies),
			SecuritySchemes: openapi3.SecuritySchemes{
				"basic_auth": &openapi3.SecuritySchemeRef{
					Value: &openapi3.SecurityScheme{Type: "http", Scheme: "basic"},
				},
			},
		},
		// Global security requirement, mirroring the original specification.
		// This is declarative in the published document; authentication is enforced
		// by a separate middleware layer, not the OpenAPI request validator.
		Security: openapi3.SecurityRequirements{
			openapi3.SecurityRequirement{"basic_auth": []string{}},
		},
	}

	cnSchemas, cnErr := loadClientNativeSchemas(specDir)
	if cnErr != nil {
		// Proceeding without these schemas would silently produce a structurally
		// different spec (x-go-type stubs left unreplaced).
		return fmt.Errorf("loading client-native schemas: %w", cnErr)
	}

	// Include everything declared in models/ (ensures completeness even for unreferenced items)
	if err := filepath.Walk(modelsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".yaml" {
			return err
		}
		return mergeModels(loader, path, combined)
	}); err != nil {
		return fmt.Errorf("processing models: %w", err)
	}

	// Merge paths from every oapi.yaml
	if err := filepath.Walk(pathsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "oapi.yaml" {
			return err
		}
		specPath, cleanup, prepErr := prepareSpec(path, filepath.Base(filepath.Dir(path)))
		if prepErr != nil {
			return fmt.Errorf("preparing %s: %w", path, prepErr)
		}
		defer cleanup()
		return mergeAPIPaths(ctx, loader, specPath, combined)
	}); err != nil {
		return fmt.Errorf("processing paths: %w", err)
	}

	// Replace x-go-type stubs with their full definitions from client-native.
	if cnSchemas != nil {
		replaceStubSchemas(combined.Components.Schemas, cnSchemas)
		// Seed any remaining cn schemas not already present so that internal $refs
		// within replaced schemas (e.g. "#/components/schemas/balance_base") resolve.
		for k, v := range cnSchemas {
			if _, exists := combined.Components.Schemas[k]; !exists {
				combined.Components.Schemas[k] = v
			}
		}
	}

	// Marshal to an intermediate yaml.Node so we can control key ordering.
	var raw bytes.Buffer
	rawEnc := yaml.NewEncoder(&raw)
	rawEnc.SetIndent(2)
	if err := rawEnc.Encode(combined); err != nil {
		return fmt.Errorf("marshaling combined spec: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(raw.Bytes(), &doc); err != nil {
		return fmt.Errorf("parsing spec for reorder: %w", err)
	}
	if len(doc.Content) > 0 {
		// Move components and paths to the end of the top-level mapping.
		sortMappingKeysLast(doc.Content[0], "components", "paths")
	}

	outPath := filepath.Join(buildDir, "dataplane_spec.yaml")
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return fmt.Errorf("encoding reordered spec: %w", err)
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil { //nolint:gosec // generated spec is non-sensitive and meant to be world-readable
		return fmt.Errorf("writing %s: %w", outPath, err)
	}
	fmt.Printf("wrote combined spec: %s\n", outPath)
	return nil
}

func mergeModels(loader *openapi3.Loader, path string, combined *openapi3.T) error {
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("loading %s: %w", path, err)
	}
	if doc.Components == nil {
		return nil
	}
	if err := mergeComponents(combined.Components.Schemas, doc.Components.Schemas, "schema", path); err != nil {
		return err
	}
	return mergeComponents(combined.Components.Responses, doc.Components.Responses, "response", path)
}

func mergeAPIPaths(ctx context.Context, loader *openapi3.Loader, specPath string, combined *openapi3.T) error {
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return fmt.Errorf("loading %s: %w", specPath, err)
	}

	if doc.Components == nil {
		doc.Components = &openapi3.Components{
			Schemas:   make(openapi3.Schemas),
			Responses: make(openapi3.ResponseBodies),
		}
	}

	// Rewrite external $refs (e.g. models.yaml#/components/responses/DefaultError)
	// into internal #/components/... refs and pull them into doc.Components.
	doc.InternalizeRefs(ctx, componentBaseName)

	if err := mergeComponents(combined.Components.Schemas, doc.Components.Schemas, "schema", specPath); err != nil {
		return err
	}
	if err := mergeComponents(combined.Components.Responses, doc.Components.Responses, "response", specPath); err != nil {
		return err
	}

	// Merge paths
	for p, item := range doc.Paths.Map() {
		if combined.Paths.Value(p) != nil {
			return fmt.Errorf("duplicate path %s (from %s)", p, specPath)
		}
		combined.Paths.Set(p, item)
	}

	return nil
}

// mergeComponents copies src components into dst. The same component arriving from several
// files is fine as long as every copy is identical (InternalizeRefs pulls shared models like
// DefaultError into each package's doc); two different definitions flattened to the same
// name would silently shadow each other, so that is a hard error.
func mergeComponents[V any](dst, src map[string]V, kind, source string) error {
	for name, ref := range src {
		existing, ok := dst[name]
		if !ok {
			dst[name] = ref
			continue
		}
		a, err := json.Marshal(existing)
		if err != nil {
			return fmt.Errorf("marshaling %s component %q: %w", kind, name, err)
		}
		b, err := json.Marshal(ref)
		if err != nil {
			return fmt.Errorf("marshaling %s component %q (from %s): %w", kind, name, source, err)
		}
		if !bytes.Equal(a, b) {
			return fmt.Errorf("%s component %q from %s conflicts with an existing definition of the same name", kind, name, source)
		}
	}
	return nil
}

// componentBaseName extracts just the final component name from a ref string.
// "models.yaml#/components/responses/DefaultError" → "DefaultError"
func componentBaseName(_ *openapi3.T, ref openapi3.ComponentRef) string {
	s := ref.RefString()
	if i := strings.LastIndex(s, "/"); i >= 0 {
		return s[i+1:]
	}
	return s
}

// sortMappingKeysLast moves the named keys to the end of a YAML mapping node,
// preserving the relative order of both the moved and unmoved keys.
func sortMappingKeysLast(node *yaml.Node, keys ...string) {
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}
	last := make(map[string]bool, len(keys))
	for _, k := range keys {
		last[k] = true
	}

	type pair struct{ k, v *yaml.Node }
	var head, tail []pair
	for i := 0; i+1 < len(node.Content); i += 2 {
		p := pair{node.Content[i], node.Content[i+1]}
		if last[p.k.Value] {
			tail = append(tail, p)
		} else {
			head = append(head, p)
		}
	}

	node.Content = node.Content[:0]
	for _, p := range append(head, tail...) {
		node.Content = append(node.Content, p.k, p.v)
	}
}
