package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/aryadisastra/main-app-be/internal/dto"
	"github.com/aryadisastra/main-app-be/internal/httpx"
	"github.com/aryadisastra/main-app-be/internal/models"
)

func genTrackingNumber() string {
	return time.Now().UTC().Format("20060102T150405.000000000")
}

// CreateShipment godoc
// @Summary      Create Shipment
// @Description  Membuat pesanan pengiriman untuk user yang sedang login.
// @Tags         shipments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.CreateShipmentRequest  true  "Create shipment payload"
// @Success      201  {object}  httpx.Envelope{data=dto.ShipmentResponse}
// @Failure      400  {object}  httpx.Envelope
// @Failure      401  {object}  httpx.Envelope
// @Router       /api/v1/shipments [post]
func CreateShipment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateShipmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		uidAny, _ := c.Get("user_id")
		uid, _ := uidAny.(string)

		ship := models.Shipment{
			TrackingNumber:  genTrackingNumber(),
			UserID:          uid,
			SenderName:      req.SenderName,
			SenderAddress:   req.SenderAddress,
			ReceiverName:    req.ReceiverName,
			ReceiverAddress: req.ReceiverAddress,
			ItemDescription: req.ItemDescription,
			Status:          "Created",
		}
		if err := db.Create(&ship).Error; err != nil {
			httpx.Fail(c, http.StatusBadRequest, "could not create shipment")
			return
		}
		httpx.Created(c, dto.ShipmentResponse{
			ID: ship.ID, TrackingNumber: ship.TrackingNumber, UserID: ship.UserID,
			SenderName: ship.SenderName, SenderAddress: ship.SenderAddress,
			ReceiverName: ship.ReceiverName, ReceiverAddress: ship.ReceiverAddress,
			ItemDescription: ship.ItemDescription, Status: ship.Status,
		})
	}
}

// UpdateShipmentStatus godoc
// @Summary      Update Shipment Status (admin only)
// @Description  Mengubah status pengiriman. Hanya role `admin` yang diizinkan.
// @Tags         shipments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        trackingNumber  path      string                 true  "Tracking Number"
// @Param        payload         body      dto.UpdateStatusRequest true  "Update status payload"
// @Success      200  {object}  httpx.Envelope{data=dto.ShipmentResponse}
// @Failure      400  {object}  httpx.Envelope
// @Failure      401  {object}  httpx.Envelope
// @Failure      403  {object}  httpx.Envelope
// @Failure      404  {object}  httpx.Envelope
// @Router       /api/v1/shipments/{trackingNumber}/status [patch]
func UpdateShipmentStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := c.Param("trackingNumber")
		var req dto.UpdateStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			httpx.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		var ship models.Shipment
		if err := db.First(&ship, "tracking_number = ?", tr).Error; err != nil {
			httpx.Fail(c, http.StatusNotFound, "shipment not found")
			return
		}
		ship.Status = req.Status
		if err := db.Save(&ship).Error; err != nil {
			httpx.Fail(c, http.StatusBadRequest, "could not update status")
			return
		}
		httpx.OK(c, dto.ShipmentResponse{
			ID: ship.ID, TrackingNumber: ship.TrackingNumber, UserID: ship.UserID,
			SenderName: ship.SenderName, SenderAddress: ship.SenderAddress,
			ReceiverName: ship.ReceiverName, ReceiverAddress: ship.ReceiverAddress,
			ItemDescription: ship.ItemDescription, Status: ship.Status,
		})
	}
}

// TrackShipment godoc
// @Summary      Track Shipment
// @Description  Melacak status pengiriman berdasarkan nomor resi.
// @Tags         shipments
// @Security     BearerAuth
// @Produce      json
// @Param        trackingNumber  path  string  true  "Tracking Number"
// @Success      200  {object}  httpx.Envelope{data=dto.ShipmentResponse}
// @Failure      401  {object}  httpx.Envelope
// @Failure      404  {object}  httpx.Envelope
// @Router       /api/v1/shipments/track/{trackingNumber} [get]
func TrackShipment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tr := c.Param("trackingNumber")
		var ship models.Shipment
		if err := db.First(&ship, "tracking_number = ?", tr).Error; err != nil {
			httpx.Fail(c, http.StatusNotFound, "shipment not found")
			return
		}
		httpx.OK(c, dto.ShipmentResponse{
			ID: ship.ID, TrackingNumber: ship.TrackingNumber, UserID: ship.UserID,
			SenderName: ship.SenderName, SenderAddress: ship.SenderAddress,
			ReceiverName: ship.ReceiverName, ReceiverAddress: ship.ReceiverAddress,
			ItemDescription: ship.ItemDescription, Status: ship.Status,
		})
	}
}

// GetShipments godoc
// @Summary      Get My Shipments
// @Description  Mengambil semua pengiriman milik user yang sedang login.
// @Tags         shipments
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  httpx.Envelope{data=[]dto.ShipmentResponse}
// @Failure      401  {object}  httpx.Envelope
// @Router       /api/v1/shipments [get]
func GetShipments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidAny, _ := c.Get("user_id")
		uid, _ := uidAny.(string)
		var ships []models.Shipment
		if err := db.Where("user_id = ?", uid).Order("created_at DESC").Find(&ships).Error; err != nil {
			httpx.Fail(c, http.StatusBadRequest, "could not get shipments")
			return
		}
		resp := make([]dto.ShipmentResponse, 0, len(ships))
		for _, s := range ships {
			resp = append(resp, dto.ShipmentResponse{
				ID: s.ID, TrackingNumber: s.TrackingNumber, UserID: s.UserID,
				SenderName: s.SenderName, SenderAddress: s.SenderAddress,
				ReceiverName: s.ReceiverName, ReceiverAddress: s.ReceiverAddress,
				ItemDescription: s.ItemDescription, Status: s.Status,
			})
		}
		httpx.OK(c, resp)
	}
}
