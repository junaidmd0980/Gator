package main

import (
	"config"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/junaidmd0980/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := state{
		db:     dbQueries,
		cfgPtr: &cfg,
	}

	cmds := commands{
		RegCmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough args provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
