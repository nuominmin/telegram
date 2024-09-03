package telegram

import "gopkg.in/telebot.v3"

type Commands struct {
	bot      *Bot
	commands []command
}
type command struct {
	endpoint    string
	description string
	h           telebot.HandlerFunc
	m           []telebot.MiddlewareFunc
	resetTrace  bool
}

func (b *Bot) NewCommands() *Commands {
	return &Commands{
		bot:      b,
		commands: make([]command, 0),
	}
}

func (c *Commands) AddCommand(text, description string, handler telebot.HandlerFunc, middlewares ...telebot.MiddlewareFunc) *Commands {
	c.commands = append(c.commands, command{text, description, handler, middlewares, false})
	return c
}

func (c *Commands) AddResetCommand(text, description string, handler telebot.HandlerFunc, middlewares ...telebot.MiddlewareFunc) *Commands {
	c.commands = append(c.commands, command{text, description, handler, middlewares, true})
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
		endpoint := c.commands[i].endpoint
		handler := func(ctx telebot.Context) error {

			return c.commands[i].h(ctx)

			//userId := ctx.Sender().ID
			//step, err := c.commands[i].h(ctx)
			//if err != nil {
			//	return err
			//}

			//if c.commands[i].resetTrace {
			//	c.bot.trace.Reset(userId, endpoint)
			//	return nil
			//}
			//c.bot.trace.Add(userId, endpoint, step)
			//return nil
		}

		err := c.bot.Handle(endpoint, handler, c.commands[i].m...)
		if err != nil {
			return err
		}
	}

	return c.bot.bot.SetCommands(opts)
}
