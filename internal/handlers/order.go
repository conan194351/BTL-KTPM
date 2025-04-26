package handlers

import (
	"encoding/json"
	"github.com/conan194351/BTL-KTPM/internal/dto"
	"github.com/conan194351/BTL-KTPM/internal/dto/response"
	"github.com/conan194351/BTL-KTPM/internal/errs"
	"github.com/conan194351/BTL-KTPM/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandlerImpl struct {
	orderSer *services.OrderServiceImpl
}

func NewOrderHandler(orderSer *services.OrderServiceImpl) *OrderHandlerImpl {
	return &OrderHandlerImpl{
		orderSer: orderSer,
	}
}

func (ah *OrderHandlerImpl) Create(ctx *gin.Context) {
	loginReq := dto.OderRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&loginReq); err != nil {
		response.SetHttpStatusError(ctx, errs.ErrMalformedJSON, err)
		return
	}
	res, err := ah.orderSer.CreateOrder(ctx, &loginReq)
	if err != nil {
		response.SetHttpStatusError(ctx, errs.ErrInternalServer, err)
		return
	}
	response.SetHttpStatusOK(ctx, http.StatusOK, "", res)
}
