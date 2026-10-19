package main

import (
	"config"
	"fmt"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}
	cfg.SetUser("Junaid")
	cfg, err = config.Read()
	fmt.Println("DB URL:", cfg.DbUrl)
	fmt.Println("Username:", cfg.CurrentUserName)
}
