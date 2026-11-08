package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"incident-tracker-api/config"
	"incident-tracker-api/dto"
	"incident-tracker-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateIncident(c *gin.Context) {
	var input dto.CreateIncidentDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	incident := models.Incident{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
	}

	if err := config.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incident"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func ListIncidents(c *gin.Context) {
	var incidents []models.Incident

	if err := config.DB.Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch incidents"})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func GetIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid incident ID"})
		return
	}

	var incident models.Incident

	if err := config.DB.First(&incident, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch incident"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func UpdateIncidentStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid incident ID"})
		return
	}

	var input dto.UpdateIncidentStatusDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var incident models.Incident

	if err := config.DB.First(&incident, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch incident"})
		return
	}

	incident.Status = input.Status

	if err := config.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update incident"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func DeleteIncident(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid incident ID"})
		return
	}

	result := config.DB.Delete(&models.Incident{}, id)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incident"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
