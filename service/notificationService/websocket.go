package notificationService

type subEventHandler struct{}

/*
func (h *subEventHandler) OnMessage(sub *centrifuge.Sub, msg *centrifuge.Message) {
	log.Println(fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), msg))
}

func (h *subEventHandler) OnJoin(sub *centrifuge.Sub, msg *centrifuge.ClientInfo) {
	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", msg.User, msg.Client, sub.Channel()))
}

func (h *subEventHandler) OnLeave(sub *centrifuge.Sub, msg *centrifuge.ClientInfo) {
	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", msg.User, msg.Client, sub.Channel()))
}

// In production you need to receive credentials from application backend.
func credentials(id string) *centrifuge.Credentials {
	// Never show secret to client of your application. Keep it on your application backend only.
	secret := "secret"
	// Application user ID.
	// id := "58f87135cfc9c544988ba1a8"
	// Current timestamp as string.
	timestamp := centrifuge.Timestamp()
	// Empty info.
	info := ""
	// Generate client token so Centrifugo server can trust connection parameters received from client.
	token := "" //auth.GenerateClientToken(secret, id, timestamp, info)

	return &centrifuge.Credentials{
		User:      id,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}
}
//*/

// PushWebsocketNotification publish message to business id
func PushWebsocketNotification(channel string, message interface{}) {
	// creds := credentials(channel)

	// wsURL := config.CentrifugoURL
	// c := centrifuge.New(wsURL, creds, nil, centrifuge.DefaultConfig())
	// defer c.Close()

	// err := c.Connect()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// events := centrifuge.NewSubEventHandler()
	// subEventHandler := &subEventHandler{}
	// events.OnMessage(subEventHandler)
	// events.OnJoin(subEventHandler)
	// events.OnLeave(subEventHandler)

	// sub, err := c.Subscribe("public:"+channel, events)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// dataBytes, _ := json.Marshal(message)
	// err = sub.Publish(dataBytes)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("ws sync publish is successful!")
	// }

	// err = sub.Unsubscribe()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
