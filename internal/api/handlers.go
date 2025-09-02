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

// GetPacketByID handles GET /packets/:id
// @Summary Get packet by ID
// @Description Retrieve a single packet by its unique ID
// @Tags packets
// @Accept json
// @Produce json
// @Param id path string true "Packet ID"
// @Success 200 {object} models.Packet
// @Failure 404 {object} ErrorResponse "Packet not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /packets/{id} [get]
func (h *Handler) GetPacketByID(c *gin.Context) {
	id := c.Param("id")
	packet, err := h.packetService.GetPacketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to retrieve packet"})
		return
	}
	if packet == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not Found", Message: "Packet not found"})
		return
	}
	c.JSON(http.StatusOK, packet)
}

// DeletePacketByID handles DELETE /packets/:id
// @Summary Delete packet by ID
// @Description Delete a single packet by its unique ID
// @Tags packets
// @Param id path string true "Packet ID"
// @Success 204 "Deleted"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /packets/{id} [delete]
func (h *Handler) DeletePacketByID(c *gin.Context) {
	id := c.Param("id")
	if err := h.packetService.DeletePacketByID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to delete packet"})
		return
	}
	c.Status(http.StatusNoContent)
}

// ClearPackets handles DELETE /packets
// @Summary Clear all packets
// @Description Remove all packets from storage
// @Tags packets
// @Success 204 "Cleared"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /packets [delete]
func (h *Handler) ClearPackets(c *gin.Context) {
	if err := h.packetService.ClearPackets(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to clear packets"})
		return
	}
	c.Status(http.StatusNoContent)
}

// Health handles GET /health
// @Summary Health check
// @Description Service health status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	status := map[string]interface{}{
		"status":  "ok",
		"running": h.packetService.IsSniffingRunning(),
	}
	c.JSON(http.StatusOK, status)
}

// Stats handles GET /stats
// @Summary Storage statistics
// @Description Get current storage statistics
// @Tags system
// @Produce json
// @Success 200 {object} models.Stats
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stats [get]
func (h *Handler) Stats(c *gin.Context) {
	stats, err := h.packetService.StorageStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to get stats"})
		return
	}
	if stats == nil {
		stats = &models.Stats{TotalPackets: 0}
	}
	c.JSON(http.StatusOK, stats)
}

// StartSniffing handles POST /sniffing/start
// @Summary Start sniffing
// @Description Start the packet sniffing process
// @Tags sniffing
// @Success 204 "Started"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /sniffing/start [post]
func (h *Handler) StartSniffing(c *gin.Context) {
	if err := h.packetService.StartSniffing(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to start sniffing"})
		return
	}
	c.Status(http.StatusNoContent)
}

// StopSniffing handles POST /sniffing/stop
// @Summary Stop sniffing
// @Description Stop the packet sniffing process
// @Tags sniffing
// @Success 204 "Stopped"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /sniffing/stop [post]
func (h *Handler) StopSniffing(c *gin.Context) {
	if err := h.packetService.StopSniffing(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error", Message: "Failed to stop sniffing"})
		return
	}
	c.Status(http.StatusNoContent)
}

// SniffingStatus handles GET /sniffing/status
// @Summary Sniffing status
// @Description Get current sniffing service status
// @Tags sniffing
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /sniffing/status [get]
func (h *Handler) SniffingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]bool{"running": h.packetService.IsSniffingRunning()})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
