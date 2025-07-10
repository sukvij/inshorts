package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	Version     string `json:"version"`
	LatencyMs   int64  `json:"latencyMs"`
	Environment string `json:"environment,omitempty"`
	// Pagination  *Pagination `json:"pagination,omitempty"`
	// RateLimit   *RateLimit  `json:"rateLimit,omitempty"`
}

// type Pagination struct {
// 	Total      int64 `json:"total"`
// 	Page       int   `json:"page"`
// 	PerPage    int   `json:"perPage"`
// 	TotalPages int   `json:"totalPages"`
// }

// type RateLimit struct {
// 	Limit     int   `json:"limit"`
// 	Remaining int   `json:"remaining"`
// 	Reset     int64 `json:"reset"`
// }

type FinalResponse struct {
	Success    bool      `json:"success"`
	StatusCode int       `json:"statusCode"`
	Data       any       `json:"data,omitempty"`
	Error      *AppError `json:"error,omitempty"`
	Meta       *Meta     `json:"meta"`
}

func JSONResponse(ctx *gin.Context, err error, data interface{}, totalTime int64) {
	response := &FinalResponse{Data: data, Meta: &Meta{LatencyMs: totalTime}}
	if err == nil {
		response.Success = true
		response.StatusCode = 200
		ctx.JSON(response.StatusCode, response)
		return
	}
	ReturnErrorWithCode(ctx, err, response)
	if response.StatusCode == 204 {
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(response.StatusCode, response)
}
