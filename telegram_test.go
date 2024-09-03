package telegram_test

import (
	"fmt"
	"testing"

	"github.com/nuominmin/telegram"
	"github.com/nuominmin/telegram/middleware"
	"gopkg.in/telebot.v3"
)

var bot *telegram.Bot
var pageAInlineKeyboard [][]telebot.InlineButton
var pageBInlineKeyboard [][]telebot.InlineButton

func TestNewTelebot(t *testing.T) {
	var err error
	bot, err = telegram.NewBot(telegram.DevToken, telegram.WithProxy("http://127.0.0.1:7890"))
	if err != nil {
		t.Error(err)
		return
	}

	bot.Use(middleware.Logger)

	err = bot.NewCommands().
		//AddCommand("/start", "show menu", startHandler).
		AddCommand("/start", "show menu", func(ctx telebot.Context) error {
			fmt.Println("========", ctx.Get("key"))
			return start(ctx)
		}, func(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
			return func(ctx telebot.Context) error {
				ctx.Set("key", "xxxxxxxx")
				return handlerFunc(ctx)
			}
		}).
		Commit()

	if err != nil {
		t.Error(err)
		return
	}

	pageAInlineKeyboard, err = bot.NewInlineKeyboard().
		NewRow().
		AddBackBtn("返回上一页").
		NewRow().
		AddReplyBtnWithData("pageACCCCCC", "pageACCCCCC", "pageACCCCCC", pageA).
		AddReplyBtnWithData("pageADDDDDD", "pageADDDDDD", "pageADDDDDD", pageA).
		Commit()
	if err != nil {
		t.Error(err)
		return
	}

	pageBInlineKeyboard, err = bot.NewInlineKeyboard().
		NewRow().
		AddBackBtn("返回上一页").
		NewRow().
		AddReplyBtnWithData("pageBXXXXXX", "pageBXXXXXX", "pageBXXXXXX", pageB).
		AddReplyBtnWithData("pageBWWWWWW", "pageBWWWWWW", "pageBWWWWWW", pageB).
		Commit()

	if err != nil {
		t.Error(err)
		return
	}

	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		fmt.Println("========", ctx.Get("key"))
		return nil
	})

	t.Logf("me: %+v", bot.Me())

	bot.Start()
}

func start(ctx telebot.Context) error {
	return ctx.Send(ctx.Text(), &telebot.ReplyMarkup{
		InlineKeyboard: pageAInlineKeyboard,
	}, telebot.ModeHTML)
}

func pageA(ctx telebot.Context) error {
	return ctx.Send(ctx.Text(), &telebot.ReplyMarkup{
		InlineKeyboard: pageBInlineKeyboard,
	}, telebot.ModeHTML)
}

func pageB(ctx telebot.Context) error {
	return ctx.Send(ctx.Callback().Data)
}

func replyCallbackData(ctx telebot.Context) error {
	data := ctx.Callback().Data
	if data == "" {
		data = "callback data is empty"
	}
	return ctx.Send(data)
}

func getRefreshDataFun(data string) telegram.CallbackDataFunc {
	return func(ctx telebot.Context) (string, error) {
		return fmt.Sprintf("cmd: %s, chat id: %d, username: %s", data, ctx.Chat().ID, ctx.Chat().Username), nil
	}
}
