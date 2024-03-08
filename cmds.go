package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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

func addNewTeamKeyboard(cmd string, teamsList []*Team) tgbotapi.InlineKeyboardMarkup {
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
		msg.ReplyMarkup = addNewTeamKeyboard("addTeam", teamsList)
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
		msg.ReplyMarkup = addNewTeamKeyboard("deleteTeam", teams)
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

func SendSchedule(ctx context.Context, b *Bot, telegramId int) error {
	teams, err := b.store.GetAccountFavouriteTeams(ctx, telegramId)
	var schedule string
	for _, team := range teams {
		s, err := b.store.GetSchedule(ctx, team.Abbr)
		if err != nil {
			continue
		}
		schedule += fmt.Sprintln(s.team, s.date.Format("2006-01-02 15:04:05"), s.ot, "\n")
	}
	if len(schedule) == 0 {
		schedule = "No matches for your favourite teams tomorrow."
	}

	_, err = b.api.Send(tgbotapi.NewMessage(int64(telegramId), schedule))
	return err
}
