// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

var serverGroups = map[string]string{
	"EnabledListeners":  "dataplaneapi",
	"CleanupTimeout":    "dataplaneapi",
	"GracefulTimeout":   "dataplaneapi",
	"MaxHeaderSize":     "dataplaneapi",
	"SocketPath":        "dataplaneapi",
	"DebugSocketPath":   "dataplaneapi",
	"Host":              "dataplaneapi",
	"Port":              "dataplaneapi",
	"ListenLimit":       "dataplaneapi",
	"KeepAlive":         "dataplaneapi",
	"ReadTimeout":       "dataplaneapi",
	"WriteTimeout":      "dataplaneapi",
	"TLSHost":           "tls",
	"TLSPort":           "tls",
	"TLSCertificate":    "tls",
	"TLSCertificateKey": "tls",
	"TLSCACertificate":  "tls",
	"TLSListenLimit":    "tls",
	"TLSKeepAlive":      "tls",
	"TLSReadTimeout":    "tls",
	"TLSWriteTimeout":   "tls",
}

var groupParents = map[string]string{
	"advertised":  "dataplaneapi",
	"userlist":    "dataplaneapi",
	"resources":   "dataplaneapi",
	"transaction": "dataplaneapi",
	"tls":         "dataplaneapi",
	"user":        "dataplaneapi",
	"reload":      "haproxy",
	"syslog":      "log",
	"log_targets": "log",
}

var itemDefaults = map[string]interface{}{
	"port":              80,
	"listen_limit":      1024,
	"tls_host":          "null",
	"tls_port":          6443,
	"tls_certificate":   "null",
	"tls_key":           "null",
	"tls_ca":            "null",
	"tls_listen_limit":  10,
	"tls_keep_alive":    "1m",
	"tls_read_timeout":  "10s",
	"tls_write_timeout": "10s",
	"userlist_file":     "null",
	"backups_dir":       "/tmp/backups",
	"scheme":            "http",
}

type Attribute struct {
	Default    string
	Group      string
	Type       string
	Long       string
	FileName   string
	Short      string
	AttName    string
	ENV        string
	Name       string
	StructName string
	SpecName   string
	Example    string
	Deprecated bool
	Save       bool
}

type ParseGroup struct {
	OriginalGroup string
	Name          string
	AttName       string
	Parent        string
	Elements      []string
	Attributes    []Attribute
	MaxSize       int
	MaxTypeSize   int
	HasACLKey     bool
	IsList        bool
}

type ParseData struct {
	Groups []ParseGroup
}

func readServerData(filePath string, pd *ParseData, structName string, attName string, groupName string, isList bool) {
	typeStruct := fmt.Sprintf("type %s struct {", structName)
	dat, err := os.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	lines := strings.Split(string(dat), "\n")
	inside := false
	for _, line := range lines {
		if strings.HasPrefix(line, typeStruct) {
			inside = true
			continue
		}
		if strings.HasPrefix(line, "}") {
			inside = false
			continue
		}
		if inside && strings.Contains(line, "`") {
			res, err := processLine(line)
			if err != nil {
				continue
			}
			if structName == "Server" {
				group, ok := serverGroups[res.Name]
				if !ok {
					log.Panicf("group not defined for %s in serverGroups", res.Name)
				}
				res.Group = group
			}
			res.AttName = attName
			res.StructName = structName
			found := false
			if groupName == "" {
				groupName = res.Group
			}
			for i, g := range pd.Groups {
				if g.Name == res.Group {
					found = true
					g.Attributes = append(g.Attributes, res)
					if g.MaxSize < len(res.Name) {
						g.MaxSize = len(res.Name)
					}
					if g.MaxTypeSize < len(res.Type) {
						g.MaxSize = len(res.Type)
					}
					pd.Groups[i] = g
					break
				}
			}
			if !found {
				pd.Groups = append(pd.Groups, ParseGroup{
					OriginalGroup: structName,
					Name:          res.Group,
					AttName:       attName,
					Parent:        groupParents[res.Group],
					MaxSize:       len(res.Name),
					MaxTypeSize:   len(res.Type),
					Attributes:    []Attribute{res},
					IsList:        isList,
				})
			}
		}
	}
}

func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}
	var res strings.Builder
	parts := strings.FieldsFunc(str, split)
	for _, part := range parts {
		res.WriteString(capitalizeChunk(part))
	}
	return res.String()
}

func split(r rune) bool {
	return r == '_' || r == '-'
}

func capitalizeChunk(str string) string {
	if len(str) == 0 {
		return ""
	}
	result := []rune(str)
	result[0] = unicode.ToUpper(result[0])
	return string(result)
}

func stripAtomic(str string) string {
	if len(str) == 0 {
		return ""
	}
	if strings.HasPrefix(str, "Atomic") {
		return strings.ToLower(strings.TrimPrefix(str, "Atomic"))
	}
	return str
}

func isListItem(att Attribute) string {
	return "  "
}

func getExample(att Attribute) string {
	if att.Example != "" {
		if att.Type == "int" {
			i, _ := strconv.ParseInt(att.Example, 10, 64)
			return strconv.FormatInt(i, 10)
		}
		return att.Example
	}
	if att.Default != "" {
		switch att.Type {
		case "int", "int64":
			return att.Default
		default:
			return fmt.Sprintf(`"%s"`, att.Default)
		}
	}
	if att.Type == "bool" {
		return `false`
	}
	if strings.HasPrefix(att.Type, "[]*models.") {
		return `[]`
	}
	if v, ok := itemDefaults[att.FileName]; ok {
		switch t := v.(type) {
		case int:
			return strconv.Itoa(v.(int))
		default:
			return t.(string)
		}
	}
	return "null"
}

func getQuotedExample(att Attribute) string {
	if att.Type == "int" || att.Type == "int64" {
		switch {
		case len(att.Example) > 0:
			return att.Example
		case len(att.Default) > 0:
			return att.Default
		}
	}
	if att.Type == "bool" || att.Type == "AtomicBool" {
		if len(att.Example) == 0 {
			return `false`
		}
		return att.Example
	}
	if strings.HasPrefix(att.Type, "[]*models.") {
		return `[]`
	}
	if att.Example != "" {
		return fmt.Sprintf(`"%s"`, att.Example)
	}
	if att.Default != "" {
		return fmt.Sprintf(`"%s"`, att.Default)
	}
	if v, ok := itemDefaults[att.FileName]; ok {
		switch v.(type) {
		case int:
			return fmt.Sprintf("%d", v)
		default:
			return fmt.Sprintf(`"%s"`, v)
		}
	}
	return `"null"`
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasSuffix(dir, "generate") {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	pd := &ParseData{
		Groups: []ParseGroup{},
	}

	pd.Groups = append(pd.Groups, ParseGroup{
		Name:        "",
		AttName:     "-",
		Parent:      "",
		MaxSize:     4,
		MaxTypeSize: 12,
		Attributes:  []Attribute{},
	})
	// ######################################## server.go
	filePath := path.Join(dir, "server.go")
	readServerData(filePath, pd, "Server", "", "", false)
	// ######################################## configuration.go
	filePath = path.Join(dir, "configuration", "configuration.go")
	readServerData(filePath, pd, "Configuration", "-", "", false)
	readServerData(filePath, pd, "User", "Users", "", true)
	readServerData(filePath, pd, "HAProxyConfiguration", "HAProxy", "", false)
	readServerData(filePath, pd, "APIConfiguration", "APIOptions", "", false)
	readServerData(filePath, pd, "ServiceDiscovery", "ServiceDiscovery", "", false)
	readServerData(filePath, pd, "ClusterConfiguration", "Cluster", "cluster", false)
	readServerData(filePath, pd, "SyslogOptions", "Syslog", "", false)
	readServerData(filePath, pd, "LoggingOptions", "Logging", "", false)
	// readServerData(filePath, pd, "User", "User", "")
	// ########################################

	// prepare template function
	funcMap := template.FuncMap{
		"Capitalize":    capitalize,
		"StripAtomic":   stripAtomic,
		"Example":       getExample,
		"QuotedExample": getQuotedExample,
		"IsListItem":    isListItem,
	}
	// create configuration_generated
	templatePath := path.Join(dir, "generate", "configuration.tmpl")
	tmpl, err := template.New("configuration.tmpl").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		log.Panic(err)
	}
	tmpl = tmpl.Funcs(funcMap)
	filePath = path.Join(dir, "configuration", "configuration_generated.go")
	f, err := os.Create(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, pd)
	if err != nil {
		log.Panic(err)
	}
	fmtFile(filePath)
	// create dataplaneapi_generated
	templatePath = path.Join(dir, "generate", "dataplaneapi.tmpl")
	tmpl, err = template.New("dataplaneapi.tmpl").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		log.Panic(err)
	}
	tmpl = tmpl.Funcs(funcMap)
	filePath = path.Join(dir, "dataplaneapi_generated.go")
	f, err = os.Create(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, pd)
	if err != nil {
		log.Panic(err)
	}
	fmtFile(filePath)
	// create configuration example yaml
	templatePath = path.Join(dir, "generate", "configuration-example-yaml.tmpl")
	tmpl, err = template.New("configuration-example-yaml.tmpl").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		log.Panic(err)
	}
	tmpl = tmpl.Funcs(funcMap)
	filePath = path.Join(dir, "configuration/examples/example-full.yaml")
	f, err = os.Create(filePath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	err = tmpl.Execute(f, pd)
	if err != nil {
		log.Panic(err)
	}
}

func splitBy(c rune) bool {
	return c == ' ' || c == '\t'
}

func processLine(line string) (Attribute, error) {
	parts := strings.FieldsFunc(line, splitBy)
	var longName string
	var envName string
	var shortName string
	var defaultName string
	var group string
	var specName string
	var save bool
	var deprecated bool
	var example string
	for _, part := range parts[2:] {
		if strings.Contains(part, "long:") {
			p := strings.Split(part, `"`)
			longName = p[1]
		}
		if strings.Contains(part, "env:") {
			p := strings.Split(part, `"`)
			envName = p[1]
		}
		if strings.Contains(part, "short:") {
			p := strings.Split(part, `"`)
			shortName = p[1]
		}
		if strings.Contains(part, "default:") {
			p := strings.Split(part, `"`)
			defaultName = p[1]
		}
		if strings.Contains(part, "example:") {
			p := strings.Split(part, `"`)
			example = p[1]
		}
		if strings.Contains(part, "group:") {
			p := strings.Split(part, `"`)
			group = p[1]
		}
		if strings.Contains(part, `save:"true"`) {
			save = true
		}
		if strings.Contains(part, `deprecated:"true"`) {
			deprecated = true
		}
		if strings.Contains(part, `yaml:"-"`) {
			return Attribute{}, errors.New("ignore this attribute")
		} else if strings.Contains(part, `yaml:"`) {
			p := strings.Split(part, `"`)
			if strings.Contains(part, `,omitempty`) {
				p[1] = strings.Replace(p[1], ",omitempty", "", 1)
			}
			specName = p[1]
		}
	}
	fileName := longName
	if longName == "" {
		fileName = specName
	}

	// fmt.Println(parts[0], parts[1])
	return Attribute{
		Name:       parts[0],
		Type:       parts[1],
		Long:       longName,
		FileName:   strings.ReplaceAll(fileName, "-", "_"),
		ENV:        envName,
		Short:      shortName,
		Default:    defaultName,
		Group:      group,
		SpecName:   specName,
		Save:       save,
		Deprecated: deprecated,
		Example:    example,
	}, nil
}

func fmtFile(filename string) {
	cmd := exec.Command("gofmt", "-s", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
