package telegram

import "gopkg.in/telebot.v3"

const STEP = ""

func ContextAddStep(ctx telebot.Context, step string) {
	ctx.Set("step", step)
}

func ContextGetStep(ctx telebot.Context) string {
	if step := ctx.Get("step"); step != nil {
		return step.(string)
	}
	return ""
}
