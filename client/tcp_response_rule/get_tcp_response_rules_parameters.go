// Code generated by go-swagger; DO NOT EDIT.

package tcp_response_rule

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

// NewGetTCPResponseRulesParams creates a new GetTCPResponseRulesParams object
// with the default values initialized.
func NewGetTCPResponseRulesParams() *GetTCPResponseRulesParams {
	var ()
	return &GetTCPResponseRulesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetTCPResponseRulesParamsWithTimeout creates a new GetTCPResponseRulesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetTCPResponseRulesParamsWithTimeout(timeout time.Duration) *GetTCPResponseRulesParams {
	var ()
	return &GetTCPResponseRulesParams{

		timeout: timeout,
	}
}

// NewGetTCPResponseRulesParamsWithContext creates a new GetTCPResponseRulesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetTCPResponseRulesParamsWithContext(ctx context.Context) *GetTCPResponseRulesParams {
	var ()
	return &GetTCPResponseRulesParams{

		Context: ctx,
	}
}

// NewGetTCPResponseRulesParamsWithHTTPClient creates a new GetTCPResponseRulesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetTCPResponseRulesParamsWithHTTPClient(client *http.Client) *GetTCPResponseRulesParams {
	var ()
	return &GetTCPResponseRulesParams{
		HTTPClient: client,
	}
}

/*GetTCPResponseRulesParams contains all the parameters to send to the API endpoint
for the get TCP response rules operation typically these are written to a http.Request
*/
type GetTCPResponseRulesParams struct {

	/*Backend
	  Parent backend name

	*/
	Backend string
	/*TransactionID
	  ID of the transaction where we want to add the operation. Cannot be used when version is specified.

	*/
	TransactionID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get TCP response rules params
func (o *GetTCPResponseRulesParams) WithTimeout(timeout time.Duration) *GetTCPResponseRulesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get TCP response rules params
func (o *GetTCPResponseRulesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get TCP response rules params
func (o *GetTCPResponseRulesParams) WithContext(ctx context.Context) *GetTCPResponseRulesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get TCP response rules params
func (o *GetTCPResponseRulesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get TCP response rules params
func (o *GetTCPResponseRulesParams) WithHTTPClient(client *http.Client) *GetTCPResponseRulesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get TCP response rules params
func (o *GetTCPResponseRulesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBackend adds the backend to the get TCP response rules params
func (o *GetTCPResponseRulesParams) WithBackend(backend string) *GetTCPResponseRulesParams {
	o.SetBackend(backend)
	return o
}

// SetBackend adds the backend to the get TCP response rules params
func (o *GetTCPResponseRulesParams) SetBackend(backend string) {
	o.Backend = backend
}

// WithTransactionID adds the transactionID to the get TCP response rules params
func (o *GetTCPResponseRulesParams) WithTransactionID(transactionID *string) *GetTCPResponseRulesParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the get TCP response rules params
func (o *GetTCPResponseRulesParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WriteToRequest writes these params to a swagger request
func (o *GetTCPResponseRulesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param backend
	qrBackend := o.Backend
	qBackend := qrBackend
	if qBackend != "" {
		if err := r.SetQueryParam("backend", qBackend); err != nil {
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
