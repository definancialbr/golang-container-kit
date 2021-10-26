package healthcheck

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/heptiolabs/healthcheck"
)

type ProbeServiceOption func(*ProbeService)

type ProbeService struct {
	health healthcheck.Handler
}

func WithCheckForLiveness(name string, check func() error) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, check)
	}
}

func WithCheckForReadiness(name string, check func() error) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, check)
	}
}

func WithGoroutineCountCheckForLiveness(name string, threshold int) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.GoroutineCountCheck(threshold))
	}
}

func WithGoroutineCountCheckForReadiness(name string, threshold int) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.GoroutineCountCheck(threshold))
	}
}

func WithHTTPGetCheckForLiveness(name string, url string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.HTTPGetCheck(url, timeout))
	}
}

func WithHTTPGetCheckForReadiness(name string, url string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.HTTPGetCheck(url, timeout))
	}
}

func WithDNSResolveCheckForLiveness(name string, host string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.DNSResolveCheck(host, timeout))
	}
}

func WithDNSResolveCheckForReadiness(name string, host string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.DNSResolveCheck(host, timeout))
	}
}

func WithTCPDialCheckForLiveness(name string, addr string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.TCPDialCheck(addr, timeout))
	}
}

func WithTCPDialCheckForReadiness(name string, addr string, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.TCPDialCheck(addr, timeout))
	}
}

func WithDatabasePingCheckForLiveness(name string, database *sql.DB, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.DatabasePingCheck(database, timeout))
	}
}

func WithDatabasePingCheckForReadiness(name string, database *sql.DB, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.DatabasePingCheck(database, timeout))
	}
}

func WithTimeoutCheckForLiveness(name string, check healthcheck.Check, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.Timeout(check, timeout))
	}
}

func WithTimeoutCheckForReadiness(name string, check healthcheck.Check, timeout time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.Timeout(check, timeout))
	}
}

func WithAsyncCheckForLiveness(name string, check healthcheck.Check, interval time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.Async(check, interval))
	}
}

func WithAsyncCheckForReadiness(name string, check healthcheck.Check, interval time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.Async(check, interval))
	}
}

func WithAsyncWithContextCheckForLiveness(name string, ctx context.Context, check healthcheck.Check, interval time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddLivenessCheck(name, healthcheck.AsyncWithContext(ctx, check, interval))
	}
}

func WithAsyncWithContextCheckForReadiness(name string, ctx context.Context, check healthcheck.Check, interval time.Duration) ProbeServiceOption {
	return func(p *ProbeService) {
		p.health.AddReadinessCheck(name, healthcheck.AsyncWithContext(ctx, check, interval))
	}
}

func NewProbeService(options ...ProbeServiceOption) *ProbeService {

	p := &ProbeService{
		health: healthcheck.NewHandler(),
	}

	for _, option := range options {
		option(p)
	}

	return p

}

func (p *ProbeService) LivenessHandler() http.HandlerFunc {
	return p.health.LiveEndpoint
}

func (p *ProbeService) ReadinessHandler() http.HandlerFunc {
	return p.health.ReadyEndpoint
}
