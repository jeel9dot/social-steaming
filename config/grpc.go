package config

type GrpcConfig struct {
	Port string `envconfig:"GRPC_PORT"`
}