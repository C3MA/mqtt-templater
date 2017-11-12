package main

import (
	"fmt"
	"github.com/yosssi/gmq/mqtt/client"
	"crypto/tls"
	"math/rand"
)

type Message struct {
	Key string
	Value []byte
}

func buildClientID(prefix string) string {
	var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	hash := make([]byte, 7)
	for i := range hash {
		hash[i] = letters[rand.Intn(len(letters))]
	}
	return prefix + "-" + string(hash)
}

func connectToBroker(addr string, useTLS bool) *client.Client {
	tlsConfig := tls.Config{}

	mqttOptions := client.ConnectOptions{
		Network: "tcp",
		Address: addr,
		KeepAlive: 30,
		CONNACKTimeout: 15,
		PINGRESPTimeout: 15,
		ClientID: []byte(buildClientID("zlr-responder")),
	}
	
	if useTLS {
		mqttOptions.TLSConfig = &tlsConfig
	}

	clientOptions := client.Options{
		ErrorHandler: func(err error) {
			panic(fmt.Errorf("MQTT Error: %v", err))
		},
	}
	
	// Create new MQTT client
	mqtt := client.New(&clientOptions)
	
	// Connect to MQTT server
	err := mqtt.Connect(&mqttOptions)
	if err != nil {
		panic(fmt.Errorf("MQTT Connect() Error: %v", err))
	}			
	
	return mqtt
}

func receiveMessages(mqtt *client.Client, messages chan Message) {
	err := mqtt.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("/#"),
				Handler: func(topicName, message []byte) {
					messages <- Message{Key: string(topicName), Value: message}
				},
			},
		},
	})

	if err != nil {
		panic(fmt.Errorf("MQTT Subscribe() Error: %v", err))
	}
}