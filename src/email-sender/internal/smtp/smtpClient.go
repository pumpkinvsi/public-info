package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"src/email-sender/internal/config"
)

type SmtpClient struct {
	cfg *config.SmtpConfig
}

func NewSmtpClient(cfg *config.SmtpConfig) *SmtpClient {
	return &SmtpClient{cfg: cfg}
}

func (c *SmtpClient) Send(msg *Message) error {
	server := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	raw := buildRawMessage(c.cfg.User, msg)

	return c.sendSSL(server, raw, msg.To)
}

func (c *SmtpClient) sendSSL(server string, raw []byte, recipient string) error {
	tlsCfg := &tls.Config{
		ServerName: c.cfg.Host,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", server, tlsCfg)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}

	host, _, err := net.SplitHostPort(server)
	if err != nil {
		return fmt.Errorf("split host port: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}
	defer client.Quit() //nolint:errcheck

	auth := smtp.PlainAuth("", c.cfg.User, c.cfg.Password, c.cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(c.cfg.User); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}

	if err := client.Rcpt(recipient); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	defer wc.Close()

	if _, err := wc.Write(raw); err != nil {
		return fmt.Errorf("smtp write body: %w", err)
	}

	return nil
}

func buildRawMessage(from string, msg *Message) []byte {
	var sb strings.Builder
	writeHeader := func(key, value string) {
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("\r\n")
	}

	writeHeader("From", from)
	writeHeader("To", msg.To)
	writeHeader("Subject", "Email from web page")
	writeHeader("MIME-Version", "1.0")
	writeHeader("Content-Type", "text/html; charset=UTF-8")
	sb.WriteString("\r\n")

	body := strings.ReplaceAll(msg.Body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\n", "\r\n")
	sb.WriteString(body)

	return []byte(sb.String())
}
