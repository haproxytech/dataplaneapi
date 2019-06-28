// Code generated by go-swagger; DO NOT EDIT.

package tcp_request_rule

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTCPRequestRulesParams creates a new GetTCPRequestRulesParams object
// with the default values initialized.
func NewGetTCPRequestRulesParams() *GetTCPRequestRulesParams {
	var ()
	return &GetTCPRequestRulesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetTCPRequestRulesParamsWithTimeout creates a new GetTCPRequestRulesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetTCPRequestRulesParamsWithTimeout(timeout time.Duration) *GetTCPRequestRulesParams {
	var ()
	return &GetTCPRequestRulesParams{

		timeout: timeout,
	}
}

// NewGetTCPRequestRulesParamsWithContext creates a new GetTCPRequestRulesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetTCPRequestRulesParamsWithContext(ctx context.Context) *GetTCPRequestRulesParams {
	var ()
	return &GetTCPRequestRulesParams{

		Context: ctx,
	}
}

// NewGetTCPRequestRulesParamsWithHTTPClient creates a new GetTCPRequestRulesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetTCPRequestRulesParamsWithHTTPClient(client *http.Client) *GetTCPRequestRulesParams {
	var ()
	return &GetTCPRequestRulesParams{
		HTTPClient: client,
	}
}

/*GetTCPRequestRulesParams contains all the parameters to send to the API endpoint
for the get TCP request rules operation typically these are written to a http.Request
*/
type GetTCPRequestRulesParams struct {

	/*ParentName
	  Parent name

	*/
	ParentName string
	/*ParentType
	  Parent type

	*/
	ParentType string
	/*TransactionID
	  ID of the transaction where we want to add the operation. Cannot be used when version is specified.

	*/
	TransactionID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithTimeout(timeout time.Duration) *GetTCPRequestRulesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithContext(ctx context.Context) *GetTCPRequestRulesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithHTTPClient(client *http.Client) *GetTCPRequestRulesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithParentName adds the parentName to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithParentName(parentName string) *GetTCPRequestRulesParams {
	o.SetParentName(parentName)
	return o
}

// SetParentName adds the parentName to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetParentName(parentName string) {
	o.ParentName = parentName
}

// WithParentType adds the parentType to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithParentType(parentType string) *GetTCPRequestRulesParams {
	o.SetParentType(parentType)
	return o
}

// SetParentType adds the parentType to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetParentType(parentType string) {
	o.ParentType = parentType
}

// WithTransactionID adds the transactionID to the get TCP request rules params
func (o *GetTCPRequestRulesParams) WithTransactionID(transactionID *string) *GetTCPRequestRulesParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the get TCP request rules params
func (o *GetTCPRequestRulesParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetTCPRequestRulesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param parent_name
	qrParentName := o.ParentName
	qParentName := qrParentName
	if qParentName != "" {
		if err := r.SetQueryParam("parent_name", qParentName); err != nil {
			return err
		}
	}

	// query param parent_type
	qrParentType := o.ParentType
	qParentType := qrParentType
	if qParentType != "" {
		if err := r.SetQueryParam("parent_type", qParentType); err != nil {
			return err
		}
	}

	if o.TransactionID != nil {

		// query param transaction_id
		var qrTransactionID string
		if o.TransactionID != nil {
			qrTransactionID = *o.TransactionID
		}
		qTransactionID := qrTransactionID
		if qTransactionID != "" {
			if err := r.SetQueryParam("transaction_id", qTransactionID); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
