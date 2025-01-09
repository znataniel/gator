package commands

import (
	"fmt"

	"github.com/znataniel/gator/internal/config"
)

type State struct {
	Cfg *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Comms map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if _, exists := c.Comms[name]; exists {
		println("command", name, "already exists")
		return
	}

	c.Comms[name] = f
	return
}

func (c *Commands) Run(s *State, cmd Command) error {
	if _, exists := c.Comms[cmd.Name]; !exists {
		return fmt.Errorf("error: command provided does not exist")
	}
	return c.Comms[cmd.Name](s, cmd)
}
