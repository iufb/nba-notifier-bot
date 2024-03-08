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
	team string
	date time.Time
	ot   string
}

var teamsList = []*Team{
	{Name: "Atlanta Hawks", Abbr: "ATL"},
	{Name: "Boston Celtics", Abbr: "BOS"},
	{Name: "Brooklyn Nets", Abbr: "BKN"},
	{Name: "Charlotte Hornets", Abbr: "CHA"},
	{Name: "Chicago Bulls", Abbr: "CHI"},
	{Name: "Cleveland Cavaliers", Abbr: "CLE"},
	{Name: "Dallas Mavericks", Abbr: "DAL"},
	{Name: "Denver Nuggets", Abbr: "DEN"},
	{Name: "Detroit Pistons", Abbr: "DET"},
	{Name: "Golden State Warriors", Abbr: "GSW"},
	{Name: "Houston Rockets", Abbr: "HOU"},
	{Name: "Indiana Pacers", Abbr: "IND"},
	{Name: "LA Clippers", Abbr: "LAC"},
	{Name: "Los Angeles Lakers", Abbr: "LAL"},
	{Name: "Memphis Grizzlies", Abbr: "MEM"},
	{Name: "Miami Heat", Abbr: "MIA"},
	{Name: "Milwaukee Bucks", Abbr: "MIL"},
	{Name: "Minnesota Timberwolves", Abbr: "MIN"},
	{Name: "New Orleans Pelicans", Abbr: "NOP"},
	{Name: "New York Knicks", Abbr: "NYK"},
	{Name: "Oklahoma City Thunder", Abbr: "OKC"},
	{Name: "Orlando Magic", Abbr: "ORL"},
	{Name: "Philadelphia 76ers", Abbr: "PHI"},
	{Name: "Phoenix Suns", Abbr: "PHX"},
	{Name: "Portland Trail Blazers", Abbr: "POR"},
	{Name: "Sacramento Kings", Abbr: "SAC"},
	{Name: "San Antonio Spurs", Abbr: "SAS"},
	{Name: "Toronto Raptors", Abbr: "TOR"},
	{Name: "Utah Jazz", Abbr: "UTA"},
	{Name: "Washington Wizards", Abbr: "WAS"},
}
