// Code generated by go-swagger; DO NOT EDIT.

package filter

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

// NewDeleteFilterParams creates a new DeleteFilterParams object
// with the default values initialized.
func NewDeleteFilterParams() *DeleteFilterParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &DeleteFilterParams{
		ForceReload: &forceReloadDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewDeleteFilterParamsWithTimeout creates a new DeleteFilterParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDeleteFilterParamsWithTimeout(timeout time.Duration) *DeleteFilterParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &DeleteFilterParams{
		ForceReload: &forceReloadDefault,

		timeout: timeout,
	}
}

// NewDeleteFilterParamsWithContext creates a new DeleteFilterParams object
// with the default values initialized, and the ability to set a context for a request
func NewDeleteFilterParamsWithContext(ctx context.Context) *DeleteFilterParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &DeleteFilterParams{
		ForceReload: &forceReloadDefault,

		Context: ctx,
	}
}

// NewDeleteFilterParamsWithHTTPClient creates a new DeleteFilterParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDeleteFilterParamsWithHTTPClient(client *http.Client) *DeleteFilterParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &DeleteFilterParams{
		ForceReload: &forceReloadDefault,
		HTTPClient:  client,
	}
}

/*DeleteFilterParams contains all the parameters to send to the API endpoint
for the delete filter operation typically these are written to a http.Request
*/
type DeleteFilterParams struct {

	/*ForceReload
	  If set, do a force reload, do not wait for the configured reload-delay. Cannot be used when transaction is specified, as changes in transaction are not applied directly to configuration.

	*/
	ForceReload *bool
	/*ID
	  Filter ID

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
	/*Version
	  Version used for checking configuration version. Cannot be used when transaction is specified, transaction has it's own version.

	*/
	Version *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the delete filter params
func (o *DeleteFilterParams) WithTimeout(timeout time.Duration) *DeleteFilterParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the delete filter params
func (o *DeleteFilterParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the delete filter params
func (o *DeleteFilterParams) WithContext(ctx context.Context) *DeleteFilterParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the delete filter params
func (o *DeleteFilterParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the delete filter params
func (o *DeleteFilterParams) WithHTTPClient(client *http.Client) *DeleteFilterParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the delete filter params
func (o *DeleteFilterParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithForceReload adds the forceReload to the delete filter params
func (o *DeleteFilterParams) WithForceReload(forceReload *bool) *DeleteFilterParams {
	o.SetForceReload(forceReload)
	return o
}

// SetForceReload adds the forceReload to the delete filter params
func (o *DeleteFilterParams) SetForceReload(forceReload *bool) {
	o.ForceReload = forceReload
}

// WithID adds the id to the delete filter params
func (o *DeleteFilterParams) WithID(id int64) *DeleteFilterParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the delete filter params
func (o *DeleteFilterParams) SetID(id int64) {
	o.ID = id
}

// WithParentName adds the parentName to the delete filter params
func (o *DeleteFilterParams) WithParentName(parentName string) *DeleteFilterParams {
	o.SetParentName(parentName)
	return o
}

// SetParentName adds the parentName to the delete filter params
func (o *DeleteFilterParams) SetParentName(parentName string) {
	o.ParentName = parentName
}

// WithParentType adds the parentType to the delete filter params
func (o *DeleteFilterParams) WithParentType(parentType string) *DeleteFilterParams {
	o.SetParentType(parentType)
	return o
}

// SetParentType adds the parentType to the delete filter params
func (o *DeleteFilterParams) SetParentType(parentType string) {
	o.ParentType = parentType
}

// WithTransactionID adds the transactionID to the delete filter params
func (o *DeleteFilterParams) WithTransactionID(transactionID *string) *DeleteFilterParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the delete filter params
func (o *DeleteFilterParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WithVersion adds the version to the delete filter params
func (o *DeleteFilterParams) WithVersion(version *int64) *DeleteFilterParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the delete filter params
func (o *DeleteFilterParams) SetVersion(version *int64) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *DeleteFilterParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

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
