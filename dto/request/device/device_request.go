package request

type SaveDeviceTokenRequest struct {
	DeviceToken string `json:"device_token" binding:"required"`
	DeviceType  string `json:"device_type" binding:"required"` // android / ios
}