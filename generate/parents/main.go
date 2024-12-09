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
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-openapi/swag"
	cnparents "github.com/haproxytech/client-native/v6/configuration/parents"
)

// To generate parent aliases parent_<childType>_generated.go in handler/
// Usage: go run generate/parents/*.go
// To generate for a new childType, update 2 files:
// - operations.go for the list of operations to generate (Create/Get/...)

func main() {
	fmt.Println("Generating parent aliases")
	children := []string{
		cnparents.ServerChildType,
		cnparents.HTTPAfterResponseRuleChildType,
		cnparents.HTTPCheckChildType,
		cnparents.HTTPErrorRuleChildType,
		cnparents.HTTPRequestRuleChildType,
		cnparents.HTTPResponseRuleChildType,
		cnparents.TCPCheckChildType,
		cnparents.TCPRequestRuleChildType,
		cnparents.TCPResponseRuleChildType,
		cnparents.QUICInitialRuleType,
		cnparents.ACLChildType,
		cnparents.BindChildType,
		cnparents.FilterChildType,
		cnparents.LogTargetChildType,
	}
	for _, childType := range children {
		res := generateAlias(childType)

		// Create or open a file
		fileName := fmt.Sprintf("handlers/parent_%s_generated.go", strings.ToLower(childType))
		fmt.Printf("Generated %s\n", fileName)
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Write the buffer's content to the file
		_, err = res.WriteTo(file)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

type TmplData struct {
	GoChildType      string
	ChildType        string
	Operations       []string
	OperationPackage string
}

func generateAlias(childType string) bytes.Buffer {
	// Initialisms used in child resources need to be added here for the generated parent functions to match with the operations params
	swag.AddInitialisms("QUIC")

	funcMap := template.FuncMap{
		"parents": cnparents.Parents,
	}
	parents := cnparents.Parents(childType)
	fmt.Printf("Generating for child %s / parents: %v\n", childType, parents)

	tmplData := TmplData{
		ChildType:        childType,
		GoChildType:      swag.ToGoName(childType),
		Operations:       operations(childType),
		OperationPackage: childType,
	}
	templateName := "parent_generated.tmpl"
	tmpl := template.Must(template.New(templateName).Funcs(funcMap).ParseFiles("generate/parents/" + templateName))
	var result bytes.Buffer
	err := tmpl.ExecuteTemplate(&result, templateName, tmplData)
	if err != nil {
		exitError(err.Error())
	}
	return result
}

func exitError(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", msg)
	os.Exit(1)
}
