package notificationService

import (
	"testing"

	"../../config"
	"../../model"
)

func TestUserNotification(t *testing.T) {
	notification := &model.OneSignalNotification{}
	notification.AppID = config.UserAppID
	notification.PlayerIds = []string{"f5a6e8a0-1f9a-4207-844d-4b631f8f9117"}
	notification.Title = map[string]interface{}{"en": config.UserAppName}
	notification.Message = map[string]interface{}{"en": "Your order is accepted."}
	notification.Data = map[string]interface{}{
		"type":    "OrderAccept",
		"orderId": "5a2fd232a118745d685fd28f",
	}

	PushOneSignalNotification(notification, config.UserAPIKey)
}

func TestDriverNotification(t *testing.T) {
	notification := &model.OneSignalNotification{}
	notification.AppID = config.DriverAppID
	notification.PlayerIds = []string{"4b974b8f-fbcf-4f58-9b46-4ad22481fab4"}
	notification.Title = map[string]interface{}{"en": config.DriverAppName}
	notification.Message = map[string]interface{}{"en": "Your order trip is accepted."}
	notification.Data = map[string]interface{}{
		"type":    "TripAccept",
		"orderId": "5a2fd232a118745d685fd28f",
	}

	PushOneSignalNotification(notification, config.DriverAPIKey)
}
