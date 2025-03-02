// Code generated by go-swagger; DO NOT EDIT.

package aks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListAKSClustersReader is a Reader for the ListAKSClusters structure.
type ListAKSClustersReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAKSClustersReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAKSClustersOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAKSClustersDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAKSClustersOK creates a ListAKSClustersOK with default headers values
func NewListAKSClustersOK() *ListAKSClustersOK {
	return &ListAKSClustersOK{}
}

/* ListAKSClustersOK describes a response with status code 200, with default header values.

AKSClusterList
*/
type ListAKSClustersOK struct {
	Payload models.AKSClusterList
}

func (o *ListAKSClustersOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/clusters][%d] listAKSClustersOK  %+v", 200, o.Payload)
}
func (o *ListAKSClustersOK) GetPayload() models.AKSClusterList {
	return o.Payload
}

func (o *ListAKSClustersOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAKSClustersDefault creates a ListAKSClustersDefault with default headers values
func NewListAKSClustersDefault(code int) *ListAKSClustersDefault {
	return &ListAKSClustersDefault{
		_statusCode: code,
	}
}

/* ListAKSClustersDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListAKSClustersDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list a k s clusters default response
func (o *ListAKSClustersDefault) Code() int {
	return o._statusCode
}

func (o *ListAKSClustersDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aks/clusters][%d] listAKSClusters default  %+v", o._statusCode, o.Payload)
}
func (o *ListAKSClustersDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAKSClustersDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
