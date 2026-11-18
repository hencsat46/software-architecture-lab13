package config

import "net/url"

type PostgresConfig struct{ Host, Port, User, Password, DBName string }

func (p PostgresConfig) DSN() string {
	u := url.URL{Scheme: "postgres", Host: p.Host + ":" + p.Port, User: url.UserPassword(p.User, p.Password), Path: p.DBName}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}

type Config struct {
	Port, LogLevel, UserURL, NotificationURL string
	Postgres                                 PostgresConfig
}

func Load() *Config {
	return &Config{Port: getString("SERVER_PORT", ":8080"), LogLevel: getString("LOGGER_LEVEL", "info"), UserURL: getString("USER_SERVICE_URL", "http://localhost:8081"), NotificationURL: getString("NOTIFICATION_SERVICE_URL", "http://localhost:8083"), Postgres: PostgresConfig{Host: getString("POSTGRES_HOST", "localhost"), Port: getString("POSTGRES_PORT", "5432"), User: getString("POSTGRES_USER", "orders"), Password: getString("POSTGRES_PASSWORD", "orders_local"), DBName: getString("POSTGRES_DB", "orders")}}
}
