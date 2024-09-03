package telegram_test

import (
	"fmt"
	"github.com/nuominmin/telegram"
	"github.com/nuominmin/telegram/middleware"
	"gopkg.in/telebot.v3"
	tele "gopkg.in/telebot.v3"
	"testing"
)

var tBot *telebot.Bot
var (
	// 主菜单按钮
	btnTrade    = tele.InlineButton{Unique: "trade", Text: "Trade"}
	btnBalance  = tele.InlineButton{Unique: "balance", Text: "View Balance"}
	btnSettings = tele.InlineButton{Unique: "settings", Text: "Settings"}

	// 买入页面的按钮
	btnBuy     = tele.InlineButton{Unique: "buy", Text: "Buy"}
	btnSell    = tele.InlineButton{Unique: "sell", Text: "Sell"}
	btnBack    = tele.InlineButton{Unique: "back", Text: "Back"}
	btnRefresh = tele.InlineButton{Unique: "refresh", Text: "Refresh"}

	// 买入细节页面按钮
	btn0_5SOL      = tele.InlineButton{Unique: "0_5sol", Text: "0.5 SOL"}
	btn1SOL        = tele.InlineButton{Unique: "1sol", Text: "1 SOL"}
	btnMaxSOL      = tele.InlineButton{Unique: "maxsol", Text: "Max SOL"}
	btnSetSlippage = tele.InlineButton{Unique: "slippage", Text: "Set Slippage"}
	btnConfirmBuy  = tele.InlineButton{Unique: "confirm_buy", Text: "BUY"}

	// 存储用户选择
	selectedAmount = ""
)

func TestNewTelebot1(t *testing.T) {
	var err error
	tBot, err = telegram.NewTelebot(telegram.DevToken, telegram.WithProxy("http://127.0.0.1:7890"))
	if err != nil {
		t.Error(err)
		return
	}

	tBot.Use(middleware.Logger)

	// 设置主菜单处理
	tBot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send("Welcome to Trojan Bot! Please select an option:", &telebot.ReplyMarkup{
			InlineKeyboard: buildMainMenu(),
		})
	})

	// 注册处理器
	setupHandlers(tBot)

	tBot.Start()
}

// 构建主菜单
func buildMainMenu() [][]tele.InlineButton {
	return [][]tele.InlineButton{
		{btnTrade},
		{btnBalance},
		{btnSettings},
	}
}

// 构建买入细节页面，根据用户的选择更新按钮文本
func buildBuyDetailsMenu() [][]tele.InlineButton {
	// 根据选择更新按钮文本
	if selectedAmount == "0.5" {
		btn0_5SOL.Text = "0.5 SOL ✔️"
		btn1SOL.Text = "1 SOL"
		btnMaxSOL.Text = "Max SOL"
	} else if selectedAmount == "1" {
		btn0_5SOL.Text = "0.5 SOL"
		btn1SOL.Text = "1 SOL ✔️"
		btnMaxSOL.Text = "Max SOL"
	} else if selectedAmount == "max" {
		btn0_5SOL.Text = "0.5 SOL"
		btn1SOL.Text = "1 SOL"
		btnMaxSOL.Text = "Max SOL ✔️"
	} else {
		// 默认状态
		btn0_5SOL.Text = "0.5 SOL"
		btn1SOL.Text = "1 SOL"
		btnMaxSOL.Text = "Max SOL"
	}

	return [][]tele.InlineButton{
		{btn0_5SOL, btn1SOL, btnMaxSOL},
		{btnSetSlippage},
		{btnConfirmBuy},
		{btnBack, btnRefresh},
	}
}

// 构建交易菜单
func buildTradeMenu() [][]tele.InlineButton {
	return [][]tele.InlineButton{
		{btnBuy, btnSell},
		{btnBack},
	}
}
func errHandlers(c tele.Context, err error) error {
	if err == nil {
		c.Respond(&tele.CallbackResponse{})
	} else {
		c.Respond(&tele.CallbackResponse{
			Text:      "An error occurred!",
			ShowAlert: true,
		})
	}
	return err
}

// 注册按钮交互处理函数
func setupHandlers(b *tele.Bot) {
	// 处理"Trade"按钮
	b.Handle(&btnTrade, func(c tele.Context) error {
		err := c.Send("Choose a trade action:", &tele.ReplyMarkup{
			InlineKeyboard: buildTradeMenu(),
		})
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})

	// 处理"Buy"按钮
	b.Handle(&btnBuy, func(c tele.Context) error {
		// 模拟显示用户余额、代币价格等信息
		err := c.Send(fmt.Sprintf("Balance: %.2f SOL\nPrice: $0.000002079\nSelect the amount of SOL to buy BONK tokens.", 3.0), &tele.ReplyMarkup{
			InlineKeyboard: buildBuyDetailsMenu(),
		})
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})

	// 处理买入页面上的 SOL 选择按钮
	b.Handle(&btn0_5SOL, func(c tele.Context) error {
		selectedAmount = "0.5" // 更新选择
		return updateSelection(c)
	})

	b.Handle(&btn1SOL, func(c tele.Context) error {
		selectedAmount = "1" // 更新选择
		return updateSelection(c)
	})

	b.Handle(&btnMaxSOL, func(c tele.Context) error {
		selectedAmount = "max" // 更新选择
		return updateSelection(c)
	})

	// 处理滑点设置
	b.Handle(&btnSetSlippage, func(c tele.Context) error {
		err := c.Send("Please input slippage percentage (e.g., 1 for 1%)")
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})

	// 处理确认买入按钮
	b.Handle(&btnConfirmBuy, func(c tele.Context) error {
		// 模拟买入操作
		err := c.Send("Buy order placed successfully.")
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})

	// 处理刷新按钮
	b.Handle(&btnRefresh, func(c tele.Context) error {
		err := c.Send("Refreshing data...", &tele.ReplyMarkup{
			InlineKeyboard: buildBuyDetailsMenu(),
		})
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})

	// 处理"Back"按钮，返回上一级菜单
	b.Handle(&btnBack, func(c tele.Context) error {
		err := c.Send("Returning to previous menu:", &tele.ReplyMarkup{
			InlineKeyboard: buildTradeMenu(),
		})
		if err == nil {
			c.Respond(&tele.CallbackResponse{})
		} else {
			c.Respond(&tele.CallbackResponse{
				Text:      "An error occurred!",
				ShowAlert: true,
			})
		}
		return err
	})
}

// 更新用户的选择，并重新渲染页面
func updateSelection(c tele.Context) error {
	// 更新按钮状态并重新显示按钮
	err := c.Edit("Select the amount of SOL to buy BONK tokens:", &tele.ReplyMarkup{
		InlineKeyboard: buildBuyDetailsMenu(),
	})
	if err == nil {
		c.Respond(&tele.CallbackResponse{})
	} else {
		c.Respond(&tele.CallbackResponse{
			Text:      "An error occurred!",
			ShowAlert: true,
		})
	}
	return err
}
