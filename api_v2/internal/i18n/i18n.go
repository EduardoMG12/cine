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

// Localizer handles internationalization
type Localizer struct {
	messages         map[string]map[string]interface{}
	defaultLocale    string
	supportedLocales []string
}

// Message represents a localized message with data
type Message struct {
	Key  string
	Data map[string]interface{}
}

// contextKey for storing locale in context
type contextKey string

const localeKey contextKey = "locale"

// NewLocalizer creates a new localizer instance
func NewLocalizer() (*Localizer, error) {
	l := &Localizer{
		messages:         make(map[string]map[string]interface{}),
		defaultLocale:    "en",
		supportedLocales: []string{"en", "pt", "es"},
	}

	// Load all locale files
	for _, locale := range l.supportedLocales {
		if err := l.loadLocale(locale); err != nil {
			return nil, fmt.Errorf("failed to load locale %s: %w", locale, err)
		}
	}

	return l, nil
}

// loadLocale loads messages for a specific locale
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

// Middleware returns an HTTP middleware that detects and sets locale
func (l *Localizer) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			locale := l.detectLocale(r)
			ctx := context.WithValue(r.Context(), localeKey, locale)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// detectLocale detects the locale from the request
func (l *Localizer) detectLocale(r *http.Request) string {
	// 1. Check query parameter
	if locale := r.URL.Query().Get("lang"); locale != "" {
		if l.isSupported(locale) {
			return locale
		}
	}

	// 2. Check Accept-Language header
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		// Parse Accept-Language header (simplified)
		languages := strings.Split(acceptLang, ",")
		for _, lang := range languages {
			// Remove quality factor (e.g., "en-US;q=0.9" -> "en-US")
			lang = strings.Split(strings.TrimSpace(lang), ";")[0]

			// Extract main language (e.g., "en-US" -> "en")
			mainLang := strings.Split(lang, "-")[0]

			if l.isSupported(mainLang) {
				return mainLang
			}
		}
	}

	// 3. Default locale
	return l.defaultLocale
}

// isSupported checks if a locale is supported
func (l *Localizer) isSupported(locale string) bool {
	for _, supported := range l.supportedLocales {
		if supported == locale {
			return true
		}
	}
	return false
}

// Localize returns a localized message
func (l *Localizer) Localize(ctx context.Context, key string, data map[string]interface{}) string {
	locale := l.getLocaleFromContext(ctx)
	return l.LocalizeWithLocale(locale, key, data)
}

// LocalizeWithLocale returns a localized message for specific locale
func (l *Localizer) LocalizeWithLocale(locale, key string, data map[string]interface{}) string {
	// Try requested locale
	if message := l.getMessage(locale, key); message != "" {
		return l.processTemplate(message, data)
	}

	// Fallback to default locale
	if locale != l.defaultLocale {
		if message := l.getMessage(l.defaultLocale, key); message != "" {
			return l.processTemplate(message, data)
		}
	}

	// Last resort: return the key itself
	return key
}

// getMessage retrieves a message from the locale messages
func (l *Localizer) getMessage(locale, key string) string {
	messages, exists := l.messages[locale]
	if !exists {
		return ""
	}

	// Split key by dots for nested access (e.g., "errors.not_found")
	parts := strings.Split(key, ".")
	current := messages

	for i, part := range parts {
		switch v := current[part].(type) {
		case string:
			if i == len(parts)-1 {
				return v
			}
			return "" // Not a leaf node but expected to be
		case map[string]interface{}:
			if i == len(parts)-1 {
				return "" // Expected string but got object
			}
			current = v
		default:
			return ""
		}
	}

	return ""
}

// processTemplate processes template variables in messages
func (l *Localizer) processTemplate(message string, data map[string]interface{}) string {
	if data == nil {
		return message
	}

	tmpl, err := template.New("message").Parse(message)
	if err != nil {
		return message // Return original if template parsing fails
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return message // Return original if template execution fails
	}

	return result.String()
}

// getLocaleFromContext extracts locale from context
func (l *Localizer) getLocaleFromContext(ctx context.Context) string {
	if locale, ok := ctx.Value(localeKey).(string); ok {
		return locale
	}
	return l.defaultLocale
}

// GetSupportedLocales returns list of supported locales
func (l *Localizer) GetSupportedLocales() []string {
	return l.supportedLocales
}

// Helper functions for common messages

// T is a shorthand for Localize
func (l *Localizer) T(ctx context.Context, key string, data ...map[string]interface{}) string {
	var templateData map[string]interface{}
	if len(data) > 0 {
		templateData = data[0]
	}
	return l.Localize(ctx, key, templateData)
}
