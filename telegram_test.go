package telegram_test

import (
	"gopkg.in/telebot.v3"
	"telegram"
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

	err = bot.NewCommands().
		AddCommand("/start", "show menu", startHandler).
		Commit()

	if err != nil {
		t.Error(err)
		return
	}

	startInlineKeyboard = bot.NewInlineKeyboard().
		NewRow().
		AddReplyBtn("reply", "reply", replyCallbackData).
		AddReplyBtnWithData("reply", "replywithdata", "datadata", replyCallbackData).
		AddReplyBtnWithData("withdrawsolana", "Withdraw", "withdraw", replyCallbackData).
		AddReplyBtnWithData("MenuRefresh", "Refresh", "Refresh data", replyCallbackData).
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

func replyCallbackData(ctx telebot.Context) error {
	data := ctx.Callback().Data
	if data == "" {
		data = "callback data is empty"
	}
	return ctx.Send(data)
}
