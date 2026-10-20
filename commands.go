package main

import (
	"config"
	"fmt"

	"github.com/junaidmd0980/gator/internal/database"
)

type state struct {
	db     *database.Queries
	cfgPtr *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	RegCmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.RegCmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.RegCmds[name] = f
}
