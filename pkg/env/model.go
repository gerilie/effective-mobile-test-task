package env

type key string

const EnvKey key = "env"

const (
	Dev  = "dev"
	Prod = "prod"
	Host = "host"
)
