// Code generated by go-swagger; DO NOT EDIT.

package azure

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListAzureSizesNoCredentialsReader is a Reader for the ListAzureSizesNoCredentials structure.
type ListAzureSizesNoCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAzureSizesNoCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAzureSizesNoCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAzureSizesNoCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAzureSizesNoCredentialsOK creates a ListAzureSizesNoCredentialsOK with default headers values
func NewListAzureSizesNoCredentialsOK() *ListAzureSizesNoCredentialsOK {
	return &ListAzureSizesNoCredentialsOK{}
}

/* ListAzureSizesNoCredentialsOK describes a response with status code 200, with default header values.

AzureSizeList
*/
type ListAzureSizesNoCredentialsOK struct {
	Payload models.AzureSizeList
}

func (o *ListAzureSizesNoCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/azure/sizes][%d] listAzureSizesNoCredentialsOK  %+v", 200, o.Payload)
}
func (o *ListAzureSizesNoCredentialsOK) GetPayload() models.AzureSizeList {
	return o.Payload
}

func (o *ListAzureSizesNoCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAzureSizesNoCredentialsDefault creates a ListAzureSizesNoCredentialsDefault with default headers values
func NewListAzureSizesNoCredentialsDefault(code int) *ListAzureSizesNoCredentialsDefault {
	return &ListAzureSizesNoCredentialsDefault{
		_statusCode: code,
	}
}

/* ListAzureSizesNoCredentialsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListAzureSizesNoCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list azure sizes no credentials default response
func (o *ListAzureSizesNoCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *ListAzureSizesNoCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/azure/sizes][%d] listAzureSizesNoCredentials default  %+v", o._statusCode, o.Payload)
}
func (o *ListAzureSizesNoCredentialsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAzureSizesNoCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
