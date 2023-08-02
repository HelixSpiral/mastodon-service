package main

type MqttMessage struct {
	ClientID     string
	ClientSecret string
	MastodonUser string
	MastodonPass string

	Message string
	Image   string
}
