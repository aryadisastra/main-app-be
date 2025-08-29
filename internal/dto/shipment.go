package dto

type CreateShipmentRequest struct {
	SenderName      string `json:"sender_name" binding:"required"`
	SenderAddress   string `json:"sender_address" binding:"required"`
	ReceiverName    string `json:"receiver_name" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	ItemDescription string `json:"item_description" binding:"required"`
}

type ShipmentResponse struct {
	ID              string `json:"id"`
	TrackingNumber  string `json:"tracking_number"`
	UserID          string `json:"user_id"`
	SenderName      string `json:"sender_name"`
	SenderAddress   string `json:"sender_address"`
	ReceiverName    string `json:"receiver_name"`
	ReceiverAddress string `json:"receiver_address"`
	ItemDescription string `json:"item_description"`
	Status          string `json:"status"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=Shipped InTransit Delivered Cancelled"`
}

type TrackRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
}
