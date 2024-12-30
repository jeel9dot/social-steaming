package nats

import (
	nc "github.com/nats-io/nats.go"
)

type NatsClient struct {
	conn *nc.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := nc.Connect(nc.DefaultURL)
	if err != nil {
		return nil, err
	}
	return &NatsClient{conn}, nil
}

func (c *NatsClient) Publish(subject, message string) error {
	return c.conn.Publish(subject, []byte(message))
}

func (c *NatsClient) Subscribe(subject string, handler func(msg *nc.Msg)) (*nc.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

func (c *NatsClient) Close() {
	c.conn.Close()
}