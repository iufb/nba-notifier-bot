package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
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
	go b.SendNotification(ctx)
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

func (b *Bot) SendNotification(ctx context.Context) {
	targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 22, 0, 0, 0, time.Local)
	accounts, err := b.store.GetAccounts(ctx)
	if err != nil {
		log.Println("No registered account found.")
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			currentTime := time.Now()
			if currentTime.After(targetTime) {
				targetTime = targetTime.Add(24 * time.Hour)
				for _, acc := range accounts {
					err := SendSchedule(ctx, b, acc.TelegramId)
					if err != nil {
						log.Println(err)
					}

				}
			}
			time.Sleep(time.Hour)
		}
	}
}

func (b *Bot) RegisterNewCommand(name string, cmd CmdViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]CmdViewFunc)
	}
	b.cmdViews[name] = cmd
}

func (b *Bot) ExecCommand(ctx context.Context, update tgbotapi.Update, cmdName string) {
	cmdView, ok := b.cmdViews[cmdName]
	if !ok {
		log.Println("Command ", cmdName, "not found.")
		return
	}
	if err := cmdView(ctx, b, update); err != nil {
		log.Printf("[ERROR] failed to execute view: %v", err)

		if _, err := b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintln("[ERROR]:", err))); err != nil {
			log.Printf("[ERROR] failed to send error message: %v", err)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()
	if update.Message != nil {
		cmd := update.Message.Command()
		b.ExecCommand(ctx, update, cmd)

	} else if update.CallbackQuery != nil {
		cmd := strings.TrimSpace(strings.Split(update.CallbackQuery.Data, " ")[1])
		b.ExecCommand(ctx, update, cmd)
	}
}
