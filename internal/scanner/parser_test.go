package scanner

import (
	"reflect"
	"testing"

	"github.com/mhshahzad/portman/internal/ports"
)

func TestParseSSOutput(t *testing.T) {
	input := `Netid  State      Recv-Q Send-Q Local Address:Port               Peer Address:Port
tcp    LISTEN     0      128    0.0.0.0:22                       0.0.0.0:*                   users:(("sshd",pid=1122,fd=3))
tcp    LISTEN     0      100    127.0.0.1:5432                   0.0.0.0:*                   users:(("postgres",pid=1234,fd=5))
udp    UNCONN     0      0      0.0.0.0:68                       0.0.0.0:*                   users:(("dhclient",pid=987,fd=6))
tcp    LISTEN     0      128    [::]:80                          [::]:*                      users:(("nginx",pid=2000,fd=7))
tcp    LISTEN     0      128    127.0.0.1:8080                   0.0.0.0:*
`
	expected := []ports.Port{
		{Number: 22, Protocol: "tcp", ProcessName: "sshd", PID: 1122},
		{Number: 5432, Protocol: "tcp", ProcessName: "postgres", PID: 1234},
		{Number: 68, Protocol: "udp", ProcessName: "dhclient", PID: 987},
		{Number: 80, Protocol: "tcp", ProcessName: "nginx", PID: 2000},
		{Number: 8080, Protocol: "tcp", ProcessName: "", PID: 0},
	}

	result := ParseSSOutput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
