package middleware

import (
	"net/http"

	"github.com/suryansh74/zomato/services/shared/helper"
	"github.com/suryansh74/zomato/services/shared/middleware"
	"github.com/suryansh74/zomato/services/shared/token"
)

func IsRestaurantOwner() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := r.Context().Value(middleware.UserContextKey).(*token.Payload)
			if payload.User.Role != "restaurant_owner" {
				helper.WriteJSON(w, http.StatusForbidden, map[string]string{
					"error": "forbidden: insufficient role",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
