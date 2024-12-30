package config

type NatsConfig struct {
	Url string `envconfig:"NATS_URL"`
}