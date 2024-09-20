package telegram_test

import (
	"fmt"
	"github.com/nuominmin/telegram"
	"github.com/nuominmin/telegram/middleware"
	"gopkg.in/telebot.v3"
	"testing"
)

var bot *telegram.Bot
var startInlineKeyboard [][]telebot.InlineButton

func TestNewTelebot(t *testing.T) {
	var err error
	bot, err = telegram.NewBot(telegram.DevToken, telegram.WithProxy("http://127.0.0.1:7890"))
	if err != nil {
		t.Error(err)
		return
	}

	bot.Use(middleware.Logger)

	if err = bot.SetWebAppMenuButton("GO", "https://goethereumbook.org/zh/smart-contract-read/"); err != nil {
		t.Error(err)
		return
	}

	//err = bot.NewCommands().
	//	//AddCommand("/start", "show menu", startHandler).
	//	AddCommand("/start1", "show menu", func(ctx telebot.Context) error {
	//		fmt.Println("========", ctx.Get("key"))
	//		return startHandler(ctx)
	//	}, func(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
	//		return func(ctx telebot.Context) error {
	//			ctx.Set("key", "xxxxxxxx")
	//			return handlerFunc(ctx)
	//		}
	//	}).
	//	Commit()

	if err != nil {
		t.Error(err)
		return
	}

	startInlineKeyboard = bot.NewInlineKeyboard().
		NewRow().
		AddReplyBtn("buy", "buy", buy).
		AddReplyBtnWithData("reply", "replywithdata", "replywithdata", replyCallbackData).
		AddReplyBtnWithData("withdrawsolana", "Withdraw", "Withdraw", replyCallbackData).
		AddReplyBtnWithDataFunc("MenuRefresh", "Refresh", getDataFun("Refresh"), replyCallbackData).
		NewRow().
		AddWebAppBtn("智能合约文档", "https://goethereumbook.org/zh/smart-contract-read/").
		AddWebAppBtn("Google", "https://google.com").
		AddInlineQueryBtn("查询", "xixi").
		Commit()

	t.Logf("me: %+v", bot.Me())

	bot.Start()
	defer bot.Stop()
}

func startHandler(ctx telebot.Context) error {
	message := `Solana · 🅴 (https://solscan.io/account/93oFkxpYEB7yjmySq5Jsdn9y4BCZa4fK28u19teutP8S)
93oFkxpYEB7yjmySq5Jsdn9y4BCZa4fK28u19teutP8S  (Tap to copy)
Total Balance: $0.00
Sol Balance: 0.000 SOL ($0.00)

Press the Refresh button to update your current balance.

Join Telegram group @sillybot_users for help and questions about Sillybot`

	return ctx.Send(message, &telebot.ReplyMarkup{
		InlineKeyboard: startInlineKeyboard,
	}, telebot.ModeHTML)
}

func buy(ctx telebot.Context) error {
	return ctx.Send("Enter a token contract address to buy")
}

func replyCallbackData(ctx telebot.Context) error {
	data := ctx.Callback().Data
	if data == "" {
		data = "callback data is empty"
	}
	return ctx.Send(data)
}

func getDataFun(data string) telegram.CallbackDataFunc {
	return func(ctx telebot.Context) (string, error) {
		return fmt.Sprintf("cmd: %s, chat id: %d, username: %s", data, ctx.Chat().ID, ctx.Chat().Username), nil
	}
}
