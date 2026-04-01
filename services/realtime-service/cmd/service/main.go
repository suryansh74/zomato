package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/suryansh74/zomato/services/realtime-service/internal/config"
	"github.com/suryansh74/zomato/services/realtime-service/internal/middleware"
	ws "github.com/suryansh74/zomato/services/realtime-service/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// VERY IMPORTANT FOR LOCAL DEVELOPMENT: Allow all origins
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Content-Type", "X-Internal-Key"},
	}))

	manager := ws.NewManager()

	// ---------------------------------------------------------
	// 1. PUBLIC WEBSOCKET ROUTE (For the React Dashboard)
	// ---------------------------------------------------------
	r.Get("/ws/restaurant/{restaurant_id}", func(w http.ResponseWriter, r *http.Request) {
		restaurantID := chi.URLParam(r, "restaurant_id")

		// Upgrade the standard HTTP request to a persistent WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket Upgrade Error:", err)
			return
		}

		// Register the new connection in our manager
		manager.AddClient(restaurantID, conn)

		// 🚨 THE FIX: Do NOT use a goroutine here!
		// We MUST block this function so Go doesn't close the connection.
		defer manager.RemoveClient(restaurantID, conn)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				// User closed the browser tab or internet dropped
				break
			}
		}
	})

	// ---------------------------------------------------------
	// 2. INTERNAL PROTECTED ROUTE (For the Go Order Service)
	// ---------------------------------------------------------
	r.Group(func(protected chi.Router) {
		// Apply the VIP Bouncer middleware!
		protected.Use(middleware.InternalAPIAuth)

		protected.Post("/api/internal/notify-order", func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid payload", http.StatusBadRequest)
				return
			}

			// Extract which restaurant needs to see the order
			restaurantID, ok := payload["restaurant_id"].(string)
			if !ok {
				http.Error(w, "Missing restaurant_id", http.StatusBadRequest)
				return
			}

			// Broadcast the order details to that specific restaurant's dashboard!
			manager.BroadcastToRestaurant(restaurantID, payload)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))
		})
	})

	log.Printf("realtime-service server started on %s:%s", cfg.Host, cfg.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), r)
}
