package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"gopkg.in/yaml.v3"
)

// loadClientNativeSchemas fetches the HAProxy client-native Swagger 2.0 spec for the pinned
// dependency version, converts it to OpenAPI 3, and returns all component schemas.
func loadClientNativeSchemas(specDir string) (openapi3.Schemas, error) {
	goModPath := filepath.Join(specDir, "..", "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("reading go.mod: %w", err)
	}

	version, modulePath, err := findClientNativeModule(goModData)
	if err != nil {
		return nil, err
	}

	data, err := fetchClientNativeSpec(modulePath, version, filepath.Dir(goModPath))
	if err != nil {
		return nil, err
	}

	v3doc, err := convertSwaggerToV3(data)
	if err != nil {
		return nil, fmt.Errorf("converting spec: %w", err)
	}

	if v3doc.Components == nil {
		return openapi3.Schemas{}, nil
	}
	return v3doc.Components.Schemas, nil
}

// findClientNativeModule parses go.mod and returns the version string and module path of the
// direct github.com/haproxytech/client-native/vN dependency, honoring replace directives.
// A local filesystem replace returns the replacement directory as modulePath with an empty
// version; the caller reads the spec from disk instead of fetching it.
func findClientNativeModule(goModData []byte) (version, modulePath string, err error) {
	f, err := modfile.Parse("go.mod", goModData, nil)
	if err != nil {
		return "", "", fmt.Errorf("parsing go.mod: %w", err)
	}

	const cnPrefix = "github.com/haproxytech/client-native/v"
	for _, req := range f.Require {
		if req.Indirect || !strings.HasPrefix(req.Mod.Path, cnPrefix) {
			continue
		}
		version, modulePath = req.Mod.Version, req.Mod.Path
		for _, rep := range f.Replace {
			if rep.Old.Path != modulePath || (rep.Old.Version != "" && rep.Old.Version != version) {
				continue
			}
			if rep.New.Version == "" {
				// Filesystem replacement: New.Path is a directory.
				return "", rep.New.Path, nil
			}
			version, modulePath = rep.New.Version, rep.New.Path
		}
		return version, modulePath, nil
	}
	return "", "", errors.New("direct client-native dependency not found in go.mod")
}

// specSubPath is the location of the HAProxy OpenAPI 2 spec inside the client-native repository.
const specSubPath = "specification/build/haproxy_spec.yaml"

// fetchClientNativeSpec returns the haproxy_spec.yaml contents for the given client-native
// module path and version. An empty version means a local filesystem replace directive:
// modulePath is then a directory (relative to goModDir unless absolute) and the spec is read
// from disk. Otherwise only the spec file is downloaded directly from the source repository's
// raw file endpoint. The whole module is never pulled: a Go module cache may hold the
// dependency as an un-extracted zip (e.g. in clean CI jobs that build no package importing
// client-native), so reading the file from disk is unreliable.
func fetchClientNativeSpec(modulePath, version, goModDir string) ([]byte, error) {
	if version == "" {
		dir := modulePath
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(goModDir, dir)
		}
		localPath := filepath.Join(dir, filepath.FromSlash(specSubPath))
		fmt.Printf("reading client-native spec from local replace: %s\n", localPath)
		return os.ReadFile(localPath) //nolint:gosec // dev-time generator; path comes from the developer's own go.mod replace directive
	}

	specURL, err := rawSpecURL(modulePath, version)
	if err != nil {
		return nil, err
	}

	fmt.Printf("fetching client-native spec: %s\n", specURL)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(specURL) //nolint:noctx // dev-time generator; timeout is set on the client
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", specURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching %s: unexpected status %s", specURL, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", specURL, err)
	}
	return data, nil
}

// rawSpecURL builds the raw-file URL for the spec from a Go module path
// (e.g. github.com/haproxytech/client-native/v6) and a module version. Tagged versions are used
// as the git ref directly; pseudo-versions (vX.Y.Z-<timestamp>-<commit>) resolve to their commit
// hash. GitHub and GitLab raw-content URL layouts are both supported.
func rawSpecURL(modulePath, version string) (string, error) {
	parts := strings.Split(modulePath, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected module path %q", modulePath)
	}
	host, owner, repo := parts[0], parts[1], parts[2]

	ref := gitRefFromVersion(version)

	if strings.Contains(host, "gitlab") {
		return fmt.Sprintf("https://%s/%s/%s/-/raw/%s/%s", host, owner, repo, ref, specSubPath), nil
	}
	// Default to GitHub's raw-content host.
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, specSubPath), nil
}

// gitRefFromVersion returns the git ref to fetch for a module version. For a Go pseudo-version
// (e.g. v6.3.1-0.20260520132134-5762461b033f) the trailing commit hash is the ref;
// for a normal tag (e.g. v6.3.8) the version string itself is the ref.
func gitRefFromVersion(version string) string {
	if module.IsPseudoVersion(version) {
		if i := strings.LastIndex(version, "-"); i >= 0 {
			return version[i+1:]
		}
	}
	return version
}

// convertSwaggerToV3 parses a Swagger 2.0 YAML document and converts it to an OpenAPI 3 doc.
// The YAML→map→JSON→openapi2.T pipeline is used for maximum compatibility.
func convertSwaggerToV3(data []byte) (*openapi3.T, error) {
	// Decode into interface{} first; yaml.v3 can produce map[interface{}]interface{} for
	// nested maps, which json.Marshal cannot handle — normalizeForJSON converts them to
	// map[string]interface{} recursively.
	var raw any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	jsonData, err := json.Marshal(normalizeForJSON(raw))
	if err != nil {
		return nil, fmt.Errorf("re-encoding as JSON: %w", err)
	}
	var v2doc openapi2.T
	if err := json.Unmarshal(jsonData, &v2doc); err != nil {
		return nil, fmt.Errorf("parsing OpenAPI 2 doc: %w", err)
	}
	return openapi2conv.ToV3(&v2doc)
}

// normalizeForJSON recursively converts map[interface{}]interface{} (produced by yaml.v3 for
// some nested structures) to map[string]interface{} so json.Marshal can serialize it.
func normalizeForJSON(v any) any {
	switch val := v.(type) {
	case map[any]any:
		out := make(map[string]any, len(val))
		for k, elem := range val {
			out[fmt.Sprintf("%v", k)] = normalizeForJSON(elem)
		}
		return out
	case map[string]any:
		for k, elem := range val {
			val[k] = normalizeForJSON(elem)
		}
		return val
	case []any:
		for i, elem := range val {
			val[i] = normalizeForJSON(elem)
		}
		return val
	default:
		return v
	}
}

// replaceStubSchemas replaces each schema in schemas that is a stub (only x-go-type, no real
// definition) with the full definition from cnSchemas, preserving all x-go-* extensions.
// It tries the schema's own name first (cn uses lowercase names), then the CamelCase Go type name.
func replaceStubSchemas(schemas openapi3.Schemas, cnSchemas openapi3.Schemas) {
	for name, ref := range schemas {
		if ref == nil || ref.Value == nil {
			continue
		}
		goType := extractGoType(ref.Value)
		if goType == "" {
			continue
		}
		// cn spec uses lowercase definition names (e.g. "backend"), but x-go-type is CamelCase
		// ("models.Backend"). Try the schema's own name first, then the CamelCase type name.
		cnRef, ok := cnSchemas[name]
		if !ok || cnRef == nil || cnRef.Value == nil {
			cnRef, ok = cnSchemas[goType]
		}
		if !ok || cnRef == nil || cnRef.Value == nil {
			continue
		}
		// Shallow-copy the cn schema value and overlay our x-go-* extensions so codegen still works.
		newVal := *cnRef.Value
		merged := make(map[string]any, len(newVal.Extensions)+len(ref.Value.Extensions))
		maps.Copy(merged, newVal.Extensions)
		// our x-go-* extensions take priority
		maps.Copy(merged, ref.Value.Extensions)
		newVal.Extensions = merged
		schemas[name] = openapi3.NewSchemaRef("", &newVal)
	}
}

// extractGoType returns the bare type name from an x-go-type extension
// (e.g. "models.Backend" → "Backend"), or "" if the schema has no such extension.
func extractGoType(schema *openapi3.Schema) string {
	if len(schema.Extensions) == 0 {
		return ""
	}
	raw, ok := schema.Extensions["x-go-type"]
	if !ok {
		return ""
	}
	var s string
	switch v := raw.(type) {
	case string:
		s = v
	case json.RawMessage:
		if err := json.Unmarshal(v, &s); err != nil {
			return ""
		}
	default:
		return ""
	}
	if i := strings.LastIndex(s, "."); i >= 0 {
		return s[i+1:]
	}
	return s
}
