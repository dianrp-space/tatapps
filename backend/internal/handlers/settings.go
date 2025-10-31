package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tatapps/internal/config"
	"tatapps/internal/models"
	"tatapps/internal/services/notification"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	db    *gorm.DB
	notif *notification.NotificationService
	cfg   *config.Config
}

func NewSettingsHandler(db *gorm.DB, notif *notification.NotificationService, cfg *config.Config) *SettingsHandler {
	return &SettingsHandler{
		db:    db,
		notif: notif,
		cfg:   cfg,
	}
}

type NotificationSettings struct {
	Enabled         bool   `json:"enabled"`
	Threshold       int    `json:"threshold"`
	CheckFrequency  string `json:"check_frequency"`
	ScheduleMode    string `json:"schedule_mode"`
	CronExpression  string `json:"cron_expression"`
	TimeZone        string `json:"timezone"`
	WhatsAppEnabled bool   `json:"whatsapp_enabled"`
	WhatsAppNumber  string `json:"whatsapp_number"`
	EmailEnabled    bool   `json:"email_enabled"`
	EmailAddress    string `json:"email_address"`
}

type SiteSettingsResponse struct {
	AppName string `json:"app_name"`
	Logo    string `json:"logo"`
	Favicon string `json:"favicon"`
}

type SiteSettingsAdminResponse struct {
	AppName string `json:"app_name"`
	Logo    string `json:"logo"`
	Favicon string `json:"favicon"`

	WhatsAppAPIURL string `json:"whatsapp_api_url"`
	WhatsAppAPIKey string `json:"whatsapp_api_key"`
	WhatsAppSender string `json:"whatsapp_sender"`

	SMTPHost      string `json:"smtp_host"`
	SMTPPort      int    `json:"smtp_port"`
	SMTPUsername  string `json:"smtp_username"`
	SMTPPassword  string `json:"smtp_password"`
	SMTPFromEmail string `json:"smtp_from_email"`
	SMTPFromName  string `json:"smtp_from_name"`
}

func (h *SettingsHandler) getOrCreateSiteSetting() (*models.SiteSetting, error) {
	var setting models.SiteSetting
	if err := h.db.First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			setting = models.SiteSetting{
				AppName:        h.cfg.AppName,
				WhatsAppAPIURL: h.cfg.WAApiURL,
				WhatsAppAPIKey: h.cfg.WAApiKey,
				WhatsAppSender: h.cfg.WASender,
				SMTPHost:       h.cfg.SMTPHost,
				SMTPPort:       h.cfg.SMTPPort,
				SMTPUsername:   h.cfg.SMTPUsername,
				SMTPPassword:   h.cfg.SMTPPassword,
				SMTPFromEmail:  h.cfg.SMTPFromEmail,
				SMTPFromName:   h.cfg.SMTPFromName,
			}
			if err := h.db.Create(&setting).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &setting, nil
}

func (h *SettingsHandler) GetSiteSettings(c *gin.Context) {
	setting, err := h.getOrCreateSiteSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load site settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": SiteSettingsResponse{
			AppName: setting.AppName,
			Logo:    setting.LogoPath,
			Favicon: setting.FaviconPath,
		},
	})
}

func (h *SettingsHandler) GetSiteSettingsAdmin(c *gin.Context) {
	setting, err := h.getOrCreateSiteSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load site settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": SiteSettingsAdminResponse{
			AppName: setting.AppName,
			Logo:    setting.LogoPath,
			Favicon: setting.FaviconPath,

			WhatsAppAPIURL: setting.WhatsAppAPIURL,
			WhatsAppAPIKey: setting.WhatsAppAPIKey,
			WhatsAppSender: setting.WhatsAppSender,

			SMTPHost:      setting.SMTPHost,
			SMTPPort:      setting.SMTPPort,
			SMTPUsername:  setting.SMTPUsername,
			SMTPPassword:  setting.SMTPPassword,
			SMTPFromEmail: setting.SMTPFromEmail,
			SMTPFromName:  setting.SMTPFromName,
		},
	})
}

func saveSiteAsset(c *gin.Context, file *multipart.FileHeader, prefix string) (string, error) {
	destDir := filepath.Join("uploads", "site")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".png"
	}

	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	fullPath := filepath.Join(destDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}

func deleteIfExists(path string) {
	if path == "" {
		return
	}
	if _, err := os.Stat(path); err == nil {
		_ = os.Remove(path)
	}
}

func (h *SettingsHandler) UpdateSiteSettings(c *gin.Context) {
	setting, err := h.getOrCreateSiteSetting()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load site settings"})
		return
	}

	if appName, exists := c.GetPostForm("app_name"); exists {
		setting.AppName = strings.TrimSpace(appName)
	}

	if value, exists := c.GetPostForm("whatsapp_api_url"); exists {
		setting.WhatsAppAPIURL = strings.TrimSpace(value)
	}
	if value, exists := c.GetPostForm("whatsapp_api_key"); exists {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			setting.WhatsAppAPIKey = trimmed
		}
	}
	if value, exists := c.GetPostForm("whatsapp_sender"); exists {
		setting.WhatsAppSender = strings.TrimSpace(value)
	}

	if value, exists := c.GetPostForm("smtp_host"); exists {
		setting.SMTPHost = strings.TrimSpace(value)
	}
	if value, exists := c.GetPostForm("smtp_port"); exists {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			setting.SMTPPort = 0
		} else {
			if port, err := strconv.Atoi(trimmed); err == nil {
				setting.SMTPPort = port
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SMTP port"})
				return
			}
		}
	}
	if value, exists := c.GetPostForm("smtp_username"); exists {
		setting.SMTPUsername = strings.TrimSpace(value)
	}
	if value, exists := c.GetPostForm("smtp_password"); exists {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			setting.SMTPPassword = trimmed
		}
	}
	if value, exists := c.GetPostForm("smtp_from_email"); exists {
		setting.SMTPFromEmail = strings.TrimSpace(value)
	}
	if value, exists := c.GetPostForm("smtp_from_name"); exists {
		setting.SMTPFromName = strings.TrimSpace(value)
	}

	if logoFile, err := c.FormFile("logo"); err == nil {
		newPath, saveErr := saveSiteAsset(c, logoFile, "logo")
		if saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload logo"})
			return
		}
		deleteIfExists(setting.LogoPath)
		setting.LogoPath = newPath
	} else if !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid logo upload"})
		return
	}

	if faviconFile, err := c.FormFile("favicon"); err == nil {
		newPath, saveErr := saveSiteAsset(c, faviconFile, "favicon")
		if saveErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload favicon"})
			return
		}
		deleteIfExists(setting.FaviconPath)
		setting.FaviconPath = newPath
	} else if !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid favicon upload"})
		return
	}

	if err := h.db.Save(setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update site settings"})
		return
	}

	h.notif.UpdateFromSiteSetting(setting)

	c.JSON(http.StatusOK, gin.H{
		"message": "Site settings updated successfully",
		"data": SiteSettingsResponse{
			AppName: setting.AppName,
			Logo:    setting.LogoPath,
			Favicon: setting.FaviconPath,
		},
	})
}

func (h *SettingsHandler) BackupDatabase(c *gin.Context) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("tatapps_backup_%s.sql", timestamp)

	args := []string{
		"--clean",
		"--if-exists",
		"--no-owner",
		"--no-privileges",
		"--host", h.cfg.DBHost,
		"--port", h.cfg.DBPort,
		"--username", h.cfg.DBUser,
		"--dbname", h.cfg.DBName,
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", h.cfg.DBPassword))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create backup: %s", message)})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/sql", output)
}

func (h *SettingsHandler) RestoreDatabase(c *gin.Context) {
	file, err := c.FormFile("backup")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup file is required"})
		return
	}

	tempFile, err := os.CreateTemp("", "tatapps-restore-*.sql")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare restore file"})
		return
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded backup"})
		return
	}

	sanitizedPath, err := h.sanitizeBackup(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(sanitizedPath)

	restoreFile, err := os.Open(sanitizedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded backup"})
		return
	}
	defer restoreFile.Close()

	args := []string{
		"--host", h.cfg.DBHost,
		"--port", h.cfg.DBPort,
		"--username", h.cfg.DBUser,
		"--dbname", h.cfg.DBName,
		"-v", "ON_ERROR_STOP=1",
		"--single-transaction",
	}

	cmd := exec.Command("psql", args...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", h.cfg.DBPassword))
	cmd.Stdin = restoreFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		log.Printf("[restore] psql failed: %s", message)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to restore database: %s", message)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database restored successfully"})
}

func quoteIdentifier(identifier string) string {
	if identifier == "" {
		return ""
	}
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}

func (h *SettingsHandler) sanitizeBackup(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Gagal membaca file backup: %w", err)
	}

	content := string(data)
	target := quoteIdentifier(h.cfg.DBUser)

	ownerPattern := regexp.MustCompile(`(?i)OWNER TO\s+("[^"]+"|\S+);`)
	content = ownerPattern.ReplaceAllString(content, fmt.Sprintf("OWNER TO %s;", target))

	grantPattern := regexp.MustCompile(`(?i)TO\s+"?postgres"?`)
	content = grantPattern.ReplaceAllString(content, "TO "+target)

	fromPattern := regexp.MustCompile(`(?i)FROM\s+"?postgres"?`)
	content = fromPattern.ReplaceAllString(content, "FROM "+target)

	rolePattern := regexp.MustCompile(`(?i)FOR ROLE\s+"?postgres"?`)
	content = rolePattern.ReplaceAllString(content, "FOR ROLE "+target)

	authPattern := regexp.MustCompile(`(?i)AUTHORIZATION\s+"?postgres"?`)
	content = authPattern.ReplaceAllString(content, "AUTHORIZATION "+target)

	setRolePattern := regexp.MustCompile(`(?i)SET ROLE\s+"?postgres"?;`)
	content = setRolePattern.ReplaceAllString(content, fmt.Sprintf("SET ROLE %s;", target))

	roleGrantPattern := regexp.MustCompile(`(?im)^\s*GRANT\s+"?postgres"?\s+TO\s+[^\n;]+;[ \t]*\n?`)
	content = roleGrantPattern.ReplaceAllString(content, "")

	roleRevokePattern := regexp.MustCompile(`(?im)^\s*REVOKE\s+[^\n;]+\s+FROM\s+"?postgres"?\s*;[ \t]*\n?`)
	content = roleRevokePattern.ReplaceAllString(content, "")

	roleAlterPattern := regexp.MustCompile(`(?im)^\s*ALTER\s+ROLE\s+"?postgres"?[^\n;]*;[ \t]*\n?`)
	content = roleAlterPattern.ReplaceAllString(content, "")

	connectPattern := regexp.MustCompile(`(?im)^\\connect\s+([^\s]+)(?:\s+[^\s]+){0,2}`)
	content = connectPattern.ReplaceAllString(content, `\connect $1`)

	dropDatabasePattern := regexp.MustCompile(`(?im)^\s*DROP\s+DATABASE\s+[^;]+;[ \t]*\n?`)
	content = dropDatabasePattern.ReplaceAllString(content, "")

	alterDatabasePattern := regexp.MustCompile(`(?im)^\s*ALTER\s+DATABASE\s+[^;]+;[ \t]*\n?`)
	content = alterDatabasePattern.ReplaceAllString(content, "")

	createDatabasePattern := regexp.MustCompile(`(?im)^\s*CREATE\s+DATABASE\s+[^;]+;[ \t]*\n?`)
	content = createDatabasePattern.ReplaceAllString(content, "")

	commentDatabasePattern := regexp.MustCompile(`(?im)^\s*COMMENT\s+ON\s+DATABASE\s+[^;]+;[ \t]*\n?`)
	content = commentDatabasePattern.ReplaceAllString(content, "")

	dbGrantPattern := regexp.MustCompile(`(?im)^\s*(GRANT|REVOKE)\s+[^\n;]+ON\s+DATABASE\s+[^\n;]+;[ \t]*\n?`)
	content = dbGrantPattern.ReplaceAllString(content, "")

	prefix := `-- tatapps restore prelude
DROP OWNED BY CURRENT_USER CASCADE;
`

	content = prefix + content

	sanitizedPath := path + ".sanitized"
	if err := os.WriteFile(sanitizedPath, []byte(content), 0600); err != nil {
		return "", fmt.Errorf("Gagal menulis file backup hasil sanitasi: %w", err)
	}
	return sanitizedPath, nil
}

// GetNotificationSettings retrieves notification settings for the current user
func (h *SettingsHandler) GetNotificationSettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var settings models.NotificationSetting
	err := h.db.Where("user_id = ?", userID).First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		// Return default settings if not found
		c.JSON(http.StatusOK, gin.H{
			"data": NotificationSettings{
				Enabled:         false,
				Threshold:       10,
				CheckFrequency:  "daily",
				ScheduleMode:    "preset",
				CronExpression:  "0 9 * * *",
				TimeZone:        "Asia/Jakarta",
				WhatsAppEnabled: false,
				WhatsAppNumber:  "",
				EmailEnabled:    false,
				EmailAddress:    "",
			},
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	}

	responseData := NotificationSettings{
		Enabled:         settings.Enabled,
		Threshold:       settings.Threshold,
		CheckFrequency:  settings.CheckFrequency,
		ScheduleMode:    settings.ScheduleMode,
		CronExpression:  settings.CronExpression,
		TimeZone:        settings.TimeZone,
		WhatsAppEnabled: settings.WhatsAppEnabled,
		WhatsAppNumber:  settings.WhatsAppNumber,
		EmailEnabled:    settings.EmailEnabled,
		EmailAddress:    settings.EmailAddress,
	}

	if strings.TrimSpace(responseData.ScheduleMode) == "" {
		responseData.ScheduleMode = "preset"
	}
	if strings.TrimSpace(responseData.TimeZone) == "" {
		responseData.TimeZone = "Asia/Jakarta"
	}
	if strings.TrimSpace(responseData.CronExpression) == "" {
		responseData.CronExpression = "0 9 * * *"
	}

	responseData.WhatsAppNumber = notification.NormalizeWhatsAppRecipients(responseData.WhatsAppNumber)

	c.JSON(http.StatusOK, gin.H{
		"data": responseData,
	})
}

// UpdateNotificationSettings updates notification settings for the current user
func (h *SettingsHandler) UpdateNotificationSettings(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req NotificationSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ScheduleMode = strings.TrimSpace(req.ScheduleMode)
	if req.ScheduleMode == "" {
		req.ScheduleMode = "preset"
	}

	req.CronExpression = strings.TrimSpace(req.CronExpression)
	req.CheckFrequency = strings.TrimSpace(req.CheckFrequency)
	req.TimeZone = strings.TrimSpace(req.TimeZone)
	if req.TimeZone == "" {
		req.TimeZone = "Asia/Jakarta"
	}

	switch req.ScheduleMode {
	case "cron":
		if req.CronExpression == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cron expression is required for custom schedule"})
			return
		}
	default:
		req.ScheduleMode = "preset"
		if strings.TrimSpace(req.CheckFrequency) == "" {
			req.CheckFrequency = "daily"
		}
	}

	req.WhatsAppNumber = notification.NormalizeWhatsAppRecipients(req.WhatsAppNumber)

	var settings models.NotificationSetting
	err := h.db.Where("user_id = ?", userID).First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		// Create new settings
		settings = models.NotificationSetting{
			UserID:          userID,
			Enabled:         req.Enabled,
			Threshold:       req.Threshold,
			CheckFrequency:  req.CheckFrequency,
			ScheduleMode:    req.ScheduleMode,
			CronExpression:  req.CronExpression,
			TimeZone:        req.TimeZone,
			WhatsAppEnabled: req.WhatsAppEnabled,
			WhatsAppNumber:  req.WhatsAppNumber,
			EmailEnabled:    req.EmailEnabled,
			EmailAddress:    req.EmailAddress,
			LastRunAt:       nil,
		}

		if err := h.db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch settings"})
		return
	} else {
		originalEnabled := settings.Enabled
		originalMode := strings.TrimSpace(settings.ScheduleMode)
		originalCron := strings.TrimSpace(settings.CronExpression)
		originalFreq := strings.TrimSpace(settings.CheckFrequency)
		originalTZ := strings.TrimSpace(settings.TimeZone)
		resetSchedule := false

		settings.Enabled = req.Enabled
		settings.Threshold = req.Threshold
		settings.CheckFrequency = req.CheckFrequency
		settings.ScheduleMode = req.ScheduleMode
		settings.CronExpression = req.CronExpression
		settings.TimeZone = req.TimeZone
		settings.WhatsAppEnabled = req.WhatsAppEnabled
		settings.WhatsAppNumber = req.WhatsAppNumber
		settings.EmailEnabled = req.EmailEnabled
		settings.EmailAddress = req.EmailAddress

		if originalEnabled != settings.Enabled {
			resetSchedule = true
		}
		if originalMode != settings.ScheduleMode ||
			originalCron != settings.CronExpression ||
			originalFreq != settings.CheckFrequency ||
			originalTZ != settings.TimeZone {
			resetSchedule = true
		}
		if resetSchedule {
			settings.LastRunAt = nil
		}

		if err := h.db.Save(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Settings updated successfully",
		"data":    settings,
	})
}

type TestNotificationRequest struct {
	WhatsAppEnabled bool   `json:"whatsapp_enabled"`
	WhatsAppNumber  string `json:"whatsapp_number"`
	EmailEnabled    bool   `json:"email_enabled"`
	EmailAddress    string `json:"email_address"`
}

// SendTestNotification sends a test notification
func (h *SettingsHandler) SendTestNotification(c *gin.Context) {
	var req TestNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	message := "This is a test notification from TatApps inventory system. If you receive this, your notification settings are working correctly."

	var errors []string
	whatsAppSent := false

	// Send WhatsApp notification
	if req.WhatsAppEnabled {
		recipients := notification.SplitWhatsAppRecipients(req.WhatsAppNumber)
		if len(recipients) == 0 {
			errors = append(errors, "WhatsApp: no valid phone numbers provided")
		} else {
			for _, phone := range recipients {
				if err := h.notif.WhatsApp.SendMessage(notification.WhatsAppMessage{
					Phone:   phone,
					Message: message,
				}); err != nil {
					errors = append(errors, fmt.Sprintf("WhatsApp (%s): %s", phone, err.Error()))
				} else {
					whatsAppSent = true
				}
			}
		}
	}

	// Send Email notification
	if req.EmailEnabled && req.EmailAddress != "" {
		if err := h.notif.Email.SendEmail(notification.EmailData{
			To:      req.EmailAddress,
			Subject: "Test Notification - TatApps",
			Body:    message,
			IsHTML:  false,
		}); err != nil {
			errors = append(errors, "Email: "+err.Error())
		}
	}

	// Log notification in history
	h.db.Create(&models.NotificationHistory{
		UserID:       userID,
		Type:         "test",
		Title:        "Test Notification",
		Message:      message,
		WhatsAppSent: whatsAppSent,
		EmailSent:    req.EmailEnabled && len(errors) == 0,
	})

	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "Some notifications failed to send",
			"errors":  errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test notification sent successfully",
	})
}

// GetNotificationHistory retrieves notification history for the current user
func (h *SettingsHandler) GetNotificationHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	limit := 10

	// Parse limit from query parameter
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var history []models.NotificationHistory
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notification history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}

type CheckLowStockRequest struct {
	SendWhatsApp bool `json:"send_whatsapp"`
	SendEmail    bool `json:"send_email"`
}

// CheckLowStock manually checks for low stock items and sends notifications
func (h *SettingsHandler) CheckLowStock(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Parse request body for channel selection
	var req CheckLowStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Default to both if not specified
		req.SendWhatsApp = true
		req.SendEmail = true
	}

	// Validate at least one channel is selected
	if !req.SendWhatsApp && !req.SendEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please select at least one notification channel",
		})
		return
	}

	// Get user data for phone and email
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validate user has the required contact info for selected channels
	if req.SendWhatsApp && user.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "WhatsApp selected but phone number not set in profile",
		})
		return
	}

	if req.SendEmail && user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email selected but email address not set in profile",
		})
		return
	}

	// Load inventory items and compute low stock using the same rules as the dashboard
	var allItems []models.InventoryItem
	if err := h.db.
		Preload("Warehouse").
		Where("is_active = ?", true).
		Find(&allItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch inventory items for low stock check",
		})
		return
	}

	lowStock := notification.ComputeLowStockEntries(allItems)

	if len(lowStock) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":         "No low stock items found",
			"low_stock_count": 0,
		})
		return
	}
	message := notification.BuildLowStockMessage(lowStock)
	lowStockCount := len(lowStock)

	var errors []string
	whatsappSent := false
	emailSent := false

	// Send WhatsApp to user's phone (from profile) - only if selected
	if req.SendWhatsApp && user.Phone != "" {
		if err := h.notif.WhatsApp.SendMessage(notification.WhatsAppMessage{
			Phone:   user.Phone,
			Message: message,
		}); err != nil {
			errors = append(errors, "WhatsApp: "+err.Error())
		} else {
			whatsappSent = true
		}
	}

	// Send Email to user's email (from profile) - only if selected
	if req.SendEmail && user.Email != "" {
		if err := h.notif.Email.SendEmail(notification.EmailData{
			To:      user.Email,
			Subject: fmt.Sprintf("Low Stock Alert - %d Items", lowStockCount),
			Body:    message,
			IsHTML:  false,
		}); err != nil {
			errors = append(errors, "Email: "+err.Error())
		} else {
			emailSent = true
		}
	}

	// Log notification in history
	h.db.Create(&models.NotificationHistory{
		UserID:       userID,
		Type:         "low_stock",
		Title:        fmt.Sprintf("Low Stock Alert - %d Items", lowStockCount),
		Message:      message,
		WhatsAppSent: whatsappSent,
		EmailSent:    emailSent,
	})

	if len(errors) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message":         "Some notifications failed to send",
			"errors":          errors,
			"low_stock_count": lowStockCount,
			"items_sent":      whatsappSent || emailSent,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Low stock notifications sent successfully",
		"low_stock_count": len(lowStock),
		"whatsapp_sent":   whatsappSent,
		"email_sent":      emailSent,
		"items":           lowStock,
	})
}
