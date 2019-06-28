// Code generated by go-swagger; DO NOT EDIT.

package acl

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetACLParams creates a new GetACLParams object
// with the default values initialized.
func NewGetACLParams() *GetACLParams {
	var ()
	return &GetACLParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetACLParamsWithTimeout creates a new GetACLParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetACLParamsWithTimeout(timeout time.Duration) *GetACLParams {
	var ()
	return &GetACLParams{

		timeout: timeout,
	}
}

// NewGetACLParamsWithContext creates a new GetACLParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetACLParamsWithContext(ctx context.Context) *GetACLParams {
	var ()
	return &GetACLParams{

		Context: ctx,
	}
}

// NewGetACLParamsWithHTTPClient creates a new GetACLParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetACLParamsWithHTTPClient(client *http.Client) *GetACLParams {
	var ()
	return &GetACLParams{
		HTTPClient: client,
	}
}

/*GetACLParams contains all the parameters to send to the API endpoint
for the get Acl operation typically these are written to a http.Request
*/
type GetACLParams struct {

	/*ID
	  ACL line ID

	*/
	ID int64
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

// WithTimeout adds the timeout to the get Acl params
func (o *GetACLParams) WithTimeout(timeout time.Duration) *GetACLParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get Acl params
func (o *GetACLParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get Acl params
func (o *GetACLParams) WithContext(ctx context.Context) *GetACLParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get Acl params
func (o *GetACLParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get Acl params
func (o *GetACLParams) WithHTTPClient(client *http.Client) *GetACLParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get Acl params
func (o *GetACLParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get Acl params
func (o *GetACLParams) WithID(id int64) *GetACLParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get Acl params
func (o *GetACLParams) SetID(id int64) {
	o.ID = id
}

// WithParentName adds the parentName to the get Acl params
func (o *GetACLParams) WithParentName(parentName string) *GetACLParams {
	o.SetParentName(parentName)
	return o
}

// SetParentName adds the parentName to the get Acl params
func (o *GetACLParams) SetParentName(parentName string) {
	o.ParentName = parentName
}

// WithParentType adds the parentType to the get Acl params
func (o *GetACLParams) WithParentType(parentType string) *GetACLParams {
	o.SetParentType(parentType)
	return o
}

// SetParentType adds the parentType to the get Acl params
func (o *GetACLParams) SetParentType(parentType string) {
	o.ParentType = parentType
}

// WithTransactionID adds the transactionID to the get Acl params
func (o *GetACLParams) WithTransactionID(transactionID *string) *GetACLParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the get Acl params
func (o *GetACLParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetACLParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

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
