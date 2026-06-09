package integration

import (
	"testing"

	"github.com/mhshahzad/portman/internal/ports"
)

func TestAggregate(t *testing.T) {
	dockerPorts := []ports.Port{
		{Number: 3000, Protocol: "tcp", Service: "auth-api", Source: "docker"},
	}
	systemPorts := []ports.Port{
		{Number: 3000, Protocol: "tcp", ProcessName: "node", PID: 1234, Source: "lsof"},
		{Number: 5000, Protocol: "tcp", ProcessName: "app", PID: 5678, Source: "lsof"},
	}

	agg := &Aggregator{}
	result := agg.Aggregate(dockerPorts, systemPorts)

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	// First entry should be port 3000, resolved to docker
	if result[0].Port != 3000 || result[0].Service != "auth-api" || result[0].Backing != "docker" {
		t.Errorf("port 3000 resolution failed: %+v", result[0])
	}

	// Second entry should be port 5000, resolved to system
	if result[1].Port != 5000 || result[1].Service != "app" || result[1].Backing != "system" {
		t.Errorf("port 5000 resolution failed: %+v", result[1])
	}
}
