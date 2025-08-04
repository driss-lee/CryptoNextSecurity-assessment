package sniffing

import (
	"context"
	"math/rand"
	"time"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
)

// Sniffer defines the interface for packet sniffing
type Sniffer interface {
	// Start begins the sniffing process
	Start(ctx context.Context) error

	// Stop stops the sniffing process
	Stop(ctx context.Context) error

	// IsRunning returns true if sniffing is active
	IsRunning() bool
}

// PacketSniffer implements the Sniffer interface with simulated packet capture
type PacketSniffer struct {
	storage     Storage
	interval    time.Duration
	isRunning   bool
	stopChan    chan struct{}
	commonIPs   []string
	commonPorts []int
	protocols   []string
}

// Storage defines the interface for packet storage
type Storage interface {
	Store(ctx context.Context, packet *models.Packet) error
}

// NewPacketSniffer creates a new packet sniffer instance
func NewPacketSniffer(storage Storage, interval time.Duration) *PacketSniffer {
	return &PacketSniffer{
		storage:  storage,
		interval: interval,
		stopChan: make(chan struct{}),
		commonIPs: []string{
			"192.168.1.1", "192.168.1.100", "192.168.1.101", "192.168.1.102",
			"10.0.0.1", "10.0.0.2", "10.0.0.3", "10.0.0.4",
			"172.16.0.1", "172.16.0.2", "172.16.0.3",
			"8.8.8.8", "1.1.1.1", "208.67.222.222", // DNS servers
			"142.250.190.78", "151.101.1.69", "104.16.124.96", // Google, Reddit, Cloudflare
		},
		commonPorts: []int{
			80, 443, 22, 21, 25, 53, 110, 143, 993, 995, // Common ports
			8080, 8443, 3000, 5000, 8000, 9000, // Development ports
		},
		protocols: []string{"TCP", "UDP", "HTTP", "HTTPS"},
	}
}

// Start begins the sniffing process
func (s *PacketSniffer) Start(ctx context.Context) error {
	if s.isRunning {
		return nil
	}

	s.isRunning = true

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.isRunning = false
				return
			case <-s.stopChan:
				s.isRunning = false
				return
			case <-ticker.C:
				s.generateAndStorePacket(ctx)
			}
		}
	}()

	return nil
}

// Stop stops the sniffing process
func (s *PacketSniffer) Stop(ctx context.Context) error {
	if !s.isRunning {
		return nil
	}

	close(s.stopChan)
	s.isRunning = false
	return nil
}

// IsRunning returns true if sniffing is active
func (s *PacketSniffer) IsRunning() bool {
	return s.isRunning
}

// generateAndStorePacket creates a simulated packet and stores it
func (s *PacketSniffer) generateAndStorePacket(ctx context.Context) {
	packet := s.generateRandomPacket()

	if err := s.storage.Store(ctx, packet); err != nil {
		// In a real application, we might log this error
		// For now, we'll just ignore it to keep the simulation running
		_ = err
	}
}

// generateRandomPacket creates a realistic packet with random data
func (s *PacketSniffer) generateRandomPacket() *models.Packet {
	// Generate random source and destination IPs
	sourceIP := s.commonIPs[rand.Intn(len(s.commonIPs))]
	destIP := s.commonIPs[rand.Intn(len(s.commonIPs))]

	// Avoid same source and destination
	for destIP == sourceIP {
		destIP = s.commonIPs[rand.Intn(len(s.commonIPs))]
	}

	// Generate random port
	port := s.commonPorts[rand.Intn(len(s.commonPorts))]

	// Generate random protocol
	protocol := s.protocols[rand.Intn(len(s.protocols))]

	// Generate random packet size (64-1500 bytes)
	size := rand.Intn(1436) + 64

	// Create packet
	packet := models.NewPacket(sourceIP, destIP, protocol, port, size)

	// Add some realistic variations
	if rand.Float32() < 0.3 {
		packet.TTL = rand.Intn(64) + 32
	}

	if rand.Float32() < 0.2 {
		flags := []string{"SYN", "ACK", "FIN", "RST", "PSH", "URG"}
		packet.Flags = flags[rand.Intn(len(flags))]
	}

	// Add payload for HTTP/HTTPS packets
	if protocol == "HTTP" || protocol == "HTTPS" {
		payloads := []string{
			"GET / HTTP/1.1",
			"POST /api/data HTTP/1.1",
			"PUT /resource HTTP/1.1",
			"DELETE /item/123 HTTP/1.1",
		}
		packet.Payload = payloads[rand.Intn(len(payloads))]
	}

	return packet
}
