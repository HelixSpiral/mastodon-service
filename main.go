package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mattn/go-mastodon"
)

func main() {
	// MQTT config setup
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	mqttTopic := os.Getenv("MQTT_TOPIC")
	mqttUsername := os.Getenv("MQTT_USERNAME")
	mqttPassword := os.Getenv("MQTT_PASSWORD")

	// Setup the MQTT client options
	options := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientId)
	options.ConnectRetry = true
	options.AutoReconnect = true

	if mqttUsername != "" {
		options.SetUsername(mqttUsername)
		if mqttPassword != "" {
			options.SetPassword(mqttPassword)
		}
	}

	options.OnConnectionLost = func(c mqtt.Client, e error) {
		log.Println("Connection lost")
	}
	options.OnConnect = func(c mqtt.Client) {
		log.Println("Connected")

		t := c.Subscribe(mqttTopic, 2, nil)
		go func() {
			_ = t.Wait()
			if t.Error() != nil {
				log.Printf("Error subscribing: %s\n", t.Error())
			} else {
				log.Println("Subscribed to:", mqttTopic)
			}
		}()
	}
	options.OnReconnecting = func(_ mqtt.Client, co *mqtt.ClientOptions) {
		log.Println("Attempting to reconnect")
	}
	options.DefaultPublishHandler = func(_ mqtt.Client, m mqtt.Message) {
		// Unmarshal the received json into a struct
		var mqttMsg MqttMessage
		err := json.Unmarshal(m.Payload(), &mqttMsg)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received: %s->%s\n", m.Topic(), mqttMsg.Message)

		// If we get a message with no Mastodon info, ignore it
		if mqttMsg.MastodonClientID == "" {
			return
		}

		// Default to botsin.space if not provided.
		if mqttMsg.MastodonServer == "" {
			mqttMsg.MastodonServer = "https://botsin.space"
		}

		// Create mastodon client with provided credentials
		c := mastodon.NewClient(&mastodon.Config{
			Server:       mqttMsg.MastodonServer,
			ClientID:     mqttMsg.MastodonClientID,
			ClientSecret: mqttMsg.MastodonClientSecret,
		})

		// Authenticate
		err = c.Authenticate(context.Background(), mqttMsg.MastodonUser, mqttMsg.MastodonPass)
		if err != nil {
			log.Fatal(err)
		}

		// Define the base toot
		mastodonToot := &mastodon.Toot{
			Status: mqttMsg.Message,
		}

		// If we've been given an image, upload it and attach the media id to our toot
		if len(mqttMsg.Image) > 0 {
			mediaAttachment, err := c.UploadMediaFromBytes(context.Background(), mqttMsg.Image)
			if err != nil {
				log.Fatal(err)
			}

			mastodonToot.MediaIDs = append(mastodonToot.MediaIDs, mediaAttachment.ID)
		}

		status, err := c.PostStatus(context.Background(), mastodonToot)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(status)
	}

	// Setup the MQTT client with the options we set
	mqttClient := mqtt.NewClient(options)

	// Connect to the MQTT server
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Connected")

	// Block indefinitely until something above errors, or we close out.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig

	log.Println("Signal caught -> Exit")
	mqttClient.Disconnect(1000)
}
