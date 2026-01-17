package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client wraps the MQTT client functionality
type Client struct {
	client mqtt.Client
	broker string
}

// Config holds MQTT connection configuration
type Config struct {
	Broker   string // e.g., "tcp://localhost:1883"
	ClientID string
	Username string
	Password string
}

// NewClient creates a new MQTT client with the given configuration
func NewClient(config Config) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}
	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	// Set connection options
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetKeepAlive(30 * time.Second)

	// Set connection callbacks
	opts.OnConnect = func(c mqtt.Client) {
		log.Println("MQTT: Connected to broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Printf("MQTT: Connection lost: %v\n", err)
	}
	opts.OnReconnecting = func(c mqtt.Client, opts *mqtt.ClientOptions) {
		log.Println("MQTT: Reconnecting to broker...")
	}

	client := mqtt.NewClient(opts)

	return &Client{
		client: client,
		broker: config.Broker,
	}, nil
}

// Connect establishes connection to the MQTT broker
func (c *Client) Connect() error {
	log.Printf("MQTT: Connecting to broker at %s...\n", c.broker)

	token := c.client.Connect()
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("connection timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	return nil
}

// Disconnect closes the connection to the MQTT broker
func (c *Client) Disconnect() {
	log.Println("MQTT: Disconnecting from broker...")
	c.client.Disconnect(250)
}

// Publish publishes a message to the specified topic
func (c *Client) Publish(topic string, payload interface{}) error {
	if !c.client.IsConnected() {
		return fmt.Errorf("client not connected")
	}

	// Convert payload to string if it's not already a byte slice
	var data interface{}
	switch v := payload.(type) {
	case []byte:
		data = v
	case string:
		data = v
	default:
		data = fmt.Sprintf("%v", v)
	}

	token := c.client.Publish(topic, 0, false, data)
	if !token.WaitTimeout(2 * time.Second) {
		return fmt.Errorf("publish timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}

	log.Printf("MQTT: Published to topic '%s'\n", topic)
	return nil
}

// IsConnected returns whether the client is currently connected
func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}
