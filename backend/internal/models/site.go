package models

import "time"

// SiteSetting stores site-wide branding configuration.
type SiteSetting struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AppName     string    `gorm:"default:'TatApps'" json:"app_name"`
	LogoPath    string    `json:"logo_path"`
	FaviconPath string    `json:"favicon_path"`

	// WhatsApp API
	WhatsAppAPIURL string `json:"whatsapp_api_url"`
	WhatsAppAPIKey string `json:"whatsapp_api_key"`
	WhatsAppSender string `json:"whatsapp_sender"`

	// SMTP Email
	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	SMTPFromEmail string `json:"smtp_from_email"`
	SMTPFromName  string `json:"smtp_from_name"`
}
