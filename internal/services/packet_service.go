package services

import (
	"context"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
	"github.com/cryptonextsecurity/network-sniffer/internal/storage"
	"github.com/cryptonextsecurity/network-sniffer/pkg/sniffing"
)

// PacketService handles business logic for packet operations
type PacketService struct {
	storage storage.Storage
	sniffer sniffing.Sniffer
}

// NewPacketService creates a new packet service instance
func NewPacketService(storage storage.Storage, sniffer sniffing.Sniffer, logger any) *PacketService {
	return &PacketService{
		storage: storage,
		sniffer: sniffer,
	}
}

// StartSniffing begins the packet sniffing process
func (s *PacketService) StartSniffing(ctx context.Context) error {
	return s.sniffer.Start(ctx)
}

// StopSniffing stops the packet sniffing process
func (s *PacketService) StopSniffing(ctx context.Context) error {
	return s.sniffer.Stop(ctx)
}

// IsSniffingRunning returns true if sniffing is active
func (s *PacketService) IsSniffingRunning() bool {
	return s.sniffer.IsRunning()
}

// GetPackets retrieves packets with optional filtering
func (s *PacketService) GetPackets(ctx context.Context, filter *models.PacketFilter) (*models.PacketResponse, error) {
	return s.storage.Get(ctx, filter)
}
