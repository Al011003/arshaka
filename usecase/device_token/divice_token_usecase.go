package devicetoken

import (
	"backend/model"
	"backend/repo"
)

type SaveDeviceTokenUsecase interface {
    Save(userID uint, deviceToken string, deviceType string) error
}

type saveDeviceTokenUC struct {
    repo repo.DeviceTokenRepo
}

func NewSaveDeviceTokenUC(r repo.DeviceTokenRepo) SaveDeviceTokenUsecase {
    return &saveDeviceTokenUC{
        repo: r,
    }
}

func (uc *saveDeviceTokenUC) Save(userID uint, deviceToken, deviceType string) error {

    data := &model.DeviceToken{
        UserID:      userID,
        DeviceToken: deviceToken,
        DeviceType:  deviceType,
    }

    return uc.repo.SaveOrUpdate(data)
}
