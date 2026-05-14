package healthcheck

import (
	"errors"
	"net"
	"testing"
	"time"
)

func stubProbeOK(host string, timeout time.Duration) (time.Duration, error) {
	return 5 * time.Millisecond, nil
}

func stubProbeErr(host string, timeout time.Duration) (time.Duration, error) {
	return 0, errors.New("connection refused")
}

func TestNew_SetsHosts(t *testing.T) {
	hosts := []string{"host1", "host2"}
	c := New(hosts, 5*time.Second)
	if len(c.hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(c.hosts))
	}
}

func TestCheck_HealthyHost(t *testing.T) {
	c := New([]string{"host1"}, 5*time.Second)
	c.probe = stubProbeOK
	results := c.Check()
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].Status != StatusHealthy {
		t.Errorf("expected healthy, got %s", results[0].Status)
	}
	if results[0].Latency == 0 {
		t.Errorf("expected non-zero latency")
	}
}

func TestCheck_UnreachableHost(t *testing.T) {
	c := New([]string{"bad-host"}, 5*time.Second)
	c.probe = stubProbeErr
	results := c.Check()
	if results[0].Status != StatusUnreachable {
		t.Errorf("expected unreachable, got %s", results[0].Status)
	}
	if results[0].Error == "" {
		t.Errorf("expected error message")
	}
}

func TestHealthy_AllHealthy(t *testing.T) {
	results := []Result{
		{Status: StatusHealthy},
		{Status: StatusHealthy},
	}
	if !Healthy(results) {
		t.Error("expected Healthy to return true")
	}
}

func TestHealthy_OneUnreachable(t *testing.T) {
	results := []Result{
		{Status: StatusHealthy},
		{Status: StatusUnreachable},
	}
	if Healthy(results) {
		t.Error("expected Healthy to return false")
	}
}

func TestStatus_String(t *testing.T) {
	cases := map[Status]string{
		StatusHealthy:     "healthy",
		StatusUnreachable: "unreachable",
		StatusDegraded:    "degraded",
		StatusUnknown:     "unknown",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Status(%d).String() = %q, want %q", s, got, want)
		}
	}
}

func TestDialTCP_InvalidAddr(t *testing.T) {
	_, err := dialTCP("127.0.0.1:1", 100*time.Millisecond)
	if err == nil {
		t.Skip("port 1 unexpectedly open")
	}
	var netErr net.Error
	if !errors.As(err, &netErr) {
		t.Errorf("expected net.Error, got %T", err)
	}
}
