package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	store    Storage
	cmdViews map[string]CmdViewFunc
}
type CmdViewFunc func(context.Context, *Bot, tgbotapi.Update) error

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

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, time.Second*5)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegisterNewCommand(name string, cmd CmdViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]CmdViewFunc)
	}
	b.cmdViews[name] = cmd
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()
	if update.Message == nil {
		return
	}
	if !update.Message.IsCommand() {
		return
	}
	cmd := update.Message.Command()
	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		return
	}
	if err := cmdView(ctx, b, update); err != nil {
		log.Printf("[ERROR] failed to execute view: %v", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Internal error")); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
}
