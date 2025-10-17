package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
)

func Language(localizer *i18n.Localizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acceptLanguage := r.Header.Get("Accept-Language")

			language := "en"

			if acceptLanguage != "" {
				languages := strings.Split(acceptLanguage, ",")
				for _, lang := range languages {
					lang = strings.TrimSpace(strings.Split(lang, ";")[0])

					langCode := strings.Split(lang, "-")[0]

					if langCode == "en" || langCode == "pt" || langCode == "es" {
						language = langCode
						break
					}
				}
			}

			ctx := context.WithValue(r.Context(), "language", language)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
