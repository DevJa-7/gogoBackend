package config

// Constants for API
const (
	HostURL    = "http://localhost"
	APIRoot    = "/gogo/api"
	APIVersion = "v1"
	APIURL     = APIRoot + "/" + APIVersion
)

// OneSignal parameters
const (
	OneSignalCreateNotificationURL = "https://onesignal.com/api/v1/notifications"
	UserAppName                    = "GoGo"
	UserBundleID                   = "gogo.user.com"
	UserAppID                      = "e822fe23-134d-47ac-9ed8-a04db1c3c6ba"
	UserAPIKey                     = "Y2Y5ODI0ZDctYTFjMS00ZTQ1LThkODYtNDE4NGE5N2E0Nzgw"
	DriverAppName                  = "GoGo Driver"
	DriverBundleID                 = "gogo.driver.com"
	DriverAppID                    = "9efccfe7-249e-4f74-9e62-accc8f58f038"
	DriverAPIKey                   = "YzUwNjY3MTYtZWMxNy00OGYxLTkxNjctOGQ4MWQwZWZiY2Q4"
)

// Centrifugo parameters
const (
	CentrifugoURL = "ws://18.216.19.36:9000/connection/websocket"
)
