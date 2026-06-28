package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	registrations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "uno_reverse_registrations_total",
		Help: "Total number of user registrations",
	})

	emailVerifications = promauto.NewCounter(prometheus.CounterOpts{
		Name: "uno_reverse_email_verifications_total",
		Help: "Total number of email verifications",
	})

	loginAttempts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "uno_reverse_login_attempts_total",
		Help: "Total login attempts",
	}, []string{"status"})

	verifiedUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "uno_reverse_verified_users_total",
		Help: "Total number of verified users",
	})

	unverifiedUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "uno_reverse_unverified_users_total",
		Help: "Total number of unverified users",
	})

	activeUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "uno_reverse_active_users_total",
		Help: "Total number of active users",
	})

	auditLogEntries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "uno_reverse_audit_log_entries_total",
		Help: "Total audit log entries",
	}, []string{"action"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "uno_reverse_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint", "status"})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "uno_reverse_http_requests_total",
		Help: "Total HTTP requests",
	}, []string{"method", "endpoint", "status"})
)

func incrementRegistrations() {
	registrations.Inc()
}

func incrementEmailVerifications() {
	emailVerifications.Inc()
}

func recordLoginAttempt(status string) {
	loginAttempts.WithLabelValues(status).Inc()
}

func updateVerifiedUsers(count float64) {
	verifiedUsers.Set(count)
}

func updateUnverifiedUsers(count float64) {
	unverifiedUsers.Set(count)
}

func recordAuditLog(action string) {
	auditLogEntries.WithLabelValues(action).Inc()
}

func recordRequestDuration(method, endpoint, status string, duration float64) {
	requestDuration.WithLabelValues(method, endpoint, status).Observe(duration)
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}
