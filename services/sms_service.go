package services

import (
	"fmt"
	"log"

	"github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"

	"flash-oauth2/config"
)

// SMSService defines the interface for SMS operations
type SMSService interface {
	SendVerificationCode(phone, code string) error
}

// MockSMSService is used for testing and development
type MockSMSService struct {
	// LastCode stores the last sent verification code for testing
	LastCode map[string]string
	// TestMode enables returning predictable verification codes for testing
	TestMode bool
	// TestCodes maps phone numbers to predefined verification codes for testing
	TestCodes map[string]string
}

// NewMockSMSService creates a new mock SMS service
func NewMockSMSService() *MockSMSService {
	return &MockSMSService{
		LastCode: make(map[string]string),
		TestMode: true, // Default to test mode
		TestCodes: map[string]string{
			"13800138000": "123456", // Standard test user
			"13800138001": "654321", // Premium test user
			"13800138002": "111111", // Inactive test user
			"13800138003": "999999", // New test user
			"admin":       "123456", // Admin login
		},
	}
}

// SendVerificationCode simulates sending SMS for testing
func (m *MockSMSService) SendVerificationCode(phone, code string) error {
	// In test mode, use predefined verification codes instead of the generated ones
	if m.TestMode {
		if testCode, exists := m.TestCodes[phone]; exists {
			log.Printf("MOCK SMS: Using predefined verification code %s for %s", testCode, phone)
			m.LastCode[phone] = testCode
			return nil
		}
	}

	// Fallback to using the provided code for phones not in test data
	log.Printf("MOCK SMS: Sending verification code %s to %s", code, phone)
	m.LastCode[phone] = code
	return nil
}

// GetLastCode returns the last verification code sent to a phone number (for testing)
func (m *MockSMSService) GetLastCode(phone string) string {
	return m.LastCode[phone]
}

// SetTestMode enables or disables test mode
func (m *MockSMSService) SetTestMode(enabled bool) {
	m.TestMode = enabled
}

// SetTestCode sets a predefined verification code for a specific phone number
func (m *MockSMSService) SetTestCode(phone, code string) {
	if m.TestCodes == nil {
		m.TestCodes = make(map[string]string)
	}
	m.TestCodes[phone] = code
}

// ClearTestCodes removes all predefined test codes
func (m *MockSMSService) ClearTestCodes() {
	m.TestCodes = make(map[string]string)
}

// AlibabaSMSService implements real SMS sending using Alibaba Cloud
type AlibabaSMSService struct {
	client     *dysmsapi20170525.Client
	signName   string
	templateId string
}

// NewAlibabaSMSService creates a new Alibaba Cloud SMS service
func NewAlibabaSMSService(cfg *config.SMSConfig) (*AlibabaSMSService, error) {
	config := &client.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}

	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMS client: %v", err)
	}

	return &AlibabaSMSService{
		client:     client,
		signName:   cfg.SignName,
		templateId: cfg.TemplateCode,
	}, nil
}

// SendVerificationCode sends SMS verification code using Alibaba Cloud
func (a *AlibabaSMSService) SendVerificationCode(phone, code string) error {
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(a.signName),
		TemplateCode:  tea.String(a.templateId),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}

	runtime := &util.RuntimeOptions{}

	resp, err := a.client.SendSmsWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	if resp.Body.Code != nil && *resp.Body.Code != "OK" {
		return fmt.Errorf("SMS send failed: %s - %s", *resp.Body.Code, *resp.Body.Message)
	}

	log.Printf("SMS sent successfully to %s, request ID: %s", phone, *resp.Body.RequestId)
	return nil
}

// NewSMSService creates an appropriate SMS service based on configuration
func NewSMSService(cfg *config.Config) SMSService {
	if !cfg.SMS.Enabled {
		log.Println("SMS service disabled, using mock SMS service")
		return NewMockSMSService()
	}

	if cfg.SMS.AccessKeyId == "" || cfg.SMS.AccessKeySecret == "" {
		log.Println("SMS credentials not configured, using mock SMS service")
		return NewMockSMSService()
	}

	smsService, err := NewAlibabaSMSService(cfg.SMS)
	if err != nil {
		log.Printf("Failed to create Alibaba SMS service: %v, falling back to mock", err)
		return NewMockSMSService()
	}

	log.Println("Using Alibaba Cloud SMS service")
	return smsService
}
