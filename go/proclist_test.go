package main

import "testing"

func TestParsePsOutput(t *testing.T) {
	output := `  1234     1   05:30 /usr/local/bin/node server.js
  5678     1 1-12:00:00 /usr/bin/python3 /home/user/app.py
  9012  1234      0:45 /usr/bin/ruby script.rb
`

	procs := parsePsOutput(output)

	if len(procs) != 3 {
		t.Fatalf("got %d processes, want 3", len(procs))
	}

	// Check first process
	if procs[0].PID != 1234 {
		t.Errorf("proc[0].PID = %d, want 1234", procs[0].PID)
	}
	if procs[0].PPID != 1 {
		t.Errorf("proc[0].PPID = %d, want 1", procs[0].PPID)
	}
	if procs[0].Parsed != "node:server.js" {
		t.Errorf("proc[0].Parsed = %q, want node:server.js", procs[0].Parsed)
	}
	if procs[0].ElapsedSec != 330 {
		t.Errorf("proc[0].ElapsedSec = %d, want 330", procs[0].ElapsedSec)
	}

	// Check second process (with days)
	if procs[1].PID != 5678 {
		t.Errorf("proc[1].PID = %d, want 5678", procs[1].PID)
	}
	if procs[1].ElapsedSec != 129600 {
		t.Errorf("proc[1].ElapsedSec = %d, want 129600", procs[1].ElapsedSec)
	}

	// Check third process
	if procs[2].PID != 9012 {
		t.Errorf("proc[2].PID = %d, want 9012", procs[2].PID)
	}
	if procs[2].PPID != 1234 {
		t.Errorf("proc[2].PPID = %d, want 1234", procs[2].PPID)
	}
}

func TestParsePsOutputFiltersKernel(t *testing.T) {
	output := `    1     0      0:01 /sbin/launchd
  100     1      0:00 [kworker/0:1]
  200     1      0:05 /usr/sbin/kernel_task
  300     1      0:10 /usr/local/bin/node app.js
`

	procs := parsePsOutput(output)

	// Should skip [kworker] and kernel_task
	if len(procs) != 2 {
		t.Fatalf("got %d processes, want 2", len(procs))
	}
	if procs[0].PID != 1 {
		t.Errorf("proc[0].PID = %d, want 1", procs[0].PID)
	}
	if procs[1].PID != 300 {
		t.Errorf("proc[1].PID = %d, want 300", procs[1].PID)
	}
}

func TestParsePsOutputEmpty(t *testing.T) {
	procs := parsePsOutput("")
	if len(procs) != 0 {
		t.Errorf("expected empty, got %d processes", len(procs))
	}
}

func TestSplitPsFields(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantPID int
		wantOK  bool
	}{
		{"normal", "  1234     1   05:30 /usr/bin/node", 1234, true},
		{"tight spacing", "1 2 0:01 cmd", 1, true},
		{"too few fields", "1234 5678", 0, false},
		{"empty", "", 0, false},
		{"non-numeric pid", "abc 1 0:01 cmd", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pid, _, _, _, ok := splitPsFields(tt.line)
			if ok != tt.wantOK {
				t.Errorf("ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && pid != tt.wantPID {
				t.Errorf("pid = %d, want %d", pid, tt.wantPID)
			}
		})
	}
}
