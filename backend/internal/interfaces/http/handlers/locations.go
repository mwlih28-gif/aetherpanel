package handlers

import (
	"net/http"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/gofiber/fiber/v2"
)

type CreateLocationRequest struct {
	Short string `json:"short" validate:"required,min=1,max=60"`
	Long  string `json:"long" validate:"required,min=1,max=191"`
}

type UpdateLocationRequest struct {
	Short string `json:"short" validate:"required,min=1,max=60"`
	Long  string `json:"long" validate:"required,min=1,max=191"`
}

// GetLocations returns all locations
func (h *Handler) GetLocations(c *fiber.Ctx) error {
	var locations []entities.Location
	
	if err := h.db.Find(&locations).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch locations",
		})
	}

	return c.JSON(fiber.Map{
		"data": locations,
	})
}

// CreateLocation creates a new location
func (h *Handler) CreateLocation(c *fiber.Ctx) error {
	var req CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	// Check if short code already exists
	var existingLocation entities.Location
	if err := h.db.Where("short = ?", req.Short).First(&existingLocation).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Location with this short code already exists",
		})
	}

	location := entities.Location{
		ShortCode: req.Short,
		Name:      req.Long,
	}

	if err := h.db.Create(&location).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create location",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": location,
	})
}

// GetLocation returns a specific location
func (h *Handler) GetLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var location entities.Location
	if err := h.db.Where("id = ?", id).First(&location).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": location,
	})
}

// UpdateLocation updates an existing location
func (h *Handler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var req UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	var location entities.Location
	if err := h.db.Where("id = ?", id).First(&location).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if short code already exists (excluding current location)
	var existingLocation entities.Location
	if err := h.db.Where("short = ? AND id != ?", req.Short, id).First(&existingLocation).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Location with this short code already exists",
		})
	}

	location.ShortCode = req.Short
	location.Name = req.Long

	if err := h.db.Save(&location).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update location",
		})
	}

	return c.JSON(fiber.Map{
		"data": location,
	})
}

// DeleteLocation deletes a location
func (h *Handler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	
	var location entities.Location
	if err := h.db.Where("id = ?", id).First(&location).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Location not found",
		})
	}

	// Check if location has nodes
	var nodeCount int64
	if err := h.db.Model(&entities.Node{}).Where("location_id = ?", id).Count(&nodeCount).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check location usage",
		})
	}

	if nodeCount > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete location that has nodes assigned to it",
		})
	}

	if err := h.db.Delete(&location).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete location",
		})
	}

	return c.Status(http.StatusNoContent).Send(nil)
}
