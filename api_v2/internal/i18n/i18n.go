package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed locales/*.json
var localeFS embed.FS

type Localizer struct {
	messages map[string]map[string]string
}

var globalLocalizer *Localizer

func init() {
	globalLocalizer = &Localizer{
		messages: make(map[string]map[string]string),
	}
	loadMessages()
}

func loadMessages() {
	languages := []string{"en", "pt", "es"}

	for _, lang := range languages {
		filename := fmt.Sprintf("locales/%s.json", lang)
		data, err := localeFS.ReadFile(filename)
		if err != nil {
			continue
		}

		var messages map[string]string
		if err := json.Unmarshal(data, &messages); err != nil {
			continue
		}

		globalLocalizer.messages[lang] = messages
	}
}

func T(key string, lang ...string) string {
	language := "en" // default language
	if len(lang) > 0 && lang[0] != "" {
		language = strings.ToLower(lang[0][:2])
	}

	if messages, ok := globalLocalizer.messages[language]; ok {
		if msg, exists := messages[key]; exists {
			return msg
		}
	}

	// Fallback to English
	if language != "en" {
		if messages, ok := globalLocalizer.messages["en"]; ok {
			if msg, exists := messages[key]; exists {
				return msg
			}
		}
	}

	// If no translation found, return the key
	return key
}

func Tf(key string, args []interface{}, lang ...string) string {
	template := T(key, lang...)
	return fmt.Sprintf(template, args...)
}

func GetSupportedLanguages() []string {
	return []string{"en", "pt", "es"}
}

func IsLanguageSupported(lang string) bool {
	if len(lang) < 2 {
		return false
	}
	lang = strings.ToLower(lang[:2])
	_, exists := globalLocalizer.messages[lang]
	return exists
}
