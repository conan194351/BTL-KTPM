package handlers

import (
	"encoding/json"
	"github.com/conan194351/BTL-KTPM/internal/dto"
	"github.com/conan194351/BTL-KTPM/internal/dto/response"
	"github.com/conan194351/BTL-KTPM/internal/errs"
	"github.com/conan194351/BTL-KTPM/internal/models"
	"github.com/conan194351/BTL-KTPM/internal/services"
	"github.com/conan194351/BTL-KTPM/pkg/utils"
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
	user := ctx.MustGet("user")
	userModel, err := utils.ConvertToStruct[models.User](user)
	if err != nil {
		response.SetHttpStatusError(ctx, errs.ErrInternalServer, err)
		return
	}
	orderReq := dto.OderRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&orderReq); err != nil {
		response.SetHttpStatusError(ctx, errs.ErrMalformedJSON, err)
		return
	}
	orderReq.UserID = userModel.ID
	res, err := ah.orderSer.CreateOrder(ctx, &orderReq)
	if err != nil {
		response.SetHttpStatusError(ctx, errs.ErrInternalServer, err.Error())
		return
	}
	response.SetHttpStatusOK(ctx, http.StatusOK, "", res)
}

func (ah *OrderHandlerImpl) Test(ctx *gin.Context) {
	err := ah.orderSer.Test(ctx, 2, 1)
	if err != nil {
		response.SetHttpStatusError(ctx, errs.ErrInternalServer, err)
		return
	}
	response.SetHttpStatusOK(ctx, http.StatusOK, "", nil)
}
