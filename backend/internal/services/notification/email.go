package notification

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	parent *NotificationService
}

func NewEmailService(parent *NotificationService) *EmailService {
	return &EmailService{
		parent: parent,
	}
}

type EmailData struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

func (s *EmailService) SendEmail(data EmailData) error {
	runtimeCfg := s.parent.emailConfig()
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", runtimeCfg.FromName, runtimeCfg.FromEmail))
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	if data.IsHTML {
		m.SetBody("text/html", data.Body)
	} else {
		m.SetBody("text/plain", data.Body)
	}

	d := gomail.NewDialer(
		runtimeCfg.Host,
		runtimeCfg.Port,
		runtimeCfg.Username,
		runtimeCfg.Password,
	)

	// Skip TLS verification for development
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Email Templates
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome to TatApps!</h2>
			<p>Hi %s,</p>
			<p>Your account has been created successfully.</p>
			<p>You can now login to the system and start using TatApps.</p>
			<br>
			<p>Best regards,<br>TatApps Team</p>
		</body>
		</html>
	`, name)

	return s.SendEmail(EmailData{
		To:      to,
		Subject: "Welcome to TatApps",
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailService) SendPOApprovalRequest(to, poNumber, requesterName string, amount float64) error {
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Purchase Order Approval Request</h2>
			<p>A new purchase order requires your approval:</p>
			<ul>
				<li><strong>PO Number:</strong> %s</li>
				<li><strong>Requested By:</strong> %s</li>
				<li><strong>Total Amount:</strong> Rp %.2f</li>
			</ul>
			<p>Please login to the system to review and approve this PO.</p>
			<br>
			<p>Best regards,<br>TatApps System</p>
		</body>
		</html>
	`, poNumber, requesterName, amount)

	return s.SendEmail(EmailData{
		To:      to,
		Subject: fmt.Sprintf("PO Approval Required: %s", poNumber),
		Body:    body,
		IsHTML:  true,
	})
}

func (s *EmailService) SendLowStockAlert(to, itemName, warehouseName string, currentQty, minQty float64) error {
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>⚠️ Low Stock Alert</h2>
			<p>The following item is running low on stock:</p>
			<ul>
				<li><strong>Item:</strong> %s</li>
				<li><strong>Warehouse:</strong> %s</li>
				<li><strong>Current Stock:</strong> %.2f</li>
				<li><strong>Minimum Stock:</strong> %.2f</li>
			</ul>
			<p>Please reorder this item as soon as possible.</p>
			<br>
			<p>Best regards,<br>TatApps System</p>
		</body>
		</html>
	`, itemName, warehouseName, currentQty, minQty)

	return s.SendEmail(EmailData{
		To:      to,
		Subject: fmt.Sprintf("Low Stock Alert: %s", itemName),
		Body:    body,
		IsHTML:  true,
	})
}
