package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "broker_address"
	port := 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername("conf42")
	opts.SetPassword("Something")

	opts.SetDefaultPublishHandler(messaheHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub(client)
	publish(client)
	client.Disconnect(250)
}

var messaheHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	fmt.Printf("Received message %s from topic %s ", m.Payload(), m.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(c mqtt.Client) {
	options := c.OptionsReader()
	fmt.Printf("Connected to: %v", options.Servers())
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
	fmt.Printf("Error occurred %v", err)
}

func sub(client mqtt.Client) {
	topic := "conf42/#"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subsribed to topic %s", topic)

}

func publish(client mqtt.Client) {
	num := 10
	topic := ""
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}
