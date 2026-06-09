package scanner

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestResolveServiceName(t *testing.T) {
	s := &DockerScanner{}

	tests := []struct {
		name     string
		summary  container.Summary
		expected string
	}{
		{
			name: "Compose label priority",
			summary: container.Summary{
				Labels: map[string]string{
					"com.docker.compose.service": "web-api",
				},
				Names: []string{"/web-api-container"},
				Image: "nginx:latest",
			},
			expected: "web-api",
		},
		{
			name: "Container name fallback",
			summary: container.Summary{
				Labels: map[string]string{},
				Names:  []string{"/my-custom-app"},
				Image:  "my-app:1.0",
			},
			expected: "my-custom-app",
		},
		{
			name: "Image name fallback",
			summary: container.Summary{
				Labels: map[string]string{},
				Names:  []string{},
				Image:  "redis:alpine",
			},
			expected: "redis:alpine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.resolveServiceName(tt.summary); got != tt.expected {
				t.Errorf("resolveServiceName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
