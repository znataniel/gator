package commands

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("error: wrong number of arguments provided")
	}

	err := s.Cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	println("user", cmd.Args[0], "has logged in")
	return nil
}
