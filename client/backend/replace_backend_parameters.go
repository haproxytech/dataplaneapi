// Code generated by go-swagger; DO NOT EDIT.

package backend

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

	"github.com/haproxytech/models"
)

// NewReplaceBackendParams creates a new ReplaceBackendParams object
// with the default values initialized.
func NewReplaceBackendParams() *ReplaceBackendParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &ReplaceBackendParams{
		ForceReload: &forceReloadDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewReplaceBackendParamsWithTimeout creates a new ReplaceBackendParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewReplaceBackendParamsWithTimeout(timeout time.Duration) *ReplaceBackendParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &ReplaceBackendParams{
		ForceReload: &forceReloadDefault,

		timeout: timeout,
	}
}

// NewReplaceBackendParamsWithContext creates a new ReplaceBackendParams object
// with the default values initialized, and the ability to set a context for a request
func NewReplaceBackendParamsWithContext(ctx context.Context) *ReplaceBackendParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &ReplaceBackendParams{
		ForceReload: &forceReloadDefault,

		Context: ctx,
	}
}

// NewReplaceBackendParamsWithHTTPClient creates a new ReplaceBackendParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewReplaceBackendParamsWithHTTPClient(client *http.Client) *ReplaceBackendParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &ReplaceBackendParams{
		ForceReload: &forceReloadDefault,
		HTTPClient:  client,
	}
}

/*ReplaceBackendParams contains all the parameters to send to the API endpoint
for the replace backend operation typically these are written to a http.Request
*/
type ReplaceBackendParams struct {

	/*Data*/
	Data *models.Backend
	/*ForceReload
	  If set, do a force reload, do not wait for the configured reload-delay. Cannot be used when transaction is specified, as changes in transaction are not applied directly to configuration.

	*/
	ForceReload *bool
	/*Name
	  Backend name

	*/
	Name string
	/*TransactionID
	  ID of the transaction where we want to add the operation. Cannot be used when version is specified.

	*/
	TransactionID *string
	/*Version
	  Version used for checking configuration version. Cannot be used when transaction is specified, transaction has it's own version.

	*/
	Version *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the replace backend params
func (o *ReplaceBackendParams) WithTimeout(timeout time.Duration) *ReplaceBackendParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the replace backend params
func (o *ReplaceBackendParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the replace backend params
func (o *ReplaceBackendParams) WithContext(ctx context.Context) *ReplaceBackendParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the replace backend params
func (o *ReplaceBackendParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the replace backend params
func (o *ReplaceBackendParams) WithHTTPClient(client *http.Client) *ReplaceBackendParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the replace backend params
func (o *ReplaceBackendParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithData adds the data to the replace backend params
func (o *ReplaceBackendParams) WithData(data *models.Backend) *ReplaceBackendParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the replace backend params
func (o *ReplaceBackendParams) SetData(data *models.Backend) {
	o.Data = data
}

// WithForceReload adds the forceReload to the replace backend params
func (o *ReplaceBackendParams) WithForceReload(forceReload *bool) *ReplaceBackendParams {
	o.SetForceReload(forceReload)
	return o
}

// SetForceReload adds the forceReload to the replace backend params
func (o *ReplaceBackendParams) SetForceReload(forceReload *bool) {
	o.ForceReload = forceReload
}

// WithName adds the name to the replace backend params
func (o *ReplaceBackendParams) WithName(name string) *ReplaceBackendParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the replace backend params
func (o *ReplaceBackendParams) SetName(name string) {
	o.Name = name
}

// WithTransactionID adds the transactionID to the replace backend params
func (o *ReplaceBackendParams) WithTransactionID(transactionID *string) *ReplaceBackendParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the replace backend params
func (o *ReplaceBackendParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WithVersion adds the version to the replace backend params
func (o *ReplaceBackendParams) WithVersion(version *int64) *ReplaceBackendParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the replace backend params
func (o *ReplaceBackendParams) SetVersion(version *int64) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *ReplaceBackendParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Data != nil {
		if err := r.SetBodyParam(o.Data); err != nil {
			return err
		}
	}

	if o.ForceReload != nil {

		// query param force_reload
		var qrForceReload bool
		if o.ForceReload != nil {
			qrForceReload = *o.ForceReload
		}
		qForceReload := swag.FormatBool(qrForceReload)
		if qForceReload != "" {
			if err := r.SetQueryParam("force_reload", qForceReload); err != nil {
				return err
			}
		}

	}

	// path param name
	if err := r.SetPathParam("name", o.Name); err != nil {
		return err
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

	if o.Version != nil {

		// query param version
		var qrVersion int64
		if o.Version != nil {
			qrVersion = *o.Version
		}
		qVersion := swag.FormatInt64(qrVersion)
		if qVersion != "" {
			if err := r.SetQueryParam("version", qVersion); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
