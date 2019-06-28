// Code generated by go-swagger; DO NOT EDIT.

package log_target

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new log target API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for log target API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
CreateLogTarget adds a new log target

Adds a new Log Target of the specified type in the specified parent.
*/
func (a *Client) CreateLogTarget(params *CreateLogTargetParams, authInfo runtime.ClientAuthInfoWriter) (*CreateLogTargetCreated, *CreateLogTargetAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewCreateLogTargetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "createLogTarget",
		Method:             "POST",
		PathPattern:        "/services/haproxy/configuration/log_targets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &CreateLogTargetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *CreateLogTargetCreated:
		return value, nil, nil
	case *CreateLogTargetAccepted:
		return nil, value, nil
	}
	return nil, nil, nil

}

/*
DeleteLogTarget deletes a log target

Deletes a Log Target configuration by it's ID from the specified parent.
*/
func (a *Client) DeleteLogTarget(params *DeleteLogTargetParams, authInfo runtime.ClientAuthInfoWriter) (*DeleteLogTargetAccepted, *DeleteLogTargetNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteLogTargetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteLogTarget",
		Method:             "DELETE",
		PathPattern:        "/services/haproxy/configuration/log_targets/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteLogTargetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *DeleteLogTargetAccepted:
		return value, nil, nil
	case *DeleteLogTargetNoContent:
		return nil, value, nil
	}
	return nil, nil, nil

}

/*
GetLogTarget returns one log target

Returns one Log Target configuration by it's ID in the specified parent.
*/
func (a *Client) GetLogTarget(params *GetLogTargetParams, authInfo runtime.ClientAuthInfoWriter) (*GetLogTargetOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetLogTargetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getLogTarget",
		Method:             "GET",
		PathPattern:        "/services/haproxy/configuration/log_targets/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetLogTargetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetLogTargetOK), nil

}

/*
GetLogTargets returns an array of all log targets

Returns all Log Targets that are configured in specified parent.
*/
func (a *Client) GetLogTargets(params *GetLogTargetsParams, authInfo runtime.ClientAuthInfoWriter) (*GetLogTargetsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetLogTargetsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getLogTargets",
		Method:             "GET",
		PathPattern:        "/services/haproxy/configuration/log_targets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetLogTargetsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetLogTargetsOK), nil

}

/*
ReplaceLogTarget replaces a log target

Replaces a Log Target configuration by it's ID in the specified parent.
*/
func (a *Client) ReplaceLogTarget(params *ReplaceLogTargetParams, authInfo runtime.ClientAuthInfoWriter) (*ReplaceLogTargetOK, *ReplaceLogTargetAccepted, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewReplaceLogTargetParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "replaceLogTarget",
		Method:             "PUT",
		PathPattern:        "/services/haproxy/configuration/log_targets/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &ReplaceLogTargetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *ReplaceLogTargetOK:
		return value, nil, nil
	case *ReplaceLogTargetAccepted:
		return nil, value, nil
	}
	return nil, nil, nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
