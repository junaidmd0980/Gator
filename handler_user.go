package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("you must provide username")
	}
	err := s.configPtr.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User '%s' has been set..\n", cmd.Args[0])
	return nil
}
