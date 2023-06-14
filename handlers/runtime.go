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

package handlers

import (
	"reflect"
	"strconv"

	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/log"
)

// RuntimeSupportedFields is a map of fields supported through the runtime API for
// it's respectable object type
var RuntimeSupportedFields = map[string][]string{
	"frontend": {"Maxconn"},
	"server":   {"Weight", "Address", "Port", "Maintenance", "AgentCheck", "AgentAddr", "AgentSend", "HealthCheckPort"},
}

// ChangeThroughRuntimeAPI checks if something can be changed through the runtime API, and
// returns false if reload is not needed, or true if needed.
func changeThroughRuntimeAPI(data, ondisk interface{}, parentName, parentType string, client client_native.HAProxyClient) (reload bool) {
	// reflect kinds and values are loosely checked as they are bound strictly in
	// schema, but in case of any panic, we will log and reload to ensure
	// changes go through
	defer func() {
		if r := recover(); r != nil {
			log.Warning("Change Through API Panic: ", r)
			// we are panicking, so reload to ensure changes are through
			reload = true
		}
	}()

	// check if runtime client is valid, if not, return reload needed
	runtime, err := client.Runtime()
	if err != nil {
		return true
	}
	// objects are the same, do nothing and don't reload
	diff := compareObjects(data, ondisk)
	if len(diff) == 0 {
		return false
	}
	switch vData := data.(type) {
	case models.Frontend:
		if compareChanged(diff, RuntimeSupportedFields["frontend"]) {
			// only fields that are supported in runtime API are changed, if successful, do not reload
			for _, field := range diff {
				fieldValue := reflect.ValueOf(vData).FieldByName(field)
				if fieldValue.IsValid() {
					//nolint:gocritic
					switch field {
					case "Maxconn":
						maxConn := fieldValue.Elem().Int()
						err := runtime.SetFrontendMaxConn(vData.Name, int(maxConn))
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					}
				}
			}
		} else {
			// there are changed fields that are not supported in runtime API, so reload
			return true
		}
	case models.Server:
		if compareChanged(diff, RuntimeSupportedFields["server"]) {
			// only fields that are supported in runtime API are changed, if successful, do not reload
			addrPortChanged := false
			for _, field := range diff {
				fieldValue := reflect.ValueOf(vData).FieldByName(field)
				if fieldValue.IsValid() {
					switch field {
					case "Weight":
						weight := strconv.FormatInt(fieldValue.Elem().Int(), 10)
						err := runtime.SetServerWeight(parentName, vData.Name, weight)
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					case "HealthCheckPort":
						portVal := fieldValue.Elem().Int()
						err := runtime.SetServerCheckPort(parentName, vData.Name, int(portVal))
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					case "Address":
						if !addrPortChanged {
							portVal := reflect.ValueOf(vData).FieldByName("Port").Elem().Int()
							addrVal := fieldValue.String()
							err := runtime.SetServerAddr(parentName, vData.Name, addrVal, int(portVal))
							if err != nil {
								// we have error's in runtime API changes, bail and reload
								return true
							}
							addrPortChanged = true
						}
					case "Port":
						if !addrPortChanged {
							portVal := fieldValue.Elem().Int()
							addrVal := reflect.ValueOf(vData).FieldByName("Address").String()
							err := runtime.SetServerAddr(parentName, vData.Name, addrVal, int(portVal))
							if err != nil {
								// we have error's in runtime API changes, bail and reload
								return true
							}
							addrPortChanged = true
						}
					case "Maintenance":
						maint := fieldValue.String()
						if maint == "enabled" {
							maint = "maint"
						} else {
							maint = "ready"
						}
						err := runtime.SetServerState(parentName, vData.Name, maint)
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					case "AgentCheck":
						aCheck := fieldValue.String()
						if aCheck == "enabled" {
							err := runtime.EnableAgentCheck(parentName, vData.Name)
							if err != nil {
								// we have error's in runtime API changes, bail and reload
								return true
							}
						}
						if aCheck == "disabled" {
							err := runtime.DisableAgentCheck(parentName, vData.Name)
							if err != nil {
								// we have error's in runtime API changes, bail and reload
								return true
							}
						}
					case "AgentAddr":
						aAddr := fieldValue.String()
						err := runtime.SetServerAgentAddr(parentName, vData.Name, aAddr)
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					case "AgentSend":
						aSend := fieldValue.String()
						err := runtime.SetServerAgentSend(parentName, vData.Name, aSend)
						if err != nil {
							// we have error's in runtime API changes, bail and reload
							return true
						}
					}
				}
			}
		} else {
			// there are changed fields that are not supported in runtime API, so reload
			return true
		}
	default:
		// for this type, we do not support changing through runtime API, so reload
		return true
	}
	return false
}

// return string of field names that have a diff
func compareObjects(data, ondisk interface{}) []string {
	diff := make([]string, 0)
	dataVal := reflect.ValueOf(data)
	ondiskVal := reflect.ValueOf(ondisk)
	for i := 0; i < dataVal.NumField(); i++ {
		fName := dataVal.Type().Field(i).Name
		dataField := dataVal.FieldByName(fName)
		ondiskField := ondiskVal.FieldByName(fName)

		dataKind := dataField.Kind()
		ondiskKind := ondiskField.Kind()

		if dataKind != ondiskKind {
			diff = append(diff, fName)
			continue
		}

		if dataKind == reflect.Ptr {
			dataField = dataField.Elem()
			ondiskField = ondiskField.Elem()

			dataKind = dataField.Kind()
			ondiskKind = ondiskField.Kind()

			if dataKind != ondiskKind {
				diff = append(diff, fName)
				continue
			}
		}

		switch dataKind {
		case reflect.Float32, reflect.Float64:
			if dataField.Float() != ondiskField.Float() {
				diff = append(diff, fName)
			}
		case reflect.Bool:
			if dataField.Bool() != ondiskField.Bool() {
				diff = append(diff, fName)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if dataField.Int() != ondiskField.Int() {
				diff = append(diff, fName)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if dataField.Uint() != ondiskField.Uint() {
				diff = append(diff, fName)
			}
		case reflect.String:
			if dataField.String() != ondiskField.String() {
				diff = append(diff, fName)
			}
		case reflect.Struct:
			diff = append(diff, compareObjects(dataField.Interface(), ondiskField.Interface())...)
		}
	}
	return diff
}

// this returns true if only changeable fields have been changed
func compareChanged(changed, changeable []string) bool {
	if len(changed) > len(changeable) {
		return false
	}

	diff := make(map[string]bool, len(changed))
	for _, elem := range changed {
		diff[elem] = true
	}

	for _, elem := range changeable {
		delete(diff, elem)
	}

	return len(diff) == 0
}
