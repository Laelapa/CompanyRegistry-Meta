package logging

import "strings"

// sanitizeLogValue replaces values that could be used for log injection or are otherwise problematic.
func sanitizeLogValue(v string) string {
	replacer := strings.NewReplacer(
		"\n", "[LF]", // Line feed
		"\r", "[CR]", // Carriage return
		"\u0000", "[NUL]", // Null byte
		"\u001b", "[ESC]", // ANSI Escape
		"\u200B", "[ZWS]", // Zero width space
		"\u2028", "[LS]", // Line separator
		"\u2029", "[PS]", // Paragraph separator
		"\u2063", "[IS]", // Invisible separator
		// JSON structural characters handled by zap
	)
	return replacer.Replace(v)
}

func truncateLogValue(v string, maxLen int) string {
	vRune := []rune(v)
	if len(vRune) <= maxLen {
		return v
	}
	return string(vRune[:maxLen]) + "[TRUNCATED]"
}

// FiletLogValue sanitizes and truncates a log value to a maximum length.
// It should be used on any data that could be controlled by an external user.
func (l *Logger) FiletLogValue(v string) string {
	return truncateLogValue(sanitizeLogValue(v), l.maxHeaderLength)
}
