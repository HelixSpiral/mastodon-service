package main

type MqttMessage struct {
	MastodonServer       string
	MastodonClientID     string
	MastodonClientSecret string
	MastodonUser         string
	MastodonPass         string

	Message string
	Image   string
}
