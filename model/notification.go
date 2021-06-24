package model

// OneSignalNotification struct.
type OneSignalNotification struct {
	AppID     string                 `json:"app_id"`
	PlayerIds []string               `json:"include_player_ids"`
	Title     map[string]interface{} `json:"headings,omitempty"`
	Message   map[string]interface{} `json:"contents,omitempty"`
	Data      interface{}            `json:"data,omitempty"` // for both
}
