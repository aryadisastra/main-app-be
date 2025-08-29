package models

import "time"

type Shipment struct {
	ID              string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TrackingNumber  string    `gorm:"uniqueIndex;not null" json:"tracking_number"`
	UserID          string    `gorm:"type:uuid;not null" json:"user_id"`
	SenderName      string    `json:"sender_name"`
	SenderAddress   string    `json:"sender_address"`
	ReceiverName    string    `json:"receiver_name"`
	ReceiverAddress string    `json:"receiver_address"`
	ItemDescription string    `json:"item_description"`
	Status          string    `json:"status"` // Created, Shipped, InTransit, Delivered, Cancelled
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Shipment) TableName() string {
	return "tb_r_shipments"
}
