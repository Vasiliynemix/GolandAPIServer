package response

import (
	"just_for_fun/internal/server/router/structures"
)

// Response model info
// @Description Response
type Response struct {
	Status string `json:"status" example:"OK" binding:"required"` // OK or ERROR
}

// ErrorResponse model info
// @Description Error Response
type ErrorResponse struct {
	Status string `json:"status" example:"ERROR" binding:"required"`        // Always ERROR
	Error  string `json:"error" example:"Error message" binding:"required"` // Error message
}

// UserResponse model info
// @Description User Response
type UserResponse struct {
	Response                     // Response info status OK
	User     structures.UserShow `json:"data,omitempty" binding:"required"` // User info, if status is OK else nil
}

func OK() Response {
	return Response{
		Status: "OK",
	}
}

func Error(msgErr string) ErrorResponse {
	return ErrorResponse{
		Status: "ERROR",
		Error:  msgErr,
	}
}
