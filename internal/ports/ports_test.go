package ports

import "testing"

func TestSuggestNext(t *testing.T) {
	tests := []struct {
		name        string
		activePorts []Port
		expected    int
	}{
		{
			name:        "Empty active ports",
			activePorts: []Port{},
			expected:    3000,
		},
		{
			name: "Port 3000 occupied",
			activePorts: []Port{
				{Number: 3000},
			},
			expected: 3001,
		},
		{
			name: "Gap in ports",
			activePorts: []Port{
				{Number: 3000},
				{Number: 3002},
			},
			expected: 3001,
		},
		{
			name: "Multiple ports occupied",
			activePorts: []Port{
				{Number: 3000},
				{Number: 3001},
				{Number: 3002},
				{Number: 3003},
			},
			expected: 3004,
		},
		{
			name: "Port outside range occupied",
			activePorts: []Port{
				{Number: 80},
				{Number: 443},
			},
			expected: 3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SuggestNext(tt.activePorts); got != tt.expected {
				t.Errorf("SuggestNext() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSortPorts(t *testing.T) {
	ports := []Port{
		{Number: 80},
		{Number: 22},
		{Number: 443},
	}
	expected := []Port{
		{Number: 22},
		{Number: 80},
		{Number: 443},
	}

	SortPorts(ports)

	for i := range ports {
		if ports[i].Number != expected[i].Number {
			t.Errorf("at index %d: expected %d, got %d", i, expected[i].Number, ports[i].Number)
		}
	}
}
