package main

import (
	"context"
	"fmt"
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

func addNewTeamKeyboard() tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup
	const TEAMS_PER_ROW = 5
	for i := 0; i < len(teamsList); i += TEAMS_PER_ROW {
		if teamsList[i].Abbr == "ATL" {
			continue
		}
		end := i + TEAMS_PER_ROW
		if end > len(teamsList) {
			end = len(teamsList)
		}
		var row []tgbotapi.InlineKeyboardButton
		for _, item := range teamsList[i:end] {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(item.Abbr, item.Abbr))
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
	err := b.store.DeleteAccount(int(update.Message.From.ID))
	if err != nil {
		return err
	}
	_, err = b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Account deleted successfully."))
	return err
}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func AddTeamToFavourite(ctx context.Context, b *Bot, update tgbotapi.Update) error {
	receivedMsg := update.Message.Text
	cmd := strings.Split(receivedMsg, " ")
	fmt.Println(cmd)
	if len(cmd) != 2 {
		b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid command type : example - /addTeam GSW"))
		return fmt.Errorf("Invalid command %s :", receivedMsg)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln("Team", cmd[1], "added"))
	_, err := b.api.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
