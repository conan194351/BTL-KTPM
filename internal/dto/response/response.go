package response

import (
	"github.com/conan194351/BTL-KTPM/internal/errs"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func DefaultResponse(ctx *gin.Context, res Response) {
	ctx.JSON(res.Status, res)
}

func SetHttpStatusError(ctx *gin.Context, e errs.ErrorCode, er interface{}) {
	err := errs.ResponseError(e)
	res := Response{
		Status:  err.HttpStatusCode,
		Message: err.Message,
		Error:   er,
	}
	DefaultResponse(ctx, res)
}

func SetHttpStatusOK(ctx *gin.Context, status int, message string, data interface{}) {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	DefaultResponse(ctx, res)
}
