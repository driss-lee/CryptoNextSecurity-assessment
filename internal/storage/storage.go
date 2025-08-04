package storage

import (
	"context"
	"sync"
	"time"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
)

// Storage defines the interface for packet storage
type Storage interface {
	// Store adds a packet to storage
	Store(ctx context.Context, packet *models.Packet) error

	// Get retrieves packets with optional filtering
	Get(ctx context.Context, filter *models.PacketFilter) (*models.PacketResponse, error)
}

// InMemoryStorage implements Storage interface with in-memory storage
type InMemoryStorage struct {
	packets map[string]*models.Packet
	mutex   sync.RWMutex
	maxSize int
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage(maxSize int) *InMemoryStorage {
	return &InMemoryStorage{
		packets: make(map[string]*models.Packet),
		maxSize: maxSize,
	}
}

// Store adds a packet to storage
func (s *InMemoryStorage) Store(ctx context.Context, packet *models.Packet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if we need to remove old packets to make room
	if len(s.packets) >= s.maxSize {
		s.removeOldestPacket()
	}

	s.packets[packet.ID] = packet
	return nil
}

// Get retrieves packets with optional filtering
func (s *InMemoryStorage) Get(ctx context.Context, filter *models.PacketFilter) (*models.PacketResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var packets []models.Packet

	for _, packet := range s.packets {
		if s.matchesFilter(packet, filter) {
			packets = append(packets, *packet)
		}
	}

	// Apply pagination
	if filter != nil && filter.Limit > 0 {
		start := filter.Offset
		end := start + filter.Limit
		if start >= len(packets) {
			packets = []models.Packet{}
		} else if end > len(packets) {
			packets = packets[start:]
		} else {
			packets = packets[start:end]
		}
	}

	return &models.PacketResponse{
		Packets:   packets,
		Total:     len(packets),
		Timestamp: time.Now(),
	}, nil
}

// matchesFilter checks if a packet matches the given filter
func (s *InMemoryStorage) matchesFilter(packet *models.Packet, filter *models.PacketFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Protocol != "" && packet.Protocol != filter.Protocol {
		return false
	}

	if filter.SourceIP != "" && packet.SourceIP != filter.SourceIP {
		return false
	}

	if filter.DestinationIP != "" && packet.DestinationIP != filter.DestinationIP {
		return false
	}

	if !filter.FromTimestamp.IsZero() && packet.Timestamp.Before(filter.FromTimestamp) {
		return false
	}

	if !filter.ToTimestamp.IsZero() && packet.Timestamp.After(filter.ToTimestamp) {
		return false
	}

	return true
}

// removeOldestPacket removes the oldest packet to make room for new ones
func (s *InMemoryStorage) removeOldestPacket() {
	var oldestID string
	var oldestTime time.Time
	first := true

	for id, packet := range s.packets {
		if first {
			oldestID = id
			oldestTime = packet.Timestamp
			first = false
		} else if packet.Timestamp.Before(oldestTime) {
			oldestID = id
			oldestTime = packet.Timestamp
		}
	}

	if oldestID != "" {
		delete(s.packets, oldestID)
	}
}
