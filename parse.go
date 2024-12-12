package tgsupergroup

import "strings"

type ParseMode uint8

const (
	ParseModeHTML ParseMode = iota + 1
	ParseModeMarkdown
	ParseModeMarkdownV2
)

// EscapeText
// Code from https://github.com/go-telegram-bot-api/telegram-bot-api/blob/master/bot.go (Line 729)
// Special thanks
func (m ParseMode) EscapeText(message string) string {
	var replacer *strings.Replacer
	switch m {
	case 0:
		return message
	case ParseModeHTML:
		replacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	case ParseModeMarkdown:
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	case ParseModeMarkdownV2:
		replacer = strings.NewReplacer(
			"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
			"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
			"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
			"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
		)
	default:
		replacer = strings.NewReplacer()
	}
	return replacer.Replace(message)
}
