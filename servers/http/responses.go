// Package http provides HTTP utilities for building APIs,
// including status codes, custom error codes, and response formatting.
package http

import "github.com/gofiber/fiber/v2"

// StatusCode defines HTTP status codes used for server responses.
type StatusCode int

const (
	OK                  StatusCode = 200 // HTTP 200 OK
	Redirect            StatusCode = 301 // HTTP 301 Moved Permanently
	BadRequest          StatusCode = 400 // HTTP 400 Bad Request
	Unauthorized        StatusCode = 401 // HTTP 401 Unauthorized
	Forbidden           StatusCode = 403 // HTTP 403 Forbidden
	NotFound            StatusCode = 404 // HTTP 404 Not Found
	InternalServerError StatusCode = 500 // HTTP 500 Internal Server Error
)

// CustomCode defines application-specific custom error codes.
type CustomCode int

const (
	CodeSuccess CustomCode = 0 // Indicates successful operation.

	CodeClientError     CustomCode = 100 // Indicates a client error.
	CodeValidationError CustomCode = 101 // Indicates a validation error in the request.
	CodeAuthError       CustomCode = 102 // Indicates an authentication error.
	CodeNotFoundError   CustomCode = 103 // Indicates a resource not found error.
	CodeUnauthorized    CustomCode = 104 // Indicates an unauthorized error.
	CodeForbidden       CustomCode = 105 // Indicates a forbidden error.

	CodeInternalError CustomCode = 200 // Indicates an internal server error.
)

// CustomCodeToHTTPStatus maps a CustomCode to its corresponding HTTP status code.
func CustomCodeToHTTPStatus(code CustomCode) StatusCode {
	switch code {
	case CodeSuccess:
		return OK
	case CodeClientError, CodeValidationError, CodeAuthError:
		return BadRequest
	case CodeInternalError:
		return InternalServerError
	case CodeUnauthorized:
		return Unauthorized
	case CodeForbidden:
		return Forbidden
	case CodeNotFoundError:
		return NotFound
	default:
		return InternalServerError
	}
}

// NewResponse creates a standardized JSON response for the API.
//
// Arguments:
//
//	ctx - The Fiber context used to send the response.
//	status - The HTTP status code to return (e.g., 200, 404).
//	result - The data payload to include in the response.
//	errCode - A custom application-specific error code.
//	errMsg - A custom error message (ignored if errCode is CodeSuccess).
//
// Behavior:
//   - If `errCode` is not `CodeSuccess` and `errMsg` is empty, a default error message is used.
//   - If `errMsg` is provided, `result` is set to nil and `errMsg` is included as the message.
//   - Returns a JSON response with the structure:
//     {
//     "code": <CustomCode>,
//     "message": <string>,
//     "data": <any>
//     }
//
// Returns:
//
//	An error if the response cannot be sent via the Fiber context.
func NewResponse(ctx *fiber.Ctx, status StatusCode, result any, errCode CustomCode, errMsg string) error {
	msg := "success"

	if errCode != CodeSuccess && errMsg == "" {
		errMsg = "An unexpected error occurred"
	}

	if errMsg != "" {
		msg = errMsg
		result = nil
	}

	response := map[string]any{
		"code":    errCode,
		"message": msg,
		"data":    result,
	}

	return ctx.Status(int(status)).JSON(response)
}
