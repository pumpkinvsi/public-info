package metrics
 
import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)
 
const (
	ReasonParseError        = "parse_error"
	ReasonSanitizationError = "sanitization_error"
	ReasonSMTPError         = "smtp_error"
)
 
// Metrics holds all Prometheus instruments for the service.
type Metrics struct {
	MessagesReceived    prometheus.Counter
	MessagesSent        prometheus.Counter
	MessagesFailed      *prometheus.CounterVec
	ProcessingDuration  prometheus.Histogram
	ActiveWorkers       prometheus.Gauge
	KafkaErrors         prometheus.Counter
}
 
// New registers all metrics with the default Prometheus registry and returns them.
func New() *Metrics {
	return &Metrics{
		MessagesReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "email_sender_messages_received_total",
			Help: "Total number of email messages consumed from Kafka.",
		}),
 
		MessagesSent: promauto.NewCounter(prometheus.CounterOpts{
			Name: "email_sender_messages_sent_total",
			Help: "Total number of email messages successfully sent via SMTP.",
		}),
 
		MessagesFailed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "email_sender_messages_failed_total",
			Help: "Total number of email messages that could not be delivered, by reason.",
		}, []string{"reason"}),
 
		ProcessingDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "email_sender_processing_duration_seconds",
			Help:    "End-to-end time from message receive to SMTP response, in seconds.",
			Buckets: prometheus.DefBuckets,
		}),
 
		ActiveWorkers: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "email_sender_active_workers",
			Help: "Number of worker goroutines currently processing a message.",
		}),
 
		KafkaErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "email_sender_kafka_errors_total",
			Help: "Total number of non-fatal Kafka consumer errors.",
		}),
	}
}