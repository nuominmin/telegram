package telegram

import (
	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
	"strings"
)

type CallbackDataFunc func(telebot.Context) (string, error)

type InlineKeyboard struct {
	bot           *Bot
	inlineButtons [][]InlineButton

	backUniques []string
}

type InlineButton struct {
	*telebot.InlineButton
	DataFunc    CallbackDataFunc
	Handler     telebot.HandlerFunc
	Middlewares []telebot.MiddlewareFunc
	resetTrace  bool
}

func (b *Bot) NewInlineKeyboard() *InlineKeyboard {
	return &InlineKeyboard{
		bot:           b,
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
		InlineButton: &telebot.InlineButton{
			Unique: unique,
			Text:   text,
		},
		Handler: handler,
	})
}

// AddReplyBtnWithData 在当前行添加回复按钮与数据
func (ik *InlineKeyboard) AddReplyBtnWithData(unique, text, data string, handler telebot.HandlerFunc) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Unique: unique,
			Text:   text,
			Data:   data,
		},
		Handler: handler,
	})
}

// AddReplyBtnWithDataFunc 在当前行添加回复按钮与数据方法
func (ik *InlineKeyboard) AddReplyBtnWithDataFunc(unique, text string, dataFunc CallbackDataFunc, handler telebot.HandlerFunc) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Unique: unique,
			Text:   text,
		},
		Handler:  handler,
		DataFunc: dataFunc,
	})
}

// AddWebAppBtn 在当前行添加小程序按钮
func (ik *InlineKeyboard) AddWebAppBtn(text, webAppURL string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Text:   text,
			WebApp: &telebot.WebApp{URL: webAppURL},
		},
	})
}

// AddUrlBtn 在当前行添加Url按钮
func (ik *InlineKeyboard) AddUrlBtn(text, url string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Text: text,
			URL:  url,
		},
	})
}

// AddInlineQueryBtn 在当前行添加内联搜索按钮
func (ik *InlineKeyboard) AddInlineQueryBtn(text, query string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Text:        text,
			InlineQuery: query,
		},
	})
}

// AddBackBtn 在当前行添加返回按钮
func (ik *InlineKeyboard) AddBackBtn(text string) *InlineKeyboard {
	return ik.AddButton(InlineButton{
		InlineButton: &telebot.InlineButton{
			Unique: ik.genBackUnique(),
			Text:   text,
		},
		Handler: func(ctx telebot.Context) error {
			userId := ctx.Sender().ID
			end := "/start"
			steps := ik.bot.trace.GetSteps(userId)
			if len(steps) < 2 {
				return ik.bot.Trigger(end, ctx)
			}

			end = steps[len(steps)-2]
			if idx := len(steps) - 3; idx >= 0 {
				ik.bot.trace.ResetSteps(userId, steps[:len(steps)-3])
			}
			return ik.bot.Trigger(end, ctx)
		},
	})
}

func (ik *InlineKeyboard) genBackUnique() string {
	backUnique := strings.ReplaceAll(uuid.New().String(), "-", "")
	ik.backUniques = append(ik.backUniques, backUnique)
	return backUnique
}

func (ik *InlineKeyboard) Commit() (inlineKeyboard [][]telebot.InlineButton, err error) {

	inlineKeyboard = make([][]telebot.InlineButton, 0)

	for i := 0; i < len(ik.inlineButtons); i++ {
		inlineButtons := make([]telebot.InlineButton, 0)

		for j := 0; j < len(ik.inlineButtons[i]); j++ {
			if ik.inlineButtons[i][j].Unique == "" || ik.inlineButtons[i][j].Handler == nil {
				continue
			}

			end := ik.inlineButtons[i][j].InlineButton.CallbackUnique()
			handler := func(ctx telebot.Context) error {
				if dataFunc := ik.inlineButtons[i][j].DataFunc; dataFunc != nil {
					data, e := dataFunc(ctx)
					if e != nil {
						return e
					}
					ctx.Callback().Data = data
				}

				if err = ik.inlineButtons[i][j].Handler(ctx); err != nil {
					return err
				}

				for n := 0; n < len(ik.backUniques); n++ {
					if ik.inlineButtons[i][j].Unique == ik.backUniques[n] {
						return nil
					}
				}

				ik.bot.trace.SaveStep(ctx.Sender().ID, end)
				ik.bot.trace.SaveData(ctx.Sender().ID, ContextGetStep(ctx))

				return nil
			}
			if err = ik.bot.Handle(end, handler, ik.inlineButtons[i][j].Middlewares...); err != nil {
				return nil, err
			}
			inlineButtons = append(inlineButtons, *ik.inlineButtons[i][j].InlineButton)
		}
		inlineKeyboard = append(inlineKeyboard, inlineButtons)
	}

	return inlineKeyboard, nil
}
