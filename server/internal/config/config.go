package config

import (
	"server/internal/server"
	"server/internal/presence"
	"server/internal/store"
)

type Environment string

const (
	Dev  = "dev"
	Prod = "prod"
)

func (e Environment) IsDev() bool {
	return e == Dev
}

type Config struct {
	Environment  Environment
	Store        *store.Store
	Presence     *presence.RedisPresence
	LocalClients *server.LocalClients
}
