package notificationService

import "testing"

func TestPushWebsocketNotification(t *testing.T) {
	data := map[string]interface{}{
		"type":    "orderRequest",
		"orderId": "order1234",
	}
	PushWebsocketNotification("abcd1234", data)
}
