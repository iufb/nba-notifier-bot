package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	store Storage
}

func NewBot(store *PostgresStore) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return nil, err
	}
	return &Bot{
		api:   bot,
		store: store,
	}, nil
}

func (b *Bot) Init(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		answerMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg := update.Message

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "register":
			m, err := b.RegisterNewAccount(ctx, msg)
			if err != nil {
				answerMsg.Text = fmt.Sprintf("Error occured while register new user, %+v", err)
			} else {
				answerMsg.Text = m
			}
		case "addTeam":
			answerMsg.Text = fmt.Sprintf("Team %s added.", answerMsg.Text)
		case "deleteTeam":
			answerMsg.Text = fmt.Sprintf("Team %s deleted.", answerMsg.Text)
		default:
			answerMsg.Text = "I don't know that command"
		}

		if _, err := b.api.Send(answerMsg); err != nil {
			log.Panic(err)
		}
	}
}

func (b *Bot) RegisterNewAccount(ctx context.Context, msg *tgbotapi.Message) (string, error) {
	account := &Account{TelegramId: int(msg.From.ID)}
	err := b.store.CreateAccount(ctx, account)
	if err != nil {
		return "", err
	}

	return "New account registered!", nil
}

func (b *Bot) DeleteAccount(ctx context.Context, msg *tgbotapi.Message) error {
	return nil
}
