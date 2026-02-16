package main

import "testing"

func TestParseProcNetTcp(t *testing.T) {
	// Real /proc/net/tcp content (abbreviated)
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:1F91 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12346 1 0000000000000000 100 0 0 10 0
   2: 0100007F:C5E8 0100007F:1F90 01 00000000:00000000 00:00000000 00000000  1000        0 99999 1 0000000000000000 100 0 0 10 0
`
	// 0x1F90 = 8080, 0x1F91 = 8081, state 0A = LISTEN, state 01 = ESTABLISHED

	result := parseProcNetTcp(content)

	if port, ok := result["12345"]; !ok || port != "8080" {
		t.Errorf("inode 12345: got %q, want 8080", port)
	}
	if port, ok := result["12346"]; !ok || port != "8081" {
		t.Errorf("inode 12346: got %q, want 8081", port)
	}
	// ESTABLISHED socket should not be included
	if _, ok := result["99999"]; ok {
		t.Error("inode 99999 should not be included (not LISTEN)")
	}
}

func TestParseProcNetTcpEmpty(t *testing.T) {
	result := parseProcNetTcp("")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}
