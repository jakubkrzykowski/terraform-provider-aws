// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cloudhsmv2

import (
	"github.com/aws/aws-sdk-go/private/protocol"
)

const (

	// ErrCodeCloudHsmAccessDeniedException for service response error code
	// "CloudHsmAccessDeniedException".
	//
	// The request was rejected because the requester does not have permission to
	// perform the requested operation.
	ErrCodeCloudHsmAccessDeniedException = "CloudHsmAccessDeniedException"

	// ErrCodeCloudHsmInternalFailureException for service response error code
	// "CloudHsmInternalFailureException".
	//
	// The request was rejected because of an AWS CloudHSM internal failure. The
	// request can be retried.
	ErrCodeCloudHsmInternalFailureException = "CloudHsmInternalFailureException"

	// ErrCodeCloudHsmInvalidRequestException for service response error code
	// "CloudHsmInvalidRequestException".
	//
	// The request was rejected because it is not a valid request.
	ErrCodeCloudHsmInvalidRequestException = "CloudHsmInvalidRequestException"

	// ErrCodeCloudHsmResourceNotFoundException for service response error code
	// "CloudHsmResourceNotFoundException".
	//
	// The request was rejected because it refers to a resource that cannot be found.
	ErrCodeCloudHsmResourceNotFoundException = "CloudHsmResourceNotFoundException"

	// ErrCodeCloudHsmServiceException for service response error code
	// "CloudHsmServiceException".
	//
	// The request was rejected because an error occurred.
	ErrCodeCloudHsmServiceException = "CloudHsmServiceException"

	// ErrCodeCloudHsmTagException for service response error code
	// "CloudHsmTagException".
	//
	// The request was rejected because of a tagging failure. Verify the tag conditions
	// in all applicable policies, and then retry the request.
	ErrCodeCloudHsmTagException = "CloudHsmTagException"
)

var exceptionFromCode = map[string]func(protocol.ResponseMetadata) error{
	"CloudHsmAccessDeniedException":     newErrorCloudHsmAccessDeniedException,
	"CloudHsmInternalFailureException":  newErrorCloudHsmInternalFailureException,
	"CloudHsmInvalidRequestException":   newErrorCloudHsmInvalidRequestException,
	"CloudHsmResourceNotFoundException": newErrorCloudHsmResourceNotFoundException,
	"CloudHsmServiceException":          newErrorCloudHsmServiceException,
	"CloudHsmTagException":              newErrorCloudHsmTagException,
}
