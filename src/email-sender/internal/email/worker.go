package email

import (
	"log/slog"
	"sync"
	"time"

	"src/email-sender/internal/metrics"
	"src/email-sender/internal/smtp"
)

const (
	MaxWorkers = 5
	recipient  = "pumpkinvsi@live.ru" // TODO: make configurable
)

type Sanitizer interface {
	Sanitize(body string) string
}

type SmtpClient interface {
	Send(msg *smtp.Message) error
}

type Pool struct {
	jobs       <-chan *Email
	sanitizer  Sanitizer
	metrics    *metrics.Metrics
	wg         sync.WaitGroup
	smtpClient SmtpClient
}

func NewPool(
	jobs <-chan *Email,
	smtpClient SmtpClient,
	m *metrics.Metrics,
) *Pool {
	return &Pool{
		jobs:      jobs,
		sanitizer: &sanitizer{},
		metrics:   m,
	}
}

func (p *Pool) Start() {
	for i := range MaxWorkers {
		p.wg.Add(1)
		go p.work(i)
	}
	slog.Info("worker pool started", "workers", MaxWorkers)
}

func (p *Pool) Wait() {
	p.wg.Wait()
	slog.Info("worker pool stopped")
}

func (p *Pool) work(id int) {
	defer p.wg.Done()
	slog.Debug("worker started", "worker_id", id)

	for msg := range p.jobs {
		p.process(id, msg)
	}

	slog.Debug("worker stopped", "worker_id", id)
}

func (p *Pool) process(workerID int, msg *Email) {
	start := time.Now()

	p.metrics.ActiveWorkers.Inc()
	defer p.metrics.ActiveWorkers.Dec()

	log := slog.With(
		"worker_id", workerID,
		"sender", msg.Sender,
		"contact", msg.Contact,
	)

	msg.Text = p.sanitizer.Sanitize(msg.Text)

	if err := p.smtpClient.Send(&smtp.Message{
		From: msg.Sender,
		To:   recipient,
		Body: msg.Text,
	}); err != nil {
		log.Error("smtp delivery failed", "error", err)
		p.metrics.MessagesFailed.WithLabelValues(metrics.ReasonSMTPError).Inc()
		return
	}

	elapsed := time.Since(start).Seconds()
	p.metrics.MessagesSent.Inc()
	p.metrics.ProcessingDuration.Observe(elapsed)
	log.Info("email delivered", "duration_seconds", elapsed)
}
