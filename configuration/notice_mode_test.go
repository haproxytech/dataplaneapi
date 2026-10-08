// Copyright 2026 HAProxy Technologies
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

package configuration

import (
	"os"
	"path/filepath"
	"testing"
)

// The bootstrap key storage-dir NOTICE file must not be world-writable.
func TestInitStorageNoticeFile_NotWorldWritable(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "anywhere", "nested")
	if err := CheckIfStorageDirIsOK(dir, nil); err != nil {
		t.Fatal(err)
	}
	if err := InitStorageNoticeFile(dir); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(filepath.Join(dir, "NOTICE"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o022 != 0 {
		t.Fatalf("NOTICE writable by group/others: %v", fi.Mode().Perm())
	}
}
