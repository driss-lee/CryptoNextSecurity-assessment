package api

import (
	"net/http"
	"strconv"

	"github.com/cryptonextsecurity/network-sniffer/internal/models"
	"github.com/cryptonextsecurity/network-sniffer/internal/services"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	packetService *services.PacketService
}

// PacketService returns the packet service instance
func (h *Handler) PacketService() *services.PacketService {
	return h.packetService
}

// NewHandler creates a new handler instance
func NewHandler(packetService *services.PacketService, logger interface{}) *Handler {
	return &Handler{
		packetService: packetService,
	}
}

// GetPackets handles GET /packets requests
// @Summary Get all packets
// @Description Retrieve all sniffed packets with optional filtering
// @Tags packets
// @Accept json
// @Produce json
// @Param protocol query string false "Filter by protocol (TCP, UDP, HTTP, HTTPS)"
// @Param source_ip query string false "Filter by source IP address"
// @Param destination_ip query string false "Filter by destination IP address"
// @Param limit query int false "Limit number of results (default: no limit)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} models.PacketResponse "List of packets"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /packets [get]
func (h *Handler) GetPackets(c *gin.Context) {
	// Parse query parameters
	filter := &models.PacketFilter{}

	if protocol := c.Query("protocol"); protocol != "" {
		filter.Protocol = protocol
	}

	if sourceIP := c.Query("source_ip"); sourceIP != "" {
		filter.SourceIP = sourceIP
	}

	if destIP := c.Query("destination_ip"); destIP != "" {
		filter.DestinationIP = destIP
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Get packets from service
	response, err := h.packetService.GetPackets(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to retrieve packets",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
