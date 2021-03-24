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
	"log"
	"os"

	"github.com/haproxytech/dataplaneapi/configuration"
)

func main() {
	if len(os.Args) < 3 {
		log.Panic("not enough params")
	}

	file := os.Args[1]
	output := os.Args[2]

	isHcl := true

	storageHCL := &configuration.StorageHCL{}
	storageYAML := &configuration.StorageYML{}
	err := storageHCL.Load(file)
	_, isPathError := err.(*os.PathError)
	if isPathError {
		log.Panic(err)
	}
	if err != nil {
		isHcl = false
		errYaml := storageYAML.Load(file)
		if errYaml != nil {
			log.Panic(err)
		}
	}

	if isHcl {
		storageYAML.Set(storageHCL.Get())
		err := storageYAML.SaveAs(output)
		if err != nil {
			log.Panic(err)
		}
	} else {
		storageHCL.Set(storageYAML.Get())
		err := storageHCL.SaveAs(output)
		if err != nil {
			log.Panic(err)
		}
	}
}
