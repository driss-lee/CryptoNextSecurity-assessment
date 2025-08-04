package sniffing

import (
	"context"
	"testing"
	"time"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStorage implements Storage interface for testing
type MockStorage struct {
	packets []*models.Packet
}

func (m *MockStorage) Store(ctx context.Context, packet *models.Packet) error {
	m.packets = append(m.packets, packet)
	return nil
}

func TestPacketSniffer_Start(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 100*time.Millisecond)
	ctx := context.Background()

	// Start sniffing
	err := sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Wait a bit for packets to be generated
	time.Sleep(200 * time.Millisecond)

	// Stop sniffing
	err = sniffer.Stop(ctx)
	require.NoError(t, err)
	assert.False(t, sniffer.IsRunning())

	// Verify packets were generated
	assert.Greater(t, len(mockStorage.packets), 0)
}

func TestPacketSniffer_Stop(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 100*time.Millisecond)
	ctx := context.Background()

	// Start sniffing
	err := sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Stop sniffing
	err = sniffer.Stop(ctx)
	require.NoError(t, err)
	assert.False(t, sniffer.IsRunning())

	// Try to stop again (should not error)
	err = sniffer.Stop(ctx)
	require.NoError(t, err)
}

func TestPacketSniffer_IsRunning(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 1*time.Second)

	// Initially not running
	assert.False(t, sniffer.IsRunning())

	// Start sniffing
	ctx := context.Background()
	err := sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Stop sniffing
	err = sniffer.Stop(ctx)
	require.NoError(t, err)
	assert.False(t, sniffer.IsRunning())
}

func TestPacketSniffer_StartWhenAlreadyRunning(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 1*time.Second)
	ctx := context.Background()

	// Start sniffing
	err := sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Try to start again (should not error)
	err = sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Clean up
	sniffer.Stop(ctx)
}

func TestPacketSniffer_ContextCancellation(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 50*time.Millisecond)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Start sniffing
	err := sniffer.Start(ctx)
	require.NoError(t, err)
	assert.True(t, sniffer.IsRunning())

	// Wait a bit for packets to be generated
	time.Sleep(100 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for sniffing to stop
	time.Sleep(100 * time.Millisecond)

	// Verify sniffing stopped
	assert.False(t, sniffer.IsRunning())

	// Verify packets were generated
	assert.Greater(t, len(mockStorage.packets), 0)
}

func TestGenerateRandomPacket(t *testing.T) {
	mockStorage := &MockStorage{}
	sniffer := NewPacketSniffer(mockStorage, 1*time.Second)

	packet := sniffer.generateRandomPacket()

	// Verify packet has required fields
	assert.NotEmpty(t, packet.ID)
	assert.NotEmpty(t, packet.SourceIP)
	assert.NotEmpty(t, packet.DestinationIP)
	assert.NotEmpty(t, packet.Protocol)
	assert.Greater(t, packet.Port, 0)
	assert.LessOrEqual(t, packet.Port, 65535)
	assert.Greater(t, packet.Size, 0)
	assert.LessOrEqual(t, packet.Size, 1500)
	assert.GreaterOrEqual(t, packet.Size, 64)
	assert.False(t, packet.Timestamp.IsZero())

	// Verify source and destination are different
	assert.NotEqual(t, packet.SourceIP, packet.DestinationIP)

	// Verify protocol is valid
	validProtocols := []string{"TCP", "UDP", "HTTP", "HTTPS"}
	assert.Contains(t, validProtocols, packet.Protocol)

	// Verify IP addresses are valid
	assert.Contains(t, sniffer.commonIPs, packet.SourceIP)
	assert.Contains(t, sniffer.commonIPs, packet.DestinationIP)

	// Verify port is from common ports
	assert.Contains(t, sniffer.commonPorts, packet.Port)
}
