package main

type MqttMessage struct {
	MastodonClientID     string
	MastodonClientSecret string
	MastodonUser         string
	MastodonPass         string

	Message string
	Image   string
}
