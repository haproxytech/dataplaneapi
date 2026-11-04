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

package map_entries

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/client-native/v6/storage"
)

type fakeClient struct {
	client_native.HAProxyClient
	runtime cn_runtime.Runtime
	storage storage.Storage
}

func (c *fakeClient) Runtime() (cn_runtime.Runtime, error) {
	return c.runtime, nil
}

func (c *fakeClient) MapStorage() (storage.Storage, error) {
	return c.storage, nil
}

type fakeRuntime struct {
	cn_runtime.Runtime
	replaceErr      error
	mapFile         string
	mapsDir         string
	replaceCalls    int
	replacedName    string
	replacedEntries models.MapEntries
}

func (r *fakeRuntime) ReplaceMapPayload(name string, entries models.MapEntries) error {
	r.replaceCalls++
	r.replacedName = name
	r.replacedEntries = entries
	return r.replaceErr
}

func (r *fakeRuntime) GetMap(string) (*models.Map, error) {
	return &models.Map{File: r.mapFile}, nil
}

func (r *fakeRuntime) GetMapsDir() (string, error) {
	return r.mapsDir, nil
}

type fakeStorage struct {
	storage.Storage
	replaceErr error
	calls      int
	name       string
	contents   string
}

func (s *fakeStorage) Replace(name, contents string) (string, error) {
	s.calls++
	s.name = name
	s.contents = contents
	return name, s.replaceErr
}

func TestReplaceRuntimeMapEntries(t *testing.T) {
	tests := []struct {
		name                string
		body                string
		forceSync           bool
		runtime             *fakeRuntime
		storage             *fakeStorage
		wantStatus          int
		wantRuntimeCalls    int
		wantStorageCalls    int
		wantStorageContents string
		wantRuntimeEntries  models.MapEntries
	}{
		{
			name:      "preserves entry order in storage",
			body:      `[{"key":"second","value":"two"},{"key":"first","value":"one"}]`,
			forceSync: true,
			runtime: &fakeRuntime{
				mapFile: "/var/lib/haproxy/maps/example.map",
				mapsDir: "/var/lib/haproxy/maps",
			},
			storage:             &fakeStorage{},
			wantStatus:          http.StatusOK,
			wantRuntimeCalls:    1,
			wantStorageCalls:    1,
			wantStorageContents: "second two\nfirst one\n",
			wantRuntimeEntries: models.MapEntries{
				&models.MapEntry{Key: "second", Value: "two"},
				&models.MapEntry{Key: "first", Value: "one"},
			},
		},
		{
			name:      "empties runtime and storage",
			body:      `[]`,
			forceSync: true,
			runtime: &fakeRuntime{
				mapFile: "/var/lib/haproxy/maps/example.map",
				mapsDir: "/var/lib/haproxy/maps",
			},
			storage:             &fakeStorage{},
			wantStatus:          http.StatusOK,
			wantRuntimeCalls:    1,
			wantStorageCalls:    1,
			wantStorageContents: "",
			wantRuntimeEntries:  models.MapEntries{},
		},
		{
			name:      "does not write storage after runtime failure",
			body:      `[{"key":"key","value":"value"}]`,
			forceSync: true,
			runtime: &fakeRuntime{
				replaceErr: errors.New("runtime replacement failed"),
				mapFile:    "/var/lib/haproxy/maps/example.map",
				mapsDir:    "/var/lib/haproxy/maps",
			},
			storage:            &fakeStorage{},
			wantStatus:         http.StatusInternalServerError,
			wantRuntimeCalls:   1,
			wantStorageCalls:   0,
			wantRuntimeEntries: models.MapEntries{&models.MapEntry{Key: "key", Value: "value"}},
		},
		{
			name:      "returns error after storage failure",
			body:      `[{"key":"key","value":"value"}]`,
			forceSync: true,
			runtime: &fakeRuntime{
				mapFile: "/var/lib/haproxy/maps/example.map",
				mapsDir: "/var/lib/haproxy/maps",
			},
			storage: &fakeStorage{
				replaceErr: errors.New("storage replacement failed"),
			},
			wantStatus:          http.StatusInternalServerError,
			wantRuntimeCalls:    1,
			wantStorageCalls:    1,
			wantStorageContents: "key value\n",
			wantRuntimeEntries:  models.MapEntries{&models.MapEntry{Key: "key", Value: "value"}},
		},
		{
			name:      "rejects force sync for unmanaged map",
			body:      `[{"key":"key","value":"value"}]`,
			forceSync: true,
			runtime: &fakeRuntime{
				mapFile: "/etc/haproxy/example.map",
				mapsDir: "/var/lib/haproxy/maps",
			},
			storage:            &fakeStorage{},
			wantStatus:         http.StatusBadRequest,
			wantRuntimeCalls:   0,
			wantStorageCalls:   0,
			wantRuntimeEntries: nil,
		},
		{
			name:      "rejects force sync without maps directory",
			body:      `[{"key":"key","value":"value"}]`,
			forceSync: true,
			runtime: &fakeRuntime{
				mapFile: "example.map",
			},
			storage:            &fakeStorage{},
			wantStatus:         http.StatusBadRequest,
			wantRuntimeCalls:   0,
			wantStorageCalls:   0,
			wantRuntimeEntries: nil,
		},
		{
			name:               "skips storage without force sync",
			body:               `[{"key":"key","value":"value"}]`,
			runtime:            &fakeRuntime{},
			storage:            &fakeStorage{},
			wantStatus:         http.StatusOK,
			wantRuntimeCalls:   1,
			wantStorageCalls:   0,
			wantRuntimeEntries: models.MapEntries{&models.MapEntry{Key: "key", Value: "value"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			client := &fakeClient{runtime: tt.runtime, storage: tt.storage}
			if err := RegisterRouter(router, client); err != nil {
				t.Fatalf("RegisterRouter() error = %v", err)
			}

			url := "/services/haproxy/runtime/maps/example.map/entries"
			if tt.forceSync {
				url += "?force_sync=true"
			}
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, tt.wantStatus, recorder.Body.String())
			}
			if tt.runtime.replaceCalls != tt.wantRuntimeCalls {
				t.Errorf("runtime replacement calls = %d, want %d", tt.runtime.replaceCalls, tt.wantRuntimeCalls)
			}
			if tt.wantRuntimeCalls > 0 && tt.runtime.replacedName != "example.map" {
				t.Errorf("runtime map name = %q, want example.map", tt.runtime.replacedName)
			}
			if !reflect.DeepEqual(tt.runtime.replacedEntries, tt.wantRuntimeEntries) {
				t.Errorf("runtime entries = %#v, want %#v", tt.runtime.replacedEntries, tt.wantRuntimeEntries)
			}
			if tt.storage.calls != tt.wantStorageCalls {
				t.Errorf("storage calls = %d, want %d", tt.storage.calls, tt.wantStorageCalls)
			}
			if tt.storage.contents != tt.wantStorageContents {
				t.Errorf("storage contents = %q, want %q", tt.storage.contents, tt.wantStorageContents)
			}
			if tt.wantStorageCalls > 0 && tt.storage.name != "example.map" {
				t.Errorf("storage name = %q, want example.map", tt.storage.name)
			}
		})
	}
}

func TestReplaceRuntimeMapEntriesRejectsNullEntry(t *testing.T) {
	rt := &fakeRuntime{}
	handler := &HandlerImpl{
		Client: &fakeClient{
			runtime: rt,
			storage: &fakeStorage{},
		},
	}
	req := httptest.NewRequest(http.MethodPut, "/services/haproxy/runtime/maps/example.map/entries", strings.NewReader(`[null]`))
	recorder := httptest.NewRecorder()

	handler.ReplaceRuntimeMapEntries(recorder, req, "example.map", ReplaceRuntimeMapEntriesParams{})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}
	if rt.replaceCalls != 0 {
		t.Errorf("runtime replacement calls = %d, want 0", rt.replaceCalls)
	}
}
