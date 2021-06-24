package notificationService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"../../config"
	"../../model"
)

// PushOneSignalNotification pushes notification to mobile device via OneSignal
func PushOneSignalNotification(notification *model.OneSignalNotification, apiKey string) {
	b, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("notification parsing is failed!")
		return
	}

	req, err := http.NewRequest("POST", config.OneSignalCreateNotificationURL, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Base "+apiKey)
	// fmt.Println(req.Header)
	// fmt.Println(bytes.NewBuffer(b))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client is not working:", err)
		return
	}

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	resp.Body.Close()
}
