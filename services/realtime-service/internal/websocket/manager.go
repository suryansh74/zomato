package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Manager keeps track of active connections grouped by Restaurant ID
type Manager struct {
	sync.RWMutex
	// Key: RestaurantID, Value: Map of connections (using map for easy deletion)
	clients map[string]map[*websocket.Conn]bool
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

// Add a new dashboard connection for a specific restaurant
func (m *Manager) AddClient(restaurantID string, conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()

	if m.clients[restaurantID] == nil {
		m.clients[restaurantID] = make(map[*websocket.Conn]bool)
	}
	m.clients[restaurantID][conn] = true
	log.Printf("New dashboard connected for restaurant: %s", restaurantID)
}

// Remove a connection when the owner closes the tab
func (m *Manager) RemoveClient(restaurantID string, conn *websocket.Conn) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[restaurantID]; ok {
		delete(m.clients[restaurantID], conn)
		conn.Close()
		log.Printf("Dashboard disconnected for restaurant: %s", restaurantID)
	}
}

// Broadcast sends a JSON payload to all open dashboards for a specific restaurant
func (m *Manager) BroadcastToRestaurant(restaurantID string, payload interface{}) {
	m.RLock()
	defer m.RUnlock()

	clients, exists := m.clients[restaurantID]
	if !exists || len(clients) == 0 {
		log.Printf("No active dashboards found for restaurant: %s. Skipping broadcast.", restaurantID)
		return
	}

	for conn := range clients {
		err := conn.WriteJSON(payload)
		if err != nil {
			log.Println("Error broadcasting to connection:", err)
			m.RemoveClient(restaurantID, conn) // Kick out dead connections
		}
	}
	log.Printf("Successfully broadcasted event to restaurant: %s", restaurantID)
}
