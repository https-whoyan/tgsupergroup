package types

import "strings"

// Four Types: ParseModeNone, ParseModeHTML, ParseModeMarkdown, ParseModeMarkdownV2
// Use zero if not need to use escape text
type ParseMode uint8

const (
	ParseModeNone ParseMode = iota
	ParseModeHTML
	ParseModeMarkdown
	ParseModeMarkdownV2
)

// Code from https://github.com/go-telegram-bot-api/telegram-bot-api/blob/master/bot.go (Line 729)
func (m ParseMode) EscapeText(message string) string {
	var replacer *strings.Replacer
	switch m {
	case ParseModeNone:
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
