// Code generated by go-swagger; DO NOT EDIT.

package eks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new eks API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for eks API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	ListEKSClusters(params *ListEKSClustersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSClustersOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  ListEKSClusters Lists EKS clusters
*/
func (a *Client) ListEKSClusters(params *ListEKSClustersParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListEKSClustersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListEKSClustersParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listEKSClusters",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/eks/clusters",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListEKSClustersReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListEKSClustersOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListEKSClustersDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
