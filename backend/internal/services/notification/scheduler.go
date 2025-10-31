package notification

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"tatapps/internal/models"

	"gorm.io/gorm"
)

type scheduledJob struct {
	setting  *models.NotificationSetting
	runTime  time.Time
	location *time.Location
}

// LowStockScheduler polls user notification preferences and dispatches low stock
// alerts according to their chosen schedule.
type LowStockScheduler struct {
	db         *gorm.DB
	notifier   *NotificationService
	parser     cron.Parser
	quit       chan struct{}
	wg         sync.WaitGroup
	locCache   map[string]*time.Location
	locCacheMu sync.RWMutex
}

// NewLowStockScheduler constructs a scheduler. Call Start to activate the loop.
func NewLowStockScheduler(db *gorm.DB, notifier *NotificationService) *LowStockScheduler {
	return &LowStockScheduler{
		db:       db,
		notifier: notifier,
		parser: cron.NewParser(
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
		),
		quit:     make(chan struct{}),
		locCache: make(map[string]*time.Location),
	}
}

// Start begins the scheduler loop.
func (s *LowStockScheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		// Run once shortly after boot.
		s.tick()

		for {
			select {
			case <-ticker.C:
				s.tick()
			case <-s.quit:
				return
			}
		}
	}()
}

// Stop stops the scheduler and waits for the loop to terminate.
func (s *LowStockScheduler) Stop(ctx context.Context) {
	close(s.quit)
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}

func (s *LowStockScheduler) tick() {
	now := time.Now().UTC()

	var settings []models.NotificationSetting
	if err := s.db.
		Where("enabled = ?", true).
		Find(&settings).Error; err != nil {
		log.Printf("[scheduler] failed to load notification settings: %v", err)
		return
	}

	if len(settings) == 0 {
		return
	}

	due := s.computeDueJobs(settings, now)
	if len(due) == 0 {
		return
	}

	var allItems []models.InventoryItem
	if err := s.db.
		Preload("Warehouse").
		Where("is_active = ?", true).
		Find(&allItems).Error; err != nil {
		log.Printf("[scheduler] failed to load inventory for low stock check: %v", err)
		return
	}

	entries := ComputeLowStockEntries(allItems)
	message := ""
	if len(entries) > 0 {
		message = BuildLowStockMessage(entries)
	}

	for _, job := range due {
		s.handleJob(job, entries, message)
	}
}

func (s *LowStockScheduler) computeDueJobs(settings []models.NotificationSetting, now time.Time) []scheduledJob {
	due := make([]scheduledJob, 0, len(settings))
	for i := range settings {
		setting := &settings[i]
		if !setting.Enabled {
			continue
		}
		if !setting.WhatsAppEnabled && !setting.EmailEnabled {
			continue
		}
		if setting.WhatsAppEnabled && NormalizeWhatsAppRecipients(setting.WhatsAppNumber) == "" && !setting.EmailEnabled {
			// No usable recipients.
			continue
		}
		if setting.EmailEnabled && strings.TrimSpace(setting.EmailAddress) == "" && !setting.WhatsAppEnabled {
			continue
		}

		loc := s.location(setting.TimeZone)
		runTime, ok := s.nextRunTime(setting, now.In(loc))
		if !ok {
			continue
		}

		due = append(due, scheduledJob{
			setting:  setting,
			runTime:  runTime,
			location: loc,
		})
	}
	return due
}

func (s *LowStockScheduler) nextRunTime(setting *models.NotificationSetting, localizedNow time.Time) (time.Time, bool) {
	lastRun := time.Time{}
	if setting.LastRunAt != nil {
		lastRun = setting.LastRunAt.In(localizedNow.Location())
	}

	mode := strings.ToLower(strings.TrimSpace(setting.ScheduleMode))
	if mode == "cron" && strings.TrimSpace(setting.CronExpression) != "" {
		return s.nextCronRun(setting, localizedNow, lastRun)
	}
	return s.nextPresetRun(setting, localizedNow, lastRun)
}

func (s *LowStockScheduler) nextCronRun(setting *models.NotificationSetting, localizedNow, lastRun time.Time) (time.Time, bool) {
	schedule, err := s.parser.Parse(strings.TrimSpace(setting.CronExpression))
	if err != nil {
		log.Printf("[scheduler] invalid cron expression for user %d: %s", setting.UserID, err)
		return time.Time{}, false
	}

	base := localizedNow.Add(-time.Minute)
	if !lastRun.IsZero() {
		base = lastRun
	}

	next := schedule.Next(base)
	if next.After(localizedNow) {
		return time.Time{}, false
	}
	if !lastRun.IsZero() && !next.After(lastRun) {
		return time.Time{}, false
	}
	return next, true
}

func (s *LowStockScheduler) nextPresetRun(setting *models.NotificationSetting, localizedNow, lastRun time.Time) (time.Time, bool) {
	switch strings.ToLower(strings.TrimSpace(setting.CheckFrequency)) {
	case "daily":
		scheduled := time.Date(localizedNow.Year(), localizedNow.Month(), localizedNow.Day(), 9, 0, 0, 0, localizedNow.Location())
		if localizedNow.Before(scheduled) {
			return time.Time{}, false
		}
		if lastRun.IsZero() || lastRun.Before(scheduled) {
			return scheduled, true
		}
	case "weekly":
		weekday := int(localizedNow.Weekday())
		offset := (weekday + 6) % 7 // days since Monday
		monday := time.Date(localizedNow.Year(), localizedNow.Month(), localizedNow.Day(), 0, 0, 0, 0, localizedNow.Location()).
			Add(-time.Duration(offset) * 24 * time.Hour)
		scheduled := monday.Add(9 * time.Hour)
		if localizedNow.Before(scheduled) {
			return time.Time{}, false
		}
		if lastRun.IsZero() || lastRun.Before(scheduled) {
			return scheduled, true
		}
	default: // hourly fallback
		if lastRun.IsZero() || localizedNow.Sub(lastRun) >= time.Hour {
			return localizedNow.Truncate(time.Minute), true
		}
	}
	return time.Time{}, false
}

func (s *LowStockScheduler) handleJob(job scheduledJob, entries []LowStockEntry, message string) {
	runTimeUTC := time.Now().UTC()

	defer s.updateLastRun(job.setting.ID, runTimeUTC)

	if len(entries) == 0 {
		return
	}

	var (
		errors       []string
		whatsAppSent bool
		emailSent    bool
	)

	if job.setting.WhatsAppEnabled {
		recipients := SplitWhatsAppRecipients(job.setting.WhatsAppNumber)
		if len(recipients) == 0 {
			errors = append(errors, "WhatsApp: no valid phone numbers configured")
		} else {
			for _, phone := range recipients {
				if err := s.notifier.WhatsApp.SendMessage(WhatsAppMessage{
					Phone:   phone,
					Message: message,
				}); err != nil {
					errors = append(errors, "WhatsApp ("+phone+"): "+err.Error())
				} else {
					whatsAppSent = true
				}
			}
		}
	}

	if job.setting.EmailEnabled && strings.TrimSpace(job.setting.EmailAddress) != "" {
		subject := "Low Stock Alert"
		if count := len(entries); count > 0 {
			subject = subject + " - " + strconv.Itoa(count) + " Items"
		}
		if err := s.notifier.Email.SendEmail(EmailData{
			To:      job.setting.EmailAddress,
			Subject: subject,
			Body:    message,
			IsHTML:  false,
		}); err != nil {
			errors = append(errors, "Email: "+err.Error())
		} else {
			emailSent = true
		}
	}

	if len(errors) > 0 {
		log.Printf("[scheduler] errors sending low stock notification for user %d: %v", job.setting.UserID, strings.Join(errors, "; "))
	}

	if !whatsAppSent && !emailSent {
		return
	}

	if err := s.db.Create(&models.NotificationHistory{
		UserID:       job.setting.UserID,
		Type:         "low_stock",
		Title:        "Low Stock Alert",
		Message:      message,
		WhatsAppSent: whatsAppSent,
		EmailSent:    emailSent,
	}).Error; err != nil {
		log.Printf("[scheduler] failed to record notification history for user %d: %v", job.setting.UserID, err)
	}
}

func (s *LowStockScheduler) updateLastRun(settingID uint, runTime time.Time) {
	if err := s.db.Model(&models.NotificationSetting{}).
		Where("id = ?", settingID).
		Updates(map[string]interface{}{"last_run_at": runTime}).Error; err != nil {
		log.Printf("[scheduler] failed to update last_run_at for setting %d: %v", settingID, err)
	}
}

func (s *LowStockScheduler) location(tz string) *time.Location {
	zone := strings.TrimSpace(tz)
	if zone == "" {
		zone = "Asia/Jakarta"
	}

	s.locCacheMu.RLock()
	loc, ok := s.locCache[zone]
	s.locCacheMu.RUnlock()
	if ok {
		return loc
	}

	loaded, err := time.LoadLocation(zone)
	if err != nil {
		log.Printf("[scheduler] invalid timezone %q, falling back to server local: %v", zone, err)
		return time.Local
	}

	s.locCacheMu.Lock()
	s.locCache[zone] = loaded
	s.locCacheMu.Unlock()
	return loaded
}
