package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// LoadSMTPConfigFromEnv loads SMTP configuration from environment variables
func LoadSMTPConfigFromEnv() *SMTPConfig {
	return &SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

// Validate checks if all required SMTP config fields are set
func (c *SMTPConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.Port == "" {
		return fmt.Errorf("SMTP_PORT is required")
	}
	if c.Username == "" {
		return fmt.Errorf("SMTP_USERNAME is required")
	}
	if c.Password == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}
	if c.From == "" {
		return fmt.Errorf("SMTP_FROM is required")
	}
	return nil
}

// Address returns the SMTP server address (host:port)
func (c *SMTPConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Auth returns smtp.Auth for authentication
func (c *SMTPConfig) Auth() smtp.Auth {
	return smtp.PlainAuth("", c.Username, c.Password, c.Host)
}
