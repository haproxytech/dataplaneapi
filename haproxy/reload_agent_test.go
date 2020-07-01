// Copyright 2020 HAProxy Technologies
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

package haproxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_copyFile(t *testing.T) {
	createFile := func() (string, error) {
		tmpFile, err := ioutil.TempFile("", "")
		if err != nil {
			return "", err
		}
		defer tmpFile.Close()
		if _, err := tmpFile.Write([]byte("Hello, world.")); err != nil {
			return "", err
		}
		if err := tmpFile.Sync(); err != nil {
			return "", err
		}
		return tmpFile.Name(), nil
	}
	srcPath, err := createFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcPath)

	dstPath := srcPath + ".2"
	defer os.RemoveAll(dstPath)

	doTheCopy := func() error {
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
		srcContents, err := ioutil.ReadFile(srcPath)
		if err != nil {
			return err
		}
		dstContents, err := ioutil.ReadFile(dstPath)
		if err != nil {
			return err
		}
		if result := bytes.Compare(srcContents, dstContents); result != 0 {
			return fmt.Errorf("files are not same: %v", result)
		}
		return nil
	}

	if err := doTheCopy(); err != nil {
		t.Fatal(err)
	}
	// Copy the file a second time.
	if err := doTheCopy(); err != nil {
		t.Fatal(err)
	}
}
