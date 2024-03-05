package main

import (
	"time"
)

type AddToFavourite struct {
	Abbr string `json:"abbr"`
}
type Account struct {
	TelegramId int `json:"id" `
}

type Team struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

type Schedule struct {
	date time.Time
	ot   Team
}
