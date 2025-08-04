package storage

import (
	"context"
	"testing"
	"time"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
)

func TestInMemoryStorage_Store(t *testing.T) {
	storage := NewInMemoryStorage(100)
	ctx := context.Background()

	packet := models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 1500)

	err := storage.Store(ctx, packet)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify packet was stored by getting all packets
	response, err := storage.Get(ctx, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Total != 1 {
		t.Errorf("Expected 1 packet, got %d", response.Total)
	}
}

func TestInMemoryStorage_Get(t *testing.T) {
	storage := NewInMemoryStorage(100)
	ctx := context.Background()

	// Store multiple packets
	packet1 := models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 1500)
	packet2 := models.NewPacket("192.168.1.2", "1.1.1.1", "UDP", 53, 512)

	storage.Store(ctx, packet1)
	storage.Store(ctx, packet2)

	// Test getting all packets
	response, err := storage.Get(ctx, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Total != 2 {
		t.Errorf("Expected 2 packets, got %d", response.Total)
	}
	if len(response.Packets) != 2 {
		t.Errorf("Expected 2 packets in response, got %d", len(response.Packets))
	}
}

func TestInMemoryStorage_GetWithFilter(t *testing.T) {
	storage := NewInMemoryStorage(100)
	ctx := context.Background()

	// Store packets with different protocols
	packet1 := models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 1500)
	packet2 := models.NewPacket("192.168.1.2", "1.1.1.1", "UDP", 53, 512)
	packet3 := models.NewPacket("192.168.1.3", "142.250.190.78", "TCP", 443, 1500)

	storage.Store(ctx, packet1)
	storage.Store(ctx, packet2)
	storage.Store(ctx, packet3)

	// Test filtering by protocol
	filter := &models.PacketFilter{Protocol: "TCP"}
	response, err := storage.Get(ctx, filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Total != 2 {
		t.Errorf("Expected 2 TCP packets, got %d", response.Total)
	}

	// Test filtering by source IP
	filter = &models.PacketFilter{SourceIP: "192.168.1.1"}
	response, err = storage.Get(ctx, filter)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Total != 1 {
		t.Errorf("Expected 1 packet from 192.168.1.1, got %d", response.Total)
	}
	if response.Packets[0].SourceIP != "192.168.1.1" {
		t.Errorf("Expected source IP 192.168.1.1, got %s", response.Packets[0].SourceIP)
	}
}

func TestInMemoryStorage_MaxSize(t *testing.T) {
	storage := NewInMemoryStorage(2) // Only allow 2 packets
	ctx := context.Background()

	// Store 3 packets
	packet1 := models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 1500)
	packet2 := models.NewPacket("192.168.1.2", "1.1.1.1", "UDP", 53, 512)
	packet3 := models.NewPacket("192.168.1.3", "142.250.190.78", "TCP", 443, 1500)

	storage.Store(ctx, packet1)
	time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	storage.Store(ctx, packet2)
	time.Sleep(10 * time.Millisecond)
	storage.Store(ctx, packet3)

	// Should only have 2 packets (oldest one should be removed)
	response, err := storage.Get(ctx, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response.Total != 2 {
		t.Errorf("Expected 2 packets, got %d", response.Total)
	}
}
