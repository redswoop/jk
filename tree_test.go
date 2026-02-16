package main

import "testing"

func TestBuildTree(t *testing.T) {
	// Build a three-level hierarchy:
	// root (pid=1, ppid=0)
	//   ├─ child1 (pid=2, ppid=1)
	//   │  └─ grandchild (pid=4, ppid=2)
	//   └─ child2 (pid=3, ppid=1)
	processes := []*Process{
		{PID: 1, PPID: 0, ElapsedSec: 100, Parsed: "root"},
		{PID: 2, PPID: 1, ElapsedSec: 50, Parsed: "child1"},
		{PID: 3, PPID: 1, ElapsedSec: 30, Parsed: "child2"},
		{PID: 4, PPID: 2, ElapsedSec: 10, Parsed: "grandchild"},
	}

	result := buildTree(processes)

	if len(result) != 4 {
		t.Fatalf("got %d processes, want 4", len(result))
	}

	// Root should be first with no prefix
	if result[0].PID != 1 || result[0].TreePrefix != "" {
		t.Errorf("root: PID=%d, prefix=%q", result[0].PID, result[0].TreePrefix)
	}

	// Children sorted by elapsed (youngest first): child2 (30s) before child1 (50s)
	if result[1].PID != 3 {
		t.Errorf("first child should be PID 3 (youngest), got %d", result[1].PID)
	}
	if result[1].TreePrefix != "├─ " {
		t.Errorf("first child prefix = %q, want '├─ '", result[1].TreePrefix)
	}

	if result[2].PID != 2 {
		t.Errorf("second child should be PID 2, got %d", result[2].PID)
	}
	if result[2].TreePrefix != "└─ " {
		t.Errorf("second child prefix = %q, want '└─ '", result[2].TreePrefix)
	}

	// Grandchild under child1
	if result[3].PID != 4 {
		t.Errorf("grandchild should be PID 4, got %d", result[3].PID)
	}
	if result[3].TreePrefix != "   └─ " {
		t.Errorf("grandchild prefix = %q, want '   └─ '", result[3].TreePrefix)
	}
}

func TestBuildTreeMultipleRoots(t *testing.T) {
	processes := []*Process{
		{PID: 10, PPID: 0, ElapsedSec: 200, Parsed: "root1"},
		{PID: 20, PPID: 0, ElapsedSec: 100, Parsed: "root2"},
		{PID: 30, PPID: 10, ElapsedSec: 50, Parsed: "child-of-root1"},
	}

	result := buildTree(processes)

	if len(result) != 3 {
		t.Fatalf("got %d processes, want 3", len(result))
	}

	// Roots sorted by elapsed: root2 (100s) before root1 (200s)
	if result[0].PID != 20 {
		t.Errorf("first root should be PID 20, got %d", result[0].PID)
	}
	if result[1].PID != 10 {
		t.Errorf("second root should be PID 10, got %d", result[1].PID)
	}
	if result[2].PID != 30 {
		t.Errorf("child should be PID 30, got %d", result[2].PID)
	}
}

func TestAddParentNames(t *testing.T) {
	processes := []*Process{
		{PID: 1, PPID: 0, Parsed: "parent-process"},
		{PID: 2, PPID: 1, Parsed: "child-process"},
		{PID: 3, PPID: 99, Parsed: "orphan"},
	}

	addParentNames(processes)

	if processes[0].ParentName != "" {
		t.Errorf("root parent name = %q, want empty", processes[0].ParentName)
	}
	if processes[1].ParentName != "parent-process" {
		t.Errorf("child parent name = %q, want parent-process", processes[1].ParentName)
	}
	if processes[2].ParentName != "" {
		t.Errorf("orphan parent name = %q, want empty", processes[2].ParentName)
	}
}

func TestAddParentNamesTruncation(t *testing.T) {
	processes := []*Process{
		{PID: 1, PPID: 0, Parsed: "a-very-long-process-name"},
		{PID: 2, PPID: 1, Parsed: "child"},
	}

	addParentNames(processes)

	if len(processes[1].ParentName) != 15 {
		t.Errorf("parent name length = %d, want 15 (truncated)", len(processes[1].ParentName))
	}
}
