package main

import "testing"

func TestParseLsofOutput(t *testing.T) {
	output := `COMMAND     PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node      12345  armen   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
node      12345  armen   24u  IPv4 0x1234567891      0t0  TCP *:3001 (LISTEN)
python3   23456  armen   5u   IPv4 0x2345678901      0t0  TCP 127.0.0.1:8080 (LISTEN)
postgres  34567  armen   6u   IPv4 0x3456789012      0t0  TCP *:5432 (LISTEN)
`

	result := parseLsofOutput(output)

	// node should have ports 3000 and 3001
	if ports, ok := result[12345]; !ok {
		t.Error("missing PID 12345")
	} else if len(ports) != 2 || ports[0] != "3000" || ports[1] != "3001" {
		t.Errorf("PID 12345 ports = %v, want [3000 3001]", ports)
	}

	// python should have port 8080
	if ports, ok := result[23456]; !ok {
		t.Error("missing PID 23456")
	} else if len(ports) != 1 || ports[0] != "8080" {
		t.Errorf("PID 23456 ports = %v, want [8080]", ports)
	}

	// postgres should have port 5432
	if ports, ok := result[34567]; !ok {
		t.Error("missing PID 34567")
	} else if len(ports) != 1 || ports[0] != "5432" {
		t.Errorf("PID 34567 ports = %v, want [5432]", ports)
	}
}

func TestParseLsofOutputEmpty(t *testing.T) {
	result := parseLsofOutput("")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestParseLsofOutputMalformedLines(t *testing.T) {
	output := `COMMAND     PID   USER   FD   TYPE
short line
node      notapid  armen   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
node      12345  armen   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
`

	result := parseLsofOutput(output)

	// Only the valid line should be parsed
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
	if _, ok := result[12345]; !ok {
		t.Error("missing PID 12345")
	}
}

func TestParseLsofOutputDedup(t *testing.T) {
	// Same port appearing multiple times for same PID
	output := `COMMAND     PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node      12345  armen   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
node      12345  armen   24u  IPv6 0x1234567891      0t0  TCP *:3000 (LISTEN)
`

	result := parseLsofOutput(output)

	if ports := result[12345]; len(ports) != 1 {
		t.Errorf("expected 1 unique port, got %v", ports)
	}
}

func TestParseSsOutput(t *testing.T) {
	output := `State  Recv-Q Send-Q  Local Address:Port   Peer Address:Port  Process
LISTEN 0      128           0.0.0.0:22          0.0.0.0:*      users:(("sshd",pid=1234,fd=3))
LISTEN 0      511           0.0.0.0:3000        0.0.0.0:*      users:(("node",pid=5678,fd=23))
LISTEN 0      128           0.0.0.0:8080        0.0.0.0:*      users:(("python3",pid=9012,fd=5))
LISTEN 0      511                 *:3001              *:*      users:(("node",pid=5678,fd=24))
`

	result := parseSsOutput(output)

	// sshd on port 22
	if ports, ok := result[1234]; !ok {
		t.Error("missing PID 1234")
	} else if len(ports) != 1 || ports[0] != "22" {
		t.Errorf("PID 1234 ports = %v, want [22]", ports)
	}

	// node on ports 3000 and 3001
	if ports, ok := result[5678]; !ok {
		t.Error("missing PID 5678")
	} else if len(ports) != 2 || ports[0] != "3000" || ports[1] != "3001" {
		t.Errorf("PID 5678 ports = %v, want [3000 3001]", ports)
	}

	// python on port 8080
	if ports, ok := result[9012]; !ok {
		t.Error("missing PID 9012")
	} else if len(ports) != 1 || ports[0] != "8080" {
		t.Errorf("PID 9012 ports = %v, want [8080]", ports)
	}
}

func TestParseSsOutputMultiplePids(t *testing.T) {
	// nginx master + worker sharing a port
	output := `State  Recv-Q Send-Q  Local Address:Port   Peer Address:Port  Process
LISTEN 0      511           0.0.0.0:80          0.0.0.0:*      users:(("nginx",pid=100,fd=6),("nginx",pid=101,fd=6))
`

	result := parseSsOutput(output)

	if ports := result[100]; len(ports) != 1 || ports[0] != "80" {
		t.Errorf("PID 100 ports = %v, want [80]", ports)
	}
	if ports := result[101]; len(ports) != 1 || ports[0] != "80" {
		t.Errorf("PID 101 ports = %v, want [80]", ports)
	}
}

func TestParseSsOutputNoProcessInfo(t *testing.T) {
	// When running as non-root, process info may be missing
	output := `State  Recv-Q Send-Q  Local Address:Port   Peer Address:Port  Process
LISTEN 0      128           0.0.0.0:22          0.0.0.0:*
LISTEN 0      511           0.0.0.0:3000        0.0.0.0:*      users:(("node",pid=5678,fd=23))
`

	result := parseSsOutput(output)

	// Only the line with process info should produce results
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
	if _, ok := result[5678]; !ok {
		t.Error("missing PID 5678")
	}
}

func TestParseSsOutputEmpty(t *testing.T) {
	result := parseSsOutput("")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
