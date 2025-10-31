package notification

import (
	"sync"
	"tatapps/internal/config"
	"tatapps/internal/models"

	"gorm.io/gorm"
)

type NotificationService struct {
	cfg         *config.Config
	db          *gorm.DB
	mu          sync.RWMutex
	siteSetting *models.SiteSetting

	Email    *EmailService
	WhatsApp *WhatsAppService
}

func NewNotificationService(cfg *config.Config, db *gorm.DB) *NotificationService {
	service := &NotificationService{
		cfg: cfg,
		db:  db,
	}
	service.Email = NewEmailService(service)
	service.WhatsApp = NewWhatsAppService(service)
	service.refreshSiteSetting()
	return service
}

func (s *NotificationService) refreshSiteSetting() {
	if s.db == nil {
		return
	}
	var setting models.SiteSetting
	if err := s.db.First(&setting).Error; err == nil {
		s.mu.Lock()
		s.siteSetting = &setting
		s.mu.Unlock()
	}
}

func (s *NotificationService) getSiteSetting() *models.SiteSetting {
	s.mu.RLock()
	setting := s.siteSetting
	s.mu.RUnlock()
	if setting != nil {
		return setting
	}

	if s.db == nil {
		return nil
	}

	var fresh models.SiteSetting
	if err := s.db.First(&fresh).Error; err == nil {
		s.mu.Lock()
		s.siteSetting = &fresh
		s.mu.Unlock()
		return &fresh
	}
	return nil
}

func (s *NotificationService) UpdateFromSiteSetting(setting *models.SiteSetting) {
	if setting == nil {
		return
	}
	s.mu.Lock()
	copied := *setting
	s.siteSetting = &copied
	s.mu.Unlock()
}

type emailRuntimeConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

func (s *NotificationService) emailConfig() emailRuntimeConfig {
	cfg := emailRuntimeConfig{
		Host:      s.cfg.SMTPHost,
		Port:      s.cfg.SMTPPort,
		Username:  s.cfg.SMTPUsername,
		Password:  s.cfg.SMTPPassword,
		FromEmail: s.cfg.SMTPFromEmail,
		FromName:  s.cfg.SMTPFromName,
	}

	if setting := s.getSiteSetting(); setting != nil {
		if setting.SMTPHost != "" {
			cfg.Host = setting.SMTPHost
		}
		if setting.SMTPPort != 0 {
			cfg.Port = setting.SMTPPort
		}
		if setting.SMTPUsername != "" {
			cfg.Username = setting.SMTPUsername
		}
		if setting.SMTPPassword != "" {
			cfg.Password = setting.SMTPPassword
		}
		if setting.SMTPFromEmail != "" {
			cfg.FromEmail = setting.SMTPFromEmail
		}
		if setting.SMTPFromName != "" {
			cfg.FromName = setting.SMTPFromName
		}
	}

	return cfg
}

type whatsappRuntimeConfig struct {
	URL    string
	APIKey string
	Sender string
}

func (s *NotificationService) whatsappConfig() whatsappRuntimeConfig {
	cfg := whatsappRuntimeConfig{
		URL:    s.cfg.WAApiURL,
		APIKey: s.cfg.WAApiKey,
		Sender: s.cfg.WASender,
	}

	if setting := s.getSiteSetting(); setting != nil {
		if setting.WhatsAppAPIURL != "" {
			cfg.URL = setting.WhatsAppAPIURL
		}
		if setting.WhatsAppAPIKey != "" {
			cfg.APIKey = setting.WhatsAppAPIKey
		}
		if setting.WhatsAppSender != "" {
			cfg.Sender = setting.WhatsAppSender
		}
	}

	return cfg
}

// SendBoth sends notification via both email and WhatsApp
func (s *NotificationService) SendBoth(email, phone, subject, message string) error {
	// Send email
	if email != "" {
		if err := s.Email.SendEmail(EmailData{
			To:      email,
			Subject: subject,
			Body:    message,
			IsHTML:  false,
		}); err != nil {
			// Log error but continue
		}
	}

	// Send WhatsApp
	if phone != "" {
		if err := s.WhatsApp.SendMessage(WhatsAppMessage{
			Phone:   phone,
			Message: message,
		}); err != nil {
			// Log error but continue
		}
	}

	return nil
}
