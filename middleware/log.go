package middleware

import (
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

func Logger(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		var command string
		switch {
		case ctx.Callback() != nil:
			command = ctx.Callback().Unique

		case ctx.Message() != nil:
			command = ctx.Message().Text
		default:
			command = ctx.Text()
		}

		start := time.Now()
		err := handlerFunc(ctx)
		log.Printf("chat id: %d, username: %s, command: %s, took %v to complete", ctx.Chat().ID, ctx.Chat().Username, command, time.Since(start))
		return err
	}
}
