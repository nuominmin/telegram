package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Bot struct {
	bot *telebot.Bot
}

func NewBot(token string, opts ...Option) (*Bot, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	options := newOptions(opts...)
	if options.proxyUrl != "" {
		uri, err := url.Parse(options.proxyUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL. error: %v", err)
		}

		client := http.DefaultClient
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}

		settings.Client = client
	}

	if options.poller > 0 {
		settings.Poller = &telebot.LongPoller{Timeout: options.poller}
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot: bot,
	}, err
}

func (b *Bot) Start() {
	log.Println("Tg server started. ")
	b.bot.Start()
}

func (b *Bot) Me() *telebot.User {
	return b.bot.Me
}

func (b *Bot) Use(middleware ...telebot.MiddlewareFunc) {
	b.bot.Use(middleware...)
}

func (b *Bot) Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
	b.bot.Handle(endpoint, h, m...)
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

func (b *Bot) SetMenuButtonWebAppWithUser(chatId int64, text, webAppUrl string) error {
	return b.bot.SetMenuButton(&telebot.User{ID: chatId}, &telebot.MenuButton{
		Type: telebot.MenuButtonWebApp,
		Text: text,
		WebApp: &telebot.WebApp{
			URL: webAppUrl,
		},
	})
}

func (b *Bot) SetMenuButtonWebApp(text, webAppUrl string) error {
	return b.bot.SetMenuButton(nil, &telebot.MenuButton{
		Type: telebot.MenuButtonWebApp,
		Text: text,
		WebApp: &telebot.WebApp{
			URL: webAppUrl,
		},
	})
}

func (b *Bot) Bot() *telebot.Bot {
	return b.bot
}
