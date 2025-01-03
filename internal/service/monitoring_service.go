package service

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MonitoringService struct {
	loginAttempts   prometheus.Counter
	failedLogins    prometheus.Counter
	activeUsers     prometheus.Gauge
	blockedUsers    prometheus.Gauge
	requestDuration prometheus.Histogram
	auditRepo       repository.AuditRepository
	securityRepo    repository.SecurityRepository
}

func NewMonitoringService(auditRepo repository.AuditRepository, securityRepo repository.SecurityRepository) *MonitoringService {
	ms := &MonitoringService{
		loginAttempts: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "auth_login_attempts_total",
			Help: "Total number of login attempts",
		}),
		failedLogins: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "auth_failed_logins_total",
			Help: "Total number of failed login attempts",
		}),
		activeUsers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "auth_active_users",
			Help: "Number of active users",
		}),
		blockedUsers: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "auth_blocked_users",
			Help: "Number of blocked users",
		}),
		requestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "auth_request_duration_seconds",
			Help:    "Time spent processing requests",
			Buckets: prometheus.DefBuckets,
		}),
		auditRepo:    auditRepo,
		securityRepo: securityRepo,
	}

	prometheus.MustRegister(
		ms.loginAttempts,
		ms.failedLogins,
		ms.activeUsers,
		ms.blockedUsers,
		ms.requestDuration,
	)

	return ms
}

func (s *MonitoringService) RecordLoginAttempt(success bool) {
	s.loginAttempts.Inc()
	if !success {
		s.failedLogins.Inc()
	}
}

func (s *MonitoringService) UpdateUserCounts(active, blocked int) {
	s.activeUsers.Set(float64(active))
	s.blockedUsers.Set(float64(blocked))
}

func (s *MonitoringService) RecordRequestDuration(duration time.Duration) {
	s.requestDuration.Observe(duration.Seconds())
}

func (s *MonitoringService) GetSecurityAlerts(ctx context.Context, from, to time.Time) ([]SecurityAlert, error) {
	// Güvenlik uyarılarını analiz et
	return s.securityRepo.GetAlerts(ctx, from, to)
}
