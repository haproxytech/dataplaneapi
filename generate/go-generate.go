// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	"reload":      "haproxy",
	"syslog":      "log",
}

type Attribute struct {
	Name       string
	AttName    string
	Type       string
	Long       string
	FileName   string
	Short      string
	Default    string
	ENV        string
	Group      string
	StructName string
	SpecName   string
	Save       bool
	Deprecated bool
}

type ParseGroup struct {
	Name        string
	AttName     string
	MaxSize     int
	MaxTypeSize int
	Parent      string
	Elements    []string
	Attributes  []Attribute
}

type ParseData struct {
	Groups []ParseGroup
}

func readServerData(filePath string, pd *ParseData, structName string, attName string, groupName string) {
	typeStruct := fmt.Sprintf("type %s struct {", structName)
	dat, err := ioutil.ReadFile(filePath)
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
					Name:        res.Group,
					AttName:     attName,
					Parent:      groupParents[res.Group],
					MaxSize:     len(res.Name),
					MaxTypeSize: len(res.Type),
					Attributes:  []Attribute{res},
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
	readServerData(filePath, pd, "Server", "", "")
	// ######################################## configuration.go
	filePath = path.Join(dir, "configuration", "configuration.go", "")
	readServerData(filePath, pd, "Configuration", "-", "")
	readServerData(filePath, pd, "HAProxyConfiguration", "HAProxy", "")
	readServerData(filePath, pd, "APIConfiguration", "APIOptions", "")
	readServerData(filePath, pd, "ServiceDiscovery", "ServiceDiscovery", "")
	readServerData(filePath, pd, "ClusterConfiguration", "Cluster", "cluster")
	readServerData(filePath, pd, "SyslogOptions", "Syslog", "")
	readServerData(filePath, pd, "LoggingOptions", "Logging", "")
	// ########################################

	// prepare template function
	funcMap := template.FuncMap{
		"Capitalize": capitalize,
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
		FileName:   fileName,
		ENV:        envName,
		Short:      shortName,
		Default:    defaultName,
		Group:      group,
		SpecName:   specName,
		Save:       save,
		Deprecated: deprecated,
	}, nil
}

func fmtFile(filename string) {
	cmd := exec.Command("gofmt", "-s", "-w", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
