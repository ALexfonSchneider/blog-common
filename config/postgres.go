package config

import "fmt"

type Postgres struct {
	User     string `koanf:"user" validate:"required"`
	Password string `koanf:"password" validate:"required"`
	Database string `koanf:"database" validate:"required"`
	Host     string `koanf:"host" validate:"required"`
	Port     int    `koanf:"port" validate:"required"`
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Database,
	)
}
