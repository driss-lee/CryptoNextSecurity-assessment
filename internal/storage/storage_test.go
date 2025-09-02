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

func TestInMemoryStorage_GetByID_And_DeleteByID(t *testing.T) {
	storage := NewInMemoryStorage(10)
	ctx := context.Background()

	p := models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 100)
	_ = storage.Store(ctx, p)

	// GetByID should return the packet
	got, err := storage.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != p.ID {
		t.Fatalf("expected packet %s, got %#v", p.ID, got)
	}

	// DeleteByID should remove it
	if err := storage.DeleteByID(ctx, p.ID); err != nil {
		t.Fatalf("unexpected error on delete: %v", err)
	}
	got, err = storage.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("unexpected error after delete: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil after delete, got %#v", got)
	}
}

func TestInMemoryStorage_Clear(t *testing.T) {
	storage := NewInMemoryStorage(10)
	ctx := context.Background()

	_ = storage.Store(ctx, models.NewPacket("192.168.1.1", "8.8.8.8", "TCP", 80, 100))
	_ = storage.Store(ctx, models.NewPacket("192.168.1.2", "1.1.1.1", "UDP", 53, 100))

	if err := storage.Clear(ctx); err != nil {
		t.Fatalf("unexpected error clearing: %v", err)
	}

	resp, err := storage.Get(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error getting after clear: %v", err)
	}
	if resp.Total != 0 {
		t.Fatalf("expected 0 total after clear, got %d", resp.Total)
	}
}

func TestInMemoryStorage_Stats(t *testing.T) {
	storage := NewInMemoryStorage(5)
	ctx := context.Background()

	// Initially empty
	s, err := storage.Stats(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting stats: %v", err)
	}
	if s.TotalPackets != 0 || s.Capacity != 5 {
		t.Fatalf("unexpected stats when empty: %#v", s)
	}
	if s.OldestAt != nil || s.NewestAt != nil {
		t.Fatalf("expected nil timestamps when empty: %#v", s)
	}

	// Add two packets
	_ = storage.Store(ctx, models.NewPacket("10.0.0.1", "8.8.4.4", "TCP", 443, 200))
	time.Sleep(2 * time.Millisecond)
	_ = storage.Store(ctx, models.NewPacket("10.0.0.2", "8.8.8.8", "UDP", 53, 60))

	s, err = storage.Stats(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting stats: %v", err)
	}
	if s.TotalPackets != 2 || s.Capacity != 5 {
		t.Fatalf("unexpected stats after add: %#v", s)
	}
	if s.OldestAt == nil || s.NewestAt == nil || !s.NewestAt.After(*s.OldestAt) && !s.NewestAt.Equal(*s.OldestAt) {
		t.Fatalf("expected non-nil timestamps with ordering, got %#v", s)
	}
}
