package notification

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type WhatsAppService struct {
	parent *NotificationService
	client *resty.Client
}

func NewWhatsAppService(parent *NotificationService) *WhatsAppService {
	client := resty.New()
	client.SetHeader("Content-Type", "application/json")

	return &WhatsAppService{
		parent: parent,
		client: client,
	}
}

type WhatsAppMessage struct {
	Phone   string
	Message string
}

type WARequest struct {
	ApiKey  string `json:"api_key"`
	Sender  string `json:"sender"`
	Number  string `json:"number"`
	Message string `json:"message"`
}

type WAResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *WhatsAppService) SendMessage(data WhatsAppMessage) error {
	runtimeCfg := s.parent.whatsappConfig()

	if runtimeCfg.APIKey == "" || runtimeCfg.URL == "" {
		return fmt.Errorf("WhatsApp API configuration not set")
	}

	// Format phone number (ensure it starts with country code)
	phone := data.Phone
	if len(phone) > 0 && phone[0:1] == "0" {
		phone = "62" + phone[1:] // Convert 08xxx to 628xxx for Indonesia
	}

	payload := WARequest{
		ApiKey:  runtimeCfg.APIKey,
		Sender:  runtimeCfg.Sender,
		Number:  phone,
		Message: data.Message,
	}

	var response WAResponse
	resp, err := s.client.R().
		SetBody(payload).
		SetResult(&response).
		Post(runtimeCfg.URL)

	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("WhatsApp API error (HTTP %d): %s", resp.StatusCode(), response.Message)
	}

	if !response.Status {
		return fmt.Errorf("failed to send message: %s", response.Message)
	}

	return nil
}

// WhatsApp Message Templates
func (s *WhatsAppService) SendPOApprovalRequest(phone, poNumber, requesterName string, amount float64) error {
	message := fmt.Sprintf(`
*🔔 PO Approval Required*

PO Number: *%s*
Requested By: %s
Total Amount: Rp %.2f

Please login to the system to review and approve this PO.

_TatApps Notification_
	`, poNumber, requesterName, amount)

	return s.SendMessage(WhatsAppMessage{
		Phone:   phone,
		Message: message,
	})
}

func (s *WhatsAppService) SendLowStockAlert(phone, itemName, warehouseName string, currentQty, minQty float64) error {
	message := fmt.Sprintf(`
*⚠️ Low Stock Alert*

Item: *%s*
Warehouse: %s
Current Stock: %.2f
Minimum Stock: %.2f

Please reorder this item as soon as possible.

_TatApps Notification_
	`, itemName, warehouseName, currentQty, minQty)

	return s.SendMessage(WhatsAppMessage{
		Phone:   phone,
		Message: message,
	})
}

func (s *WhatsAppService) SendLeadFollowUpReminder(phone, leadName, contactPerson string) error {
	message := fmt.Sprintf(`
*📋 Lead Follow-up Reminder*

Lead: *%s*
Contact Person: %s

Don't forget to follow up with this lead today!

_TatApps Notification_
	`, leadName, contactPerson)

	return s.SendMessage(WhatsAppMessage{
		Phone:   phone,
		Message: message,
	})
}

func (s *WhatsAppService) SendProjectUpdate(phone, projectName, status string, progress int) error {
	message := fmt.Sprintf(`
*🚀 Project Update*

Project: *%s*
Status: %s
Progress: %d%%

Check the system for more details.

_TatApps Notification_
	`, projectName, status, progress)

	return s.SendMessage(WhatsAppMessage{
		Phone:   phone,
		Message: message,
	})
}
