// Code generated by go-swagger; DO NOT EDIT.

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

package defaults

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"
)

// GetDefaultsSectionURL generates an URL for the get defaults section operation
type GetDefaultsSectionURL struct {
	Name string

	TransactionID *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetDefaultsSectionURL) WithBasePath(bp string) *GetDefaultsSectionURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetDefaultsSectionURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetDefaultsSectionURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/services/haproxy/configuration/named_defaults/{name}"

	name := o.Name
	if name != "" {
		_path = strings.Replace(_path, "{name}", name, -1)
	} else {
		return nil, errors.New("name is required on GetDefaultsSectionURL")
	}

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v3"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var transactionIDQ string
	if o.TransactionID != nil {
		transactionIDQ = *o.TransactionID
	}
	if transactionIDQ != "" {
		qs.Set("transaction_id", transactionIDQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetDefaultsSectionURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetDefaultsSectionURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetDefaultsSectionURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetDefaultsSectionURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetDefaultsSectionURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetDefaultsSectionURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
