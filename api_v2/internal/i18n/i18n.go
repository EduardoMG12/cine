package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Localizer struct {
	messages         map[string]map[string]interface{}
	defaultLocale    string
	supportedLocales []string
}

type Message struct {
	Key  string
	Data map[string]interface{}
}

type contextKey string

const localeKey contextKey = "locale"

func NewLocalizer() (*Localizer, error) {
	l := &Localizer{
		messages:         make(map[string]map[string]interface{}),
		defaultLocale:    "en",
		supportedLocales: []string{"en", "pt", "es"},
	}

	for _, locale := range l.supportedLocales {
		if err := l.loadLocale(locale); err != nil {
			return nil, fmt.Errorf("failed to load locale %s: %w", locale, err)
		}
	}

	return l, nil
}

func (l *Localizer) loadLocale(locale string) error {
	filename := filepath.Join("internal", "i18n", "locales", fmt.Sprintf("%s.json", locale))

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read locale file %s: %w", filename, err)
	}

	var messages map[string]interface{}
	if err := json.Unmarshal(data, &messages); err != nil {
		return fmt.Errorf("failed to parse locale file %s: %w", filename, err)
	}

	l.messages[locale] = messages
	return nil
}

func (l *Localizer) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			locale := l.detectLocale(r)
			ctx := context.WithValue(r.Context(), localeKey, locale)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (l *Localizer) detectLocale(r *http.Request) string {

	if locale := r.URL.Query().Get("lang"); locale != "" {
		if l.isSupported(locale) {
			return locale
		}
	}

	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {

		languages := strings.Split(acceptLang, ",")
		for _, lang := range languages {

			lang = strings.Split(strings.TrimSpace(lang), ";")[0]

			mainLang := strings.Split(lang, "-")[0]

			if l.isSupported(mainLang) {
				return mainLang
			}
		}
	}

	return l.defaultLocale
}

func (l *Localizer) isSupported(locale string) bool {
	for _, supported := range l.supportedLocales {
		if supported == locale {
			return true
		}
	}
	return false
}

func (l *Localizer) Localize(ctx context.Context, key string, data map[string]interface{}) string {
	locale := l.getLocaleFromContext(ctx)
	return l.LocalizeWithLocale(locale, key, data)
}

func (l *Localizer) LocalizeWithLocale(locale, key string, data map[string]interface{}) string {

	if message := l.getMessage(locale, key); message != "" {
		return l.processTemplate(message, data)
	}

	if locale != l.defaultLocale {
		if message := l.getMessage(l.defaultLocale, key); message != "" {
			return l.processTemplate(message, data)
		}
	}

	return key
}

func (l *Localizer) getMessage(locale, key string) string {
	messages, exists := l.messages[locale]
	if !exists {
		return ""
	}

	parts := strings.Split(key, ".")
	current := messages

	for i, part := range parts {
		switch v := current[part].(type) {
		case string:
			if i == len(parts)-1 {
				return v
			}
			return ""
		case map[string]interface{}:
			if i == len(parts)-1 {
				return ""
			}
			current = v
		default:
			return ""
		}
	}

	return ""
}

func (l *Localizer) processTemplate(message string, data map[string]interface{}) string {
	if data == nil {
		return message
	}

	tmpl, err := template.New("message").Parse(message)
	if err != nil {
		return message
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return message
	}

	return result.String()
}

func (l *Localizer) getLocaleFromContext(ctx context.Context) string {
	if locale, ok := ctx.Value(localeKey).(string); ok {
		return locale
	}
	return l.defaultLocale
}

func (l *Localizer) GetSupportedLocales() []string {
	return l.supportedLocales
}

func (l *Localizer) T(ctx context.Context, key string, data ...map[string]interface{}) string {
	var templateData map[string]interface{}
	if len(data) > 0 {
		templateData = data[0]
	}
	return l.Localize(ctx, key, templateData)
}
