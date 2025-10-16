package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
)

type languageContextKey string

const LanguageKey languageContextKey = "language"

// LanguageMiddleware extracts language from Accept-Language header or query parameter
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var language string

		// Check query parameter first (for explicit language selection)
		if lang := r.URL.Query().Get("lang"); lang != "" {
			language = lang
		} else {
			// Parse Accept-Language header
			acceptLang := r.Header.Get("Accept-Language")
			language = parseAcceptLanguage(acceptLang)
		}

		// Validate and set default if not supported
		if !i18n.IsLanguageSupported(language) {
			language = "en"
		}

		// Add language to context
		ctx := context.WithValue(r.Context(), LanguageKey, language)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseAcceptLanguage extracts the primary language from Accept-Language header
func parseAcceptLanguage(header string) string {
	if header == "" {
		return "en"
	}

	// Split by comma and get first preference
	langs := strings.Split(header, ",")
	if len(langs) == 0 {
		return "en"
	}

	// Get first language and remove quality value if present
	firstLang := strings.TrimSpace(langs[0])
	if idx := strings.Index(firstLang, ";"); idx != -1 {
		firstLang = firstLang[:idx]
	}

	// Extract language code (first 2 characters)
	if len(firstLang) >= 2 {
		return strings.ToLower(firstLang[:2])
	}

	return "en"
}

// GetLanguageFromContext extracts language from request context
func GetLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return "en"
}
