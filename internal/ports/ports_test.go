package ports

import "testing"

func TestSuggestNext(t *testing.T) {
	tests := []struct {
		name     string
		occupied map[int]bool
		expected int
	}{
		{
			name:     "Empty active ports",
			occupied: map[int]bool{},
			expected: 3000,
		},
		{
			name: "Port 3000 occupied",
			occupied: map[int]bool{
				3000: true,
			},
			expected: 3001,
		},
		{
			name: "Gap in ports",
			occupied: map[int]bool{
				3000: true,
				3002: true,
			},
			expected: 3001,
		},
		{
			name: "Multiple ports occupied",
			occupied: map[int]bool{
				3000: true,
				3001: true,
				3002: true,
				3003: true,
			},
			expected: 3004,
		},
		{
			name: "Port outside range occupied",
			occupied: map[int]bool{
				80:  true,
				443: true,
			},
			expected: 3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SuggestNext(tt.occupied); got != tt.expected {
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
