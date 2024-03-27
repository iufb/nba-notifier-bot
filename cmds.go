package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func includesStr(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

func GetTimeWithTimezone(s *ScheduleSc, offset int, team string) string {
	var res string
	withTimezone := s.date.Add(time.Hour * time.Duration(offset))
	var minutes string
	if withTimezone.Minute() == 0 {
		minutes = fmt.Sprintln("00")
	} else {
		minutes = fmt.Sprintln(s.date.Minute())
	}
	res = fmt.Sprint("*", team, "*", "\n", withTimezone.Month(), withTimezone.Day(), ", ", withTimezone.Hour(), ":", minutes, "\n", "*", s.ot, "*", "\n\n")
	return res
}

func CreateTeamsKeyboard(cmd string, teamsList []Team) tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup
	const TEAMS_PER_ROW = 5
	for i := 0; i < len(teamsList); i += TEAMS_PER_ROW {
		end := i + TEAMS_PER_ROW
		if end > len(teamsList) {
			end = len(teamsList)
		}
		var row []tgbotapi.InlineKeyboardButton
		for _, item := range teamsList[i:end] {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(item.Abbr, fmt.Sprintln(item.Abbr, cmd)))
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard
}

func RegisterNewAccount(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	account := &Account{TelegramId: int(update.Message.From.ID)}
	err := b.store.CreateAccount(ctx, account)
	if err != nil {
		return err
	}
	_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Registered successfully."))
	return err
}

func DeleteAccount(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	err := b.store.DeleteAccount(ctx, int(update.Message.From.ID))
	if err != nil {
		return err
	}
	_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Account deleted successfully."))
	return err
}

func AddTeamToFavourite(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose team:")
		msg.ReplyMarkup = CreateTeamsKeyboard("addTeam", teamsList)
		_, err := b.api.Send(msg)
		if err != nil {
			return err
		}
	} else {
		data := strings.Split(update.CallbackQuery.Data, " ")[0]
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, data)
		if _, err := b.api.Request(callback); err != nil {
			return err
		}
		err := b.store.AddTeamToFavourite(ctx, int(update.CallbackQuery.From.ID), data)
		if err != nil {
			log.Println(err)
			return err
		}

		// And finally, send a message containing the data received.
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintln("Team", data, "added."))
		if _, err := b.api.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func DeleteTeamFromFavourite(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose team to delete:")
		teams, err := b.store.GetAccountFavouriteTeams(ctx, int(update.Message.From.ID))
		if err != nil {
			return err
		}
		if len(teams) == 0 {
			return fmt.Errorf("No favourite teams found.")
		}
		msg.ReplyMarkup = CreateTeamsKeyboard("deleteTeam", teams)
		_, err = b.api.Send(msg)
		if err != nil {
			return err
		}
	} else {
		data := strings.Split(update.CallbackQuery.Data, " ")[0]
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, data)
		if _, err := b.api.Request(callback); err != nil {
			return err
		}
		err := b.store.DeleteTeamFromFavourite(ctx, int(update.CallbackQuery.From.ID), data)
		if err != nil {
			log.Println(err)
			return err
		}

		// And finally, send a message containing the data received.
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintln("Team", data, "deleted."))
		if _, err := b.api.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func GetNearestGame(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose team:")
		msg.ReplyMarkup = CreateTeamsKeyboard("ng", teamsList)
		_, err := b.api.Send(msg)
		if err != nil {
			return err
		}
	} else {
		data := strings.Split(update.CallbackQuery.Data, " ")[0]
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, data)
		if _, err := b.api.Request(callback); err != nil {
			return err
		}
		var schedule string
		teamUrl, err := b.store.GetTeamUrl(ctx, data)
		if err != nil {
			return err
		}
		s, err := Scrapper(teamUrl)
		if err != nil {
			return err
		}
		schedule = GetTimeWithTimezone(&s, 9, data)

		// And finally, send a message containing the data received.
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, schedule)
		msg.ParseMode = "MarkdownV2"
		if _, err := b.api.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func SendSchedule(ctx context.Context, b *Bot, telegramId int) error {
	teams, err := b.store.GetAccountFavouriteTeams(ctx, telegramId)
	log.Println("EXEC")
	var schedule string
	for _, team := range teams {
		s, err := Scrapper(team.Url)
		if err != nil {
			log.Println(err)
			continue
		}

		withTimezone := GetTimeWithTimezone(&s, 9, team.Abbr)
		if time.Now().Day()+1 == s.date.Add(time.Hour*time.Duration(9)).Day() {
			schedule += withTimezone
		}
	}
	if len(schedule) == 0 {
		schedule = "No matches for your favourite teams tomorrow."
	}
	msg := tgbotapi.NewMessage(int64(telegramId), schedule)
	msg.ParseMode = "MarkdownV2"
	_, err = b.api.Send(msg)
	return err
}
