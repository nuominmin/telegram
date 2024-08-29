package telegram

import "gopkg.in/telebot.v3"

type Commands struct {
	bot      *telebot.Bot
	commands []command
}
type command struct {
	endpoint    string
	description string
	h           telebot.HandlerFunc
	m           []telebot.MiddlewareFunc
}

func (b *Bot) NewCommands() *Commands {
	return &Commands{
		bot:      b.bot,
		commands: make([]command, 0),
	}
}

func (c *Commands) AddCommand(text, description string, handler telebot.HandlerFunc, middlewares ...telebot.MiddlewareFunc) *Commands {
	c.commands = append(c.commands, command{text, description, handler, middlewares})
	return c
}

// Commit 提交
func (c *Commands) Commit() error {
	opts := make([]telebot.Command, len(c.commands))
	for i := 0; i < len(c.commands); i++ {
		opts[i] = telebot.Command{
			Text:        c.commands[i].endpoint,
			Description: c.commands[i].description,
		}
		c.bot.Handle(c.commands[i].endpoint, c.commands[i].h, c.commands[i].m...)
	}

	return c.bot.SetCommands(opts)
}
