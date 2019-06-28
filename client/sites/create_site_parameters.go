// Code generated by go-swagger; DO NOT EDIT.

package sites

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

// NewCreateSiteParams creates a new CreateSiteParams object
// with the default values initialized.
func NewCreateSiteParams() *CreateSiteParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &CreateSiteParams{
		ForceReload: &forceReloadDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateSiteParamsWithTimeout creates a new CreateSiteParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateSiteParamsWithTimeout(timeout time.Duration) *CreateSiteParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &CreateSiteParams{
		ForceReload: &forceReloadDefault,

		timeout: timeout,
	}
}

// NewCreateSiteParamsWithContext creates a new CreateSiteParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateSiteParamsWithContext(ctx context.Context) *CreateSiteParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &CreateSiteParams{
		ForceReload: &forceReloadDefault,

		Context: ctx,
	}
}

// NewCreateSiteParamsWithHTTPClient creates a new CreateSiteParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateSiteParamsWithHTTPClient(client *http.Client) *CreateSiteParams {
	var (
		forceReloadDefault = bool(false)
	)
	return &CreateSiteParams{
		ForceReload: &forceReloadDefault,
		HTTPClient:  client,
	}
}

/*CreateSiteParams contains all the parameters to send to the API endpoint
for the create site operation typically these are written to a http.Request
*/
type CreateSiteParams struct {

	/*Data*/
	Data *models.Site
	/*ForceReload
	  If set, do a force reload, do not wait for the configured reload-delay. Cannot be used when transaction is specified, as changes in transaction are not applied directly to configuration.

	*/
	ForceReload *bool
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

// WithTimeout adds the timeout to the create site params
func (o *CreateSiteParams) WithTimeout(timeout time.Duration) *CreateSiteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create site params
func (o *CreateSiteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create site params
func (o *CreateSiteParams) WithContext(ctx context.Context) *CreateSiteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create site params
func (o *CreateSiteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create site params
func (o *CreateSiteParams) WithHTTPClient(client *http.Client) *CreateSiteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create site params
func (o *CreateSiteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithData adds the data to the create site params
func (o *CreateSiteParams) WithData(data *models.Site) *CreateSiteParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the create site params
func (o *CreateSiteParams) SetData(data *models.Site) {
	o.Data = data
}

// WithForceReload adds the forceReload to the create site params
func (o *CreateSiteParams) WithForceReload(forceReload *bool) *CreateSiteParams {
	o.SetForceReload(forceReload)
	return o
}

// SetForceReload adds the forceReload to the create site params
func (o *CreateSiteParams) SetForceReload(forceReload *bool) {
	o.ForceReload = forceReload
}

// WithTransactionID adds the transactionID to the create site params
func (o *CreateSiteParams) WithTransactionID(transactionID *string) *CreateSiteParams {
	o.SetTransactionID(transactionID)
	return o
}

// SetTransactionID adds the transactionId to the create site params
func (o *CreateSiteParams) SetTransactionID(transactionID *string) {
	o.TransactionID = transactionID
}

// WithVersion adds the version to the create site params
func (o *CreateSiteParams) WithVersion(version *int64) *CreateSiteParams {
	o.SetVersion(version)
	return o
}

// SetVersion adds the version to the create site params
func (o *CreateSiteParams) SetVersion(version *int64) {
	o.Version = version
}

// WriteToRequest writes these params to a swagger request
func (o *CreateSiteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
