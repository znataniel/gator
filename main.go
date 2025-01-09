package main

import (
	"fmt"
	"os"

	"github.com/znataniel/gator/internal/commands"
	"github.com/znataniel/gator/internal/config"
)

func printConfig(cfg config.Config) {
	println("db url:\t", cfg.DbUrl)
	println("current user:\t", cfg.CurrentUser)
}

func main() {
	configs, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	state := commands.State{
		Cfg: &configs,
	}

	c := commands.Commands{
		Comms: make(map[string]func(*commands.State, commands.Command) error),
	}
	c.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("error: no command provided")
		os.Exit(1)
	}

	err = c.Run(&state, commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)

}
