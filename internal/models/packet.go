package models

import (
	"fmt"
	"time"
)

// Packet represents a network packet with metadata
type Packet struct {
	ID            string    `json:"id" validate:"required"`
	SourceIP      string    `json:"source_ip" validate:"required,ip"`
	DestinationIP string    `json:"destination_ip" validate:"required,ip"`
	Protocol      string    `json:"protocol" validate:"required,oneof=TCP UDP ICMP HTTP HTTPS"`
	Port          int       `json:"port" validate:"min=1,max=65535"`
	Size          int       `json:"size" validate:"min=1"`
	Timestamp     time.Time `json:"timestamp" validate:"required"`
	TTL           int       `json:"ttl,omitempty"`
	Flags         string    `json:"flags,omitempty"`
	Payload       string    `json:"payload,omitempty"`
}

// PacketResponse represents the API response for packets
type PacketResponse struct {
	Packets   []Packet  `json:"packets"`
	Total     int       `json:"total"`
	Timestamp time.Time `json:"timestamp"`
}

// PacketFilter represents filtering options for packets
type PacketFilter struct {
	Protocol      string    `json:"protocol,omitempty"`
	SourceIP      string    `json:"source_ip,omitempty"`
	DestinationIP string    `json:"destination_ip,omitempty"`
	FromTimestamp time.Time `json:"from_timestamp,omitempty"`
	ToTimestamp   time.Time `json:"to_timestamp,omitempty"`
	Limit         int       `json:"limit,omitempty"`
	Offset        int       `json:"offset,omitempty"`
}

// Stats contains basic storage statistics
type Stats struct {
	TotalPackets int        `json:"total_packets"`
	Capacity     int        `json:"capacity"`
	OldestAt     *time.Time `json:"oldest_at,omitempty"`
	NewestAt     *time.Time `json:"newest_at,omitempty"`
}

// NewPacket creates a new packet with default values
func NewPacket(sourceIP, destIP, protocol string, port, size int) *Packet {
	return &Packet{
		ID:            generatePacketID(),
		SourceIP:      sourceIP,
		DestinationIP: destIP,
		Protocol:      protocol,
		Port:          port,
		Size:          size,
		Timestamp:     time.Now(),
		TTL:           64,
		Flags:         "SYN",
	}
}

// generatePacketID creates a unique packet ID
func generatePacketID() string {
	return "packet_" + time.Now().Format("20060102150405") + "_" + fmt.Sprintf("%09d", time.Now().UnixNano()%1000000000)
}
