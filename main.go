package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	store, err := NewPostgresStore()
	if err != nil {
		log.Printf("[ERROR] failed to connect to db: %v", err)
		return
	}
	defer store.db.Close()
	bot, err := NewBot(store)
	if err != nil {
		log.Printf("[ERROR] failed to create botAPI: %v", err)
		return
	}
	bot.RegisterNewCommand("register", RegisterNewAccount)
	bot.RegisterNewCommand("delete", DeleteAccount)
	bot.RegisterNewCommand("addTeam", AddTeamToFavourite)
	bot.RegisterNewCommand("deleteTeam", DeleteTeamFromFavourite)
	bot.RegisterNewCommand("ng", GetNearestGame)
	// bot.api.Debug = true
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
