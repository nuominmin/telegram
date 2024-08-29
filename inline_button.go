package telegram

import (
	"gopkg.in/telebot.v3"
)

type InlineKeyboard struct {
	bot           *telebot.Bot
	inlineButtons [][]InlineButton
}

type InlineButton struct {
	telebot.InlineButton
	telebot.ReplyButton
	Handler     telebot.HandlerFunc
	Middlewares []telebot.MiddlewareFunc
}

func (b *Bot) NewInlineKeyboard() *InlineKeyboard {
	return &InlineKeyboard{
		bot:           b.bot,
		inlineButtons: make([][]InlineButton, 0),
	}
}

// NewRow 开始一个新的按钮行
func (ik *InlineKeyboard) NewRow() *InlineKeyboard {
	ik.inlineButtons = append(ik.inlineButtons, []InlineButton{})
	return ik
}

// AddButton 在当前行添加按钮
func (ik *InlineKeyboard) AddButton(button InlineButton) *InlineKeyboard {
	if len(ik.inlineButtons) == 0 {
		// 如果没有当前行，则创建一行
		ik.NewRow()
	}

	// 追加按钮到当前行（即最后一个切片）
	currentRowIdx := len(ik.inlineButtons) - 1
	ik.inlineButtons[currentRowIdx] = append(ik.inlineButtons[currentRowIdx], button)
	return ik
}

// AddReplyBtn 在当前行添加回复按钮
func (ik *InlineKeyboard) AddReplyBtn(unique, text string, handler telebot.HandlerFunc) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: telebot.InlineButton{
			Unique: unique,
			Text:   text,
		},
		Handler: handler,
	})
}

// AddReplyBtnWithData 在当前行添加回复按钮与数据
func (ik *InlineKeyboard) AddReplyBtnWithData(unique, text, data string, handler telebot.HandlerFunc) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: telebot.InlineButton{
			Unique: unique,
			Text:   text,
			Data:   data,
		},
		Handler: handler,
	})
}

// AddWebAppBtn 在当前行添加小程序按钮
func (ik *InlineKeyboard) AddWebAppBtn(text, webAppURL string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: telebot.InlineButton{
			Text:   text,
			WebApp: &telebot.WebApp{URL: webAppURL},
		},
	})
}

// AddUrlBtn 在当前行添加Url按钮
func (ik *InlineKeyboard) AddUrlBtn(text, url string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: telebot.InlineButton{
			Text: text,
			URL:  url,
		},
	})
}

// AddInlineQueryBtn 在当前行添加内联搜索按钮
func (ik *InlineKeyboard) AddInlineQueryBtn(text, query string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: telebot.InlineButton{
			Text:        text,
			InlineQuery: query,
		},
	})
}

func (ik *InlineKeyboard) Commit() (inlineKeyboard [][]telebot.InlineButton) {
	inlineKeyboard = make([][]telebot.InlineButton, 0)

	for i := 0; i < len(ik.inlineButtons); i++ {
		inlineButtons := make([]telebot.InlineButton, 0)

		for j := 0; j < len(ik.inlineButtons[i]); j++ {
			if ik.inlineButtons[i][j].Unique != "" && ik.inlineButtons[i][j].Handler != nil {
				ik.bot.Handle(&ik.inlineButtons[i][j].InlineButton, ik.inlineButtons[i][j].Handler, ik.inlineButtons[i][j].Middlewares...)
			}
			inlineButtons = append(inlineButtons, ik.inlineButtons[i][j].InlineButton)
		}
		inlineKeyboard = append(inlineKeyboard, inlineButtons)
	}

	return inlineKeyboard
}
