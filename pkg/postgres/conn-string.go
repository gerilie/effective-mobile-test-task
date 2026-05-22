package postgres

import (
	"fmt"
	"net"
)

// BuildConnString constructs a PostgreSQL connection string from the configuration.
//
// It formats the connection parameters into a valid PostgreSQL DSN
// using the standard postgres:// protocol format.
func BuildConnString(cfg Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		net.JoinHostPort(cfg.Host, cfg.Port),
		cfg.DB,
	)
}
