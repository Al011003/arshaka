package devicetokenhandler

import (
	request "backend/dto/request/device"
	response "backend/dto/response/common"
	devicetoken "backend/usecase/device_token"

	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceTokenHandler struct {
    saveUC devicetoken.SaveDeviceTokenUsecase
}

func NewDeviceTokenHandler(saveUC devicetoken.SaveDeviceTokenUsecase) *DeviceTokenHandler {
    return &DeviceTokenHandler{
        saveUC: saveUC,
    }
}

func (h *DeviceTokenHandler) Save(c *gin.Context) {
    // Bind JSON request
    var req request.SaveDeviceTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.BaseResponse{
            Status:    "error",
            Message: "invalid request body",
        })
        return
    }

    // Get userID from JWT context (misal lu simpan di "user_id")
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, response.BaseResponse{
            Status:    "error",
            Message: "unauthorized",
        })
        return
    }

    userID := userIDValue.(uint)

    // Call usecase
    err := h.saveUC.Save(userID, req.DeviceToken, req.DeviceType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.BaseResponse{
            Status:    "error",
            Message: "failed to save device token",
        })
        return
    }

    // Success
    c.JSON(http.StatusOK, response.BaseResponse{
        Status:    "success",
        Message: "device token saved successfully",
        Data:    nil, // karena lu mau kosong
    })
}
