// Code generated by go-swagger; DO NOT EDIT.

package regions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/utils/apiclient/models"
)

// ListAWSRegionsReader is a Reader for the ListAWSRegions structure.
type ListAWSRegionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAWSRegionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAWSRegionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewListAWSRegionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewListAWSRegionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewListAWSRegionsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAWSRegionsOK creates a ListAWSRegionsOK with default headers values
func NewListAWSRegionsOK() *ListAWSRegionsOK {
	return &ListAWSRegionsOK{}
}

/* ListAWSRegionsOK describes a response with status code 200, with default header values.

Regions
*/
type ListAWSRegionsOK struct {
	Payload []models.Regions
}

func (o *ListAWSRegionsOK) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aws/regions][%d] listAWSRegionsOK  %+v", 200, o.Payload)
}
func (o *ListAWSRegionsOK) GetPayload() []models.Regions {
	return o.Payload
}

func (o *ListAWSRegionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAWSRegionsUnauthorized creates a ListAWSRegionsUnauthorized with default headers values
func NewListAWSRegionsUnauthorized() *ListAWSRegionsUnauthorized {
	return &ListAWSRegionsUnauthorized{}
}

/* ListAWSRegionsUnauthorized describes a response with status code 401, with default header values.

EmptyResponse is a empty response
*/
type ListAWSRegionsUnauthorized struct {
}

func (o *ListAWSRegionsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aws/regions][%d] listAWSRegionsUnauthorized ", 401)
}

func (o *ListAWSRegionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListAWSRegionsForbidden creates a ListAWSRegionsForbidden with default headers values
func NewListAWSRegionsForbidden() *ListAWSRegionsForbidden {
	return &ListAWSRegionsForbidden{}
}

/* ListAWSRegionsForbidden describes a response with status code 403, with default header values.

EmptyResponse is a empty response
*/
type ListAWSRegionsForbidden struct {
}

func (o *ListAWSRegionsForbidden) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aws/regions][%d] listAWSRegionsForbidden ", 403)
}

func (o *ListAWSRegionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewListAWSRegionsDefault creates a ListAWSRegionsDefault with default headers values
func NewListAWSRegionsDefault(code int) *ListAWSRegionsDefault {
	return &ListAWSRegionsDefault{
		_statusCode: code,
	}
}

/* ListAWSRegionsDefault describes a response with status code -1, with default header values.

errorResponse
*/
type ListAWSRegionsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list a w s regions default response
func (o *ListAWSRegionsDefault) Code() int {
	return o._statusCode
}

func (o *ListAWSRegionsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v2/providers/aws/regions][%d] listAWSRegions default  %+v", o._statusCode, o.Payload)
}
func (o *ListAWSRegionsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAWSRegionsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
