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
	bot   *telebot.Bot
	trace *Trace

	endpoints map[string]struct{}
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
		bot:       bot,
		trace:     NewTrace(),
		endpoints: make(map[string]struct{}),
	}, err
}

func NewTelebot(token string, opts ...Option) (*telebot.Bot, error) {
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

	return bot, err
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

func (b *Bot) Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) error {
	end := b.extractEndpoint(endpoint)
	if end == "" {
		return fmt.Errorf("telebot: unsupported endpoint")
	}
	if _, ok := b.endpoints[end]; ok {
		return fmt.Errorf("endpoint already exists, endpoint: %s", end)
	}
	b.endpoints[end] = struct{}{}
	b.bot.Handle(end, h, m...)
	return nil
}

func (b *Bot) extractEndpoint(endpoint interface{}) string {
	switch end := endpoint.(type) {
	case string:
		return end
	case telebot.CallbackEndpoint:
		return end.CallbackUnique()
	}
	return ""
}

func (b *Bot) Trigger(endpoint interface{}, ctx telebot.Context) error {
	return b.bot.Trigger(endpoint, ctx)
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
