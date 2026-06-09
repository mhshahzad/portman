package scanner

import (
	"reflect"
	"testing"

	"github.com/mhshahzad/portman/internal/ports"
)

func TestLsofParse(t *testing.T) {
	input := `COMMAND     PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
sshd        123   root    3u  IPv4 0x1234567890abcdef      0t0  TCP *:22 (LISTEN)
sshd        123   root    4u  IPv6 0x1234567890abcdef      0t0  TCP *:22 (LISTEN)
ControlCe   456  hamza   19u  IPv4 0xabcdef1234567890      0t0  TCP *:7000 (LISTEN)
ControlCe   456  hamza   21u  IPv6 0xabcdef1234567890      0t0  TCP *:7000 (LISTEN)
node       7890  hamza   23u  IPv4 0x9876543210fedcba      0t0  TCP 127.0.0.1:8080 (LISTEN)
`
	expected := []ports.Port{
		{Number: 22, Protocol: "tcp", ProcessName: "sshd", PID: 123, Source: "lsof"},
		{Number: 7000, Protocol: "tcp", ProcessName: "ControlCe", PID: 456, Source: "lsof"},
		{Number: 8080, Protocol: "tcp", ProcessName: "node", PID: 7890, Source: "lsof"},
	}

	scanner := &LsofScanner{}
	result := scanner.parse(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
