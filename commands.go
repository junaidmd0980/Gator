package main

import (
	"config"
	"fmt"
)

type state struct {
	configPtr *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	CommandName map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.CommandName[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandName[name] = f
}
