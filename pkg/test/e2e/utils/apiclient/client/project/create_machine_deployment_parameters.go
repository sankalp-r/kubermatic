// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// NewCreateMachineDeploymentParams creates a new CreateMachineDeploymentParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateMachineDeploymentParams() *CreateMachineDeploymentParams {
	return &CreateMachineDeploymentParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateMachineDeploymentParamsWithTimeout creates a new CreateMachineDeploymentParams object
// with the ability to set a timeout on a request.
func NewCreateMachineDeploymentParamsWithTimeout(timeout time.Duration) *CreateMachineDeploymentParams {
	return &CreateMachineDeploymentParams{
		timeout: timeout,
	}
}

// NewCreateMachineDeploymentParamsWithContext creates a new CreateMachineDeploymentParams object
// with the ability to set a context for a request.
func NewCreateMachineDeploymentParamsWithContext(ctx context.Context) *CreateMachineDeploymentParams {
	return &CreateMachineDeploymentParams{
		Context: ctx,
	}
}

// NewCreateMachineDeploymentParamsWithHTTPClient creates a new CreateMachineDeploymentParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateMachineDeploymentParamsWithHTTPClient(client *http.Client) *CreateMachineDeploymentParams {
	return &CreateMachineDeploymentParams{
		HTTPClient: client,
	}
}

/* CreateMachineDeploymentParams contains all the parameters to send to the API endpoint
   for the create machine deployment operation.

   Typically these are written to a http.Request.
*/
type CreateMachineDeploymentParams struct {

	// Body.
	Body *models.NodeDeployment

	// ClusterID.
	ClusterID string

	// ProjectID.
	ProjectID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create machine deployment params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateMachineDeploymentParams) WithDefaults() *CreateMachineDeploymentParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create machine deployment params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateMachineDeploymentParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithTimeout(timeout time.Duration) *CreateMachineDeploymentParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithContext(ctx context.Context) *CreateMachineDeploymentParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithHTTPClient(client *http.Client) *CreateMachineDeploymentParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithBody(body *models.NodeDeployment) *CreateMachineDeploymentParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetBody(body *models.NodeDeployment) {
	o.Body = body
}

// WithClusterID adds the clusterID to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithClusterID(clusterID string) *CreateMachineDeploymentParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithProjectID adds the projectID to the create machine deployment params
func (o *CreateMachineDeploymentParams) WithProjectID(projectID string) *CreateMachineDeploymentParams {
	o.SetProjectID(projectID)
	return o
}

// SetProjectID adds the projectId to the create machine deployment params
func (o *CreateMachineDeploymentParams) SetProjectID(projectID string) {
	o.ProjectID = projectID
}

// WriteToRequest writes these params to a swagger request
func (o *CreateMachineDeploymentParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param cluster_id
	if err := r.SetPathParam("cluster_id", o.ClusterID); err != nil {
		return err
	}

	// path param project_id
	if err := r.SetPathParam("project_id", o.ProjectID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
