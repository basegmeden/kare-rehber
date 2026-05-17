package services

import (
	"fmt"
	"log"

	"kare-rehber/backend/internal/database"
	"kare-rehber/backend/internal/models"
)

// SMSProvider is the interface for any SMS gateway.
type SMSProvider interface {
	Send(phone, message string) error
}

// ConsoleSMSProvider logs the SMS to stdout — used until a real provider is configured.
type ConsoleSMSProvider struct{}

func (p *ConsoleSMSProvider) Send(phone, message string) error {
	log.Printf("[SMS STUB] To: %s | Message: %s", phone, message)
	return nil
}

var defaultProvider SMSProvider = &ConsoleSMSProvider{}

// SetProvider replaces the active SMS provider (call this in main when configuring a real gateway).
func SetSMSProvider(p SMSProvider) {
	defaultProvider = p
}

func SendSMS(recipientID *uint, phone, message string, smsType models.SMSType) error {
	err := defaultProvider.Send(phone, message)
	status := models.SMSStatusSent
	if err != nil {
		status = models.SMSStatusFailed
		log.Printf("SMS send error: %v", err)
	}

	record := models.SMSLog{
		RecipientID: recipientID,
		Phone:       phone,
		Message:     message,
		Type:        smsType,
		Status:      status,
	}
	database.DB.Create(&record)

	if err != nil {
		return fmt.Errorf("sms send failed: %w", err)
	}
	return nil
}
